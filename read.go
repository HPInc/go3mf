package go3mf

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"sync"
)

// ExtensionDecoder is the contract that should be implemented
// in order to enable automatic extension decoding.
// NodeDecoder should return a NodeDecoder that will do the real decoding.
// DecodeAttribute should parse the attribute and update the parentNode.
type ExtensionDecoder interface {
	NodeDecoder(parentNode interface{}, nodeName string) NodeDecoder
	DecodeAttribute(s *Scanner, parentNode interface{}, attr xml.Attr)
}

var extensionDecoder = make(map[string]ExtensionDecoder)

// RegisterExtensionDecoder registers a ExtensionDecoder.
func RegisterExtensionDecoder(key string, e ExtensionDecoder) {
	extensionDecoder[key] = e
}

// A XMLDecoder is anything that can decode a stream of XML tokens, including a Decoder.
type XMLDecoder interface {
	xml.TokenReader
	// Skip reads tokens until it has consumed the end element matching the most recent start element already consumed.
	Skip() error
	// InputOffset returns the input stream byte offset of the current decoder position.
	InputOffset() int64
}

type relationship interface {
	Type() string
	TargetURI() string
}

type packageFile interface {
	Name() string
	FindFileFromRel(string) (packageFile, bool)
	FindFileFromName(string) (packageFile, bool)
	Relationships() []relationship
	Open() (io.ReadCloser, error)
}

type packageReader interface {
	Open(func(r io.Reader) io.ReadCloser) error
	FindFileFromRel(string) (packageFile, bool)
	FindFileFromName(string) (packageFile, bool)
}

// ReadCloser wrapps a Decoder than can be closed.
type ReadCloser struct {
	f *os.File
	*Decoder
}

// OpenReader will open the 3MF file specified by name and return a ReadCloser.
func OpenReader(name string) (*ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	return &ReadCloser{f: f, Decoder: NewDecoder(f, fi.Size())}, nil
}

// Close closes the 3MF file, rendering it unusable for I/O.
func (r *ReadCloser) Close() error {
	return r.f.Close()
}

type topLevelDecoder struct {
	BaseDecoder
	model  *Model
	isRoot bool
}

func (d *topLevelDecoder) Child(name xml.Name) (child NodeDecoder) {
	modelName := xml.Name{Space: ExtensionName, Local: attrModel}
	if name == modelName {
		child = &modelDecoder{model: d.model}
	}
	return
}

// modelFileDecoder cannot be reused between goroutines.
type modelFileDecoder struct {
	Scanner *Scanner
}

func (d *modelFileDecoder) Decode(ctx context.Context, x XMLDecoder, model *Model, path string, isRoot, strict bool) error {
	d.Scanner = NewScanner(model)
	d.Scanner.IsRoot = isRoot
	d.Scanner.Strict = strict
	d.Scanner.ModelPath = path
	state := make([]NodeDecoder, 0, 10)
	names := make([]xml.Name, 0, 10)

	var (
		currentDecoder NodeDecoder
		tmpDecoder     NodeDecoder
		currentName    xml.Name
		t              xml.Token
	)
	nextBytesCheck := checkEveryBytes
	currentDecoder = &topLevelDecoder{isRoot: isRoot, model: model}
	currentDecoder.SetScanner(d.Scanner)

	for {
		t, d.Scanner.Err = x.Token()
		if d.Scanner.Err != nil {
			break
		}
		switch tp := t.(type) {
		case xml.StartElement:
			tmpDecoder = currentDecoder.Child(tp.Name)
			if tmpDecoder != nil {
				tmpDecoder.SetScanner(d.Scanner)
				state = append(state, currentDecoder)
				names = append(names, currentName)
				currentName = tp.Name
				d.Scanner.Element = tp.Name.Local
				currentDecoder = tmpDecoder
				currentDecoder.Open()
				currentDecoder.Attributes(tp.Attr)
			} else {
				d.Scanner.Err = x.Skip()
			}
		case xml.CharData:
			currentDecoder.Text(tp)
		case xml.EndElement:
			if currentName == tp.Name {
				d.Scanner.Element = tp.Name.Local
				currentDecoder.Close()
				currentDecoder, state = state[len(state)-1], state[:len(state)-1]
				currentName, names = names[len(names)-1], names[:len(names)-1]
			}
			if x.InputOffset() > nextBytesCheck {
				select {
				case <-ctx.Done():
					d.Scanner.Err = ctx.Err()
				default: // Default is must to avoid blocking
				}
				nextBytesCheck += checkEveryBytes
			}
		}
		if d.Scanner.Err != nil {
			break
		}
	}
	if d.Scanner.Err == io.EOF {
		d.Scanner.Err = nil
	}
	return d.Scanner.Err
}

// Decoder implements a 3mf file decoder.
type Decoder struct {
	Strict           bool
	Warnings         []error
	p                packageReader
	x                func(r io.Reader) XMLDecoder
	flate            func(r io.Reader) io.ReadCloser
	productionModels map[string]packageFile
	ctx              context.Context
}

// NewDecoder returns a new Decoder reading a 3mf file from r.
func NewDecoder(r io.ReaderAt, size int64) *Decoder {
	return &Decoder{
		p:      &opcReader{ra: r, size: size},
		Strict: true,
	}
}

// Decode reads the 3mf file and unmarshall its content into the model.
func (d *Decoder) Decode(model *Model) error {
	return d.DecodeContext(context.Background(), model)
}

// SetXMLDecoder sets the XML decoder to use when reading XML files.
func (d *Decoder) SetXMLDecoder(x func(r io.Reader) XMLDecoder) {
	d.x = x
}

// SetDecompressor sets or overrides a custom decompressor for deflating the zip package.
func (d *Decoder) SetDecompressor(dcomp func(r io.Reader) io.ReadCloser) {
	d.flate = dcomp
}

// DecodeContext reads the 3mf file and unmarshall its content into the model.
func (d *Decoder) DecodeContext(ctx context.Context, model *Model) error {
	rootFile, err := d.processOPC(model)
	if err != nil {
		return err
	}
	if err := d.processNonRootModels(ctx, model); err != nil {
		return err
	}
	return d.processRootModel(ctx, rootFile, model)
}

func (d *Decoder) tokenReader(r io.Reader) XMLDecoder {
	if d.x == nil {
		return xml.NewDecoder(r)
	}
	return d.x(r)
}

// DecodeRawModel fills a model with the raw content of one model file.
func (d *Decoder) DecodeRawModel(ctx context.Context, model *Model, content string) error {
	return d.processRootModel(ctx, &fakePackageFile{str: content}, model)
}

func (d *Decoder) processRootModel(ctx context.Context, rootFile packageFile, model *Model) error {
	f, err := rootFile.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	mf := modelFileDecoder{}
	err = mf.Decode(ctx, d.tokenReader(f), model, rootFile.Name(), true, d.Strict)
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default: // Default is must to avoid blocking
	}
	d.addModelFile(mf.Scanner, model)
	return err
}

func (d *Decoder) addModelFile(p *Scanner, model *Model) {
	for _, bi := range p.BuildItems {
		model.Build.Items = append(model.Build.Items, bi)
	}
	for _, res := range p.Resources {
		model.Resources = append(model.Resources, res)
	}
	for _, res := range p.Warnings {
		d.Warnings = append(d.Warnings, res)
	}
}

func (d *Decoder) processNonRootModels(ctx context.Context, model *Model) (err error) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var files sync.Map
	prodAttCount := len(model.ProductionAttachments)
	wg.Add(prodAttCount)
	for i := 0; i < prodAttCount; i++ {
		go func(i int) {
			defer wg.Done()
			f, err1 := d.readProductionAttachmentModel(ctx, i, model)
			select {
			case <-ctx.Done():
				return // Error somewhere, terminate
			default: // Default is must to avoid blocking
			}
			if err1 != nil {
				err = err1
				cancel()
			}
			files.Store(i, f)
		}(i)
	}
	wg.Wait()
	if err != nil {
		return err
	}
	indices := make([]int, 0, prodAttCount)
	files.Range(func(key, value interface{}) bool {
		indices = append(indices, key.(int))
		return true
	})
	sort.Ints(indices)
	for _, index := range indices {
		f, _ := files.Load(index)
		d.addModelFile(f.(*Scanner), model)
	}
	return nil
}

func (d *Decoder) processOPC(model *Model) (packageFile, error) {
	err := d.p.Open(d.flate)
	if err != nil {
		return nil, err
	}
	rootFile, ok := d.p.FindFileFromRel(relTypeModel3D)
	if !ok {
		return nil, errors.New("go3mf: package does not have root model")
	}

	model.Path = rootFile.Name()
	d.extractTexturesAttachments(rootFile, model)
	d.extractModelAttachments(rootFile, model)
	for _, a := range model.ProductionAttachments {
		file, _ := d.p.FindFileFromName(a.Path)
		d.extractTexturesAttachments(file, model)
	}
	return rootFile, nil
}

func (d *Decoder) extractTexturesAttachments(rootFile packageFile, model *Model) {
	for _, rel := range rootFile.Relationships() {
		if rel.Type() != relTypeTexture3D && rel.Type() != relTypeThumbnail {
			continue
		}

		if file, ok := rootFile.FindFileFromName(rel.TargetURI()); ok {
			model.Attachments = d.addAttachment(model.Attachments, file, rel.Type())
		}
	}
}

func (d *Decoder) extractModelAttachments(rootFile packageFile, model *Model) {
	d.productionModels = make(map[string]packageFile)
	for _, rel := range rootFile.Relationships() {
		if rel.Type() != relTypeModel3D {
			continue
		}

		if file, ok := rootFile.FindFileFromName(rel.TargetURI()); ok {
			model.ProductionAttachments = append(model.ProductionAttachments, &ProductionAttachment{
				RelationshipType: rel.Type(),
				Path:             file.Name(),
			})
			d.productionModels[file.Name()] = file
		}
	}
}

func (d *Decoder) addAttachment(attachments []*Attachment, file packageFile, relType string) []*Attachment {
	buff, err := copyFile(file)
	if err == nil {
		return append(attachments, &Attachment{
			RelationshipType: relType,
			Path:             file.Name(),
			Stream:           buff,
		})
	}
	return attachments
}

func (d *Decoder) readProductionAttachmentModel(ctx context.Context, i int, model *Model) (*Scanner, error) {
	attachment := model.ProductionAttachments[i]
	file, err := d.productionModels[attachment.Path].Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	mf := modelFileDecoder{}
	err = mf.Decode(ctx, d.tokenReader(file), model, attachment.Path, false, d.Strict)
	return mf.Scanner, err
}

func copyFile(file packageFile) (io.Reader, error) {
	stream, err := file.Open()
	if err != nil {
		return nil, err
	}
	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, stream)
	stream.Close()
	return buff, err
}

type fakePackageFile struct {
	str string
}

func (f *fakePackageFile) Name() string                                { return "/3d/3dmodel.model" }
func (f *fakePackageFile) FindFileFromRel(string) (packageFile, bool)  { return nil, false }
func (f *fakePackageFile) FindFileFromName(string) (packageFile, bool) { return nil, false }
func (f *fakePackageFile) Relationships() []relationship               { return nil }
func (f *fakePackageFile) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString(f.str)), nil
}
