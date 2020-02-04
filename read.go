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

type extensionDecoderWrapper struct {
	newNodeDecoder  func(interface{}, string) NodeDecoder
	decodeAttribute func(*Scanner, interface{}, xml.Attr)
	fileFilter      func(string, bool) bool
}

func (e *extensionDecoderWrapper) NewNodeDecoder(parentNode interface{}, nodeName string) NodeDecoder {
	if e.newNodeDecoder != nil {
		return e.newNodeDecoder(parentNode, nodeName)
	}
	return nil
}

func (e *extensionDecoderWrapper) DecodeAttribute(s *Scanner, parentNode interface{}, attr xml.Attr) {
	if e.decodeAttribute != nil {
		e.decodeAttribute(s, parentNode, attr)
	}
}

func (e *extensionDecoderWrapper) FileFilter(relType string, isRootModel bool) bool {
	if e.fileFilter != nil {
		return e.fileFilter(relType, isRootModel)
	}
	return false
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

func (d *modelFileDecoder) Decode(ctx context.Context, x XMLDecoder, model *Model, path string, isRoot, strict bool, extensionDecoder map[string]*extensionDecoderWrapper) error {
	d.Scanner = newScanner()
	if extensionDecoder != nil {
		d.Scanner.extensionDecoder = extensionDecoder
	}
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
	nonRootModels    []packageFile
	ctx              context.Context
	extensionDecoder map[string]*extensionDecoderWrapper
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

// RegisterNodeDecoderExtension registers a node decoding function to the associated extension key.
// The registered function should return a NodeDecoder that will do the real decoding.
func (d *Decoder) RegisterNodeDecoderExtension(key string, f func(parentNode interface{}, nodeName string) NodeDecoder) {
	if e, ok := d.extensionDecoder[key]; ok {
		e.newNodeDecoder = f
	} else {
		if d.extensionDecoder == nil {
			d.extensionDecoder = make(map[string]*extensionDecoderWrapper)
		}
		d.extensionDecoder[key] = &extensionDecoderWrapper{newNodeDecoder: f}
	}
}

// RegisterDecodeAttributeExtension registers a DecodeAttribute function to the associated extension key.
// The registered function should parse the attribute and update the parentNode.
func (d *Decoder) RegisterDecodeAttributeExtension(key string, f func(s *Scanner, parentNode interface{}, attr xml.Attr)) {
	if e, ok := d.extensionDecoder[key]; ok {
		e.decodeAttribute = f
	} else {
		if d.extensionDecoder == nil {
			d.extensionDecoder = make(map[string]*extensionDecoderWrapper)
		}
		d.extensionDecoder[key] = &extensionDecoderWrapper{decodeAttribute: f}
	}
}

// RegisterFileFilterExtension registers a FileFilter function to the associated extension key.
// The registered function should return true if a file with an specific relationship with a model file
// should be preserved as an attachment or not. If the file is accepted and it is a 3dmodel
// it will processed decoded. It can happen that a file is preserved even if this method is
// not called or it is discarded as other packages could accept it.
func (d *Decoder) RegisterFileFilterExtension(key string, f func(relType string, isRootModel bool) bool) {
	if e, ok := d.extensionDecoder[key]; ok {
		e.fileFilter = f
	} else {
		if d.extensionDecoder == nil {
			d.extensionDecoder = make(map[string]*extensionDecoderWrapper)
		}
		if d.extensionDecoder == nil {
			d.extensionDecoder = make(map[string]*extensionDecoderWrapper)
		}
		d.extensionDecoder[key] = &extensionDecoderWrapper{fileFilter: f}
	}
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
	err = mf.Decode(ctx, d.tokenReader(f), model, rootFile.Name(), true, d.Strict, d.extensionDecoder)
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
	for _, ns := range p.Namespaces {
		model.Namespaces = append(model.Namespaces, ns)
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
	nonRootModelsCount := len(d.nonRootModels)
	wg.Add(nonRootModelsCount)
	for i := 0; i < nonRootModelsCount; i++ {
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
	indices := make([]int, 0, nonRootModelsCount)
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
	rootFile, ok := d.p.FindFileFromRel(RelTypeModel3D)
	if !ok {
		return nil, errors.New("go3mf: package does not have root model")
	}

	model.Path = rootFile.Name()
	d.extractCoreAttachments(rootFile, model, true)
	for _, file := range d.nonRootModels {
		d.extractCoreAttachments(file, model, false)
	}
	return rootFile, nil
}

func (d *Decoder) extractCoreAttachments(file packageFile, model *Model, isRoot bool) {
	for _, rel := range file.Relationships() {
		relType := rel.Type()
		preserve := relType == relTypePrintTicket || relType == relTypeThumbnail
		if !preserve {
			for _, ext := range d.extensionDecoder {
				if ext.FileFilter(relType, isRoot) {
					preserve = true
					break
				}
			}
		}
		if !preserve {
			continue
		}
		if file, ok := file.FindFileFromName(rel.TargetURI()); ok {
			if relType == RelTypeModel3D {
				d.nonRootModels = append(d.nonRootModels, file)
			} else {
				model.Attachments = d.addAttachment(model.Attachments, file, relType)
			}
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
	attachment := d.nonRootModels[i]
	file, err := attachment.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	mf := modelFileDecoder{}
	err = mf.Decode(ctx, d.tokenReader(file), model, attachment.Name(), false, d.Strict, d.extensionDecoder)
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

func (f *fakePackageFile) Name() string                                { return uriDefault3DModel }
func (f *fakePackageFile) FindFileFromRel(string) (packageFile, bool)  { return nil, false }
func (f *fakePackageFile) FindFileFromName(string) (packageFile, bool) { return nil, false }
func (f *fakePackageFile) Relationships() []relationship               { return nil }
func (f *fakePackageFile) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString(f.str)), nil
}
