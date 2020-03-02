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
	"strings"
	"sync"
)

var checkEveryBytes = int64(4 * 1024 * 1024)

type extensionDecoderWrapper struct {
	newNodeDecoder  func(interface{}, string) NodeDecoder
	decodeAttribute func(*Scanner, interface{}, xml.Attr)
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

// A XMLDecoder is anything that can decode a stream of XML tokens, including a Decoder.
type XMLDecoder interface {
	xml.TokenReader
	// Skip reads tokens until it has consumed the end element matching the most recent start element already consumed.
	Skip() error
	// InputOffset returns the input stream byte offset of the current decoder position.
	InputOffset() int64
}

type packageFile interface {
	Name() string
	ContentType() string
	FindFileFromName(string) (packageFile, bool)
	Relationships() []Relationship
	Open() (io.ReadCloser, error)
}

type packageReader interface {
	Open(func(r io.Reader) io.ReadCloser) error
	FindFileFromName(string) (packageFile, bool)
	Relationships() []Relationship
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
	baseDecoder
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

func decodeModelFile(ctx context.Context, x XMLDecoder, model *Model, path string, isRoot, strict bool, extensionDecoder map[string]*extensionDecoderWrapper) *Scanner {
	scanner := Scanner{
		extensionDecoder: make(map[string]*extensionDecoderWrapper),
		IsRoot:           isRoot,
		Strict:           strict,
		ModelPath:        path,
	}
	if extensionDecoder != nil {
		scanner.extensionDecoder = extensionDecoder
	}
	state, names := make([]NodeDecoder, 0, 10), make([]xml.Name, 0, 10)

	var (
		currentDecoder, tmpDecoder NodeDecoder
		currentName                xml.Name
		t                          xml.Token
	)
	nextBytesCheck := checkEveryBytes
	currentDecoder = &topLevelDecoder{isRoot: isRoot, model: model}
	currentDecoder.SetScanner(&scanner)

	for {
		t, scanner.Err = x.Token()
		if scanner.Err != nil {
			break
		}
		switch tp := t.(type) {
		case xml.StartElement:
			tmpDecoder = currentDecoder.Child(tp.Name)
			if tmpDecoder != nil {
				tmpDecoder.SetScanner(&scanner)
				state = append(state, currentDecoder)
				names = append(names, currentName)
				currentName = tp.Name
				scanner.Element = tp.Name.Local
				currentDecoder = tmpDecoder
				currentDecoder.Start(tp.Attr)
			} else {
				scanner.Err = x.Skip()
			}
		case xml.CharData:
			currentDecoder.Text(tp)
		case xml.EndElement:
			if currentName == tp.Name {
				scanner.Element = tp.Name.Local
				currentDecoder.End()
				currentDecoder, state = state[len(state)-1], state[:len(state)-1]
				currentName, names = names[len(names)-1], names[:len(names)-1]
			}
			if x.InputOffset() > nextBytesCheck {
				select {
				case <-ctx.Done():
					scanner.Err = ctx.Err()
				default: // Default is must to avoid blocking
				}
				nextBytesCheck += checkEveryBytes
			}
		}
		if scanner.Err != nil {
			break
		}
	}
	if scanner.Err == io.EOF {
		scanner.Err = nil
	}
	return &scanner
}

// Decoder implements a 3mf file decoder.
type Decoder struct {
	Strict           bool
	Warnings         []error
	p                packageReader
	x                func(r io.Reader) XMLDecoder
	flate            func(r io.Reader) io.ReadCloser
	nonRootModels    []packageFile
	extensionDecoder map[string]*extensionDecoderWrapper
}

// NewDecoder returns a new Decoder reading a 3mf file from r.
func NewDecoder(r io.ReaderAt, size int64) *Decoder {
	return &Decoder{
		p:                &opcReader{ra: r, size: size},
		Strict:           true,
		extensionDecoder: make(map[string]*extensionDecoderWrapper),
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

// UnmarshalModel fills a model with the data of a model file.
// This function does not need a decoder initialized with a reader
// so can be initialized as NewDecoder(nil, 0).
func (d *Decoder) UnmarshalModel(data []byte, model *Model) error {
	return d.processRootModel(context.Background(), &fakePackageFile{data: data}, model)
}

func (d *Decoder) processRootModel(ctx context.Context, rootFile packageFile, model *Model) error {
	f, err := rootFile.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := decodeModelFile(ctx, d.tokenReader(f), model, rootFile.Name(), true, d.Strict, d.extensionDecoder)
	select {
	case <-ctx.Done():
		scanner.Err = ctx.Err()
	default: // Default is must to avoid blocking
	}
	d.addModelFile(scanner, model)
	return scanner.Err
}

func (d *Decoder) addChildModelFile(p *Scanner, model *Model) {
	model.Childs[p.ModelPath].Resources = p.Resources
	for _, res := range p.Warnings {
		d.Warnings = append(d.Warnings, res)
	}
}

func (d *Decoder) addModelFile(p *Scanner, model *Model) {
	for _, bi := range p.BuildItems {
		model.Build.Items = append(model.Build.Items, bi)
	}
	model.Resources = p.Resources
	for _, ns := range p.Namespaces {
		model.Namespaces = append(model.Namespaces, ns)
	}
	for _, res := range p.Warnings {
		d.Warnings = append(d.Warnings, res)
	}
}

func (d *Decoder) processNonRootModels(ctx context.Context, model *Model) (err error) {
	var (
		files              sync.Map
		wg                 sync.WaitGroup
		nonRootModelsCount = len(d.nonRootModels)
	)
	wg.Add(nonRootModelsCount)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i := 0; i < nonRootModelsCount; i++ {
		go func(i int) {
			defer wg.Done()
			f, err1 := d.readChildModel(ctx, i, model)
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
		d.addChildModelFile(f.(*Scanner), model)
	}
	return nil
}

func (d *Decoder) processOPC(model *Model) (packageFile, error) {
	if err := d.p.Open(d.flate); err != nil {
		return nil, err
	}
	var rootFile packageFile
	for _, r := range d.p.Relationships() {
		if r.Type == RelType3DModel {
			var ok bool
			rootFile, ok = d.p.FindFileFromName(r.Path)
			if !ok {
				return nil, errors.New("go3mf: package root model points to an unexisting file")
			}
			model.Path = rootFile.Name()
			d.extractCoreAttachments(rootFile, model, true)
			for _, file := range d.nonRootModels {
				d.extractCoreAttachments(file, model, false)
			}
		} else if att, ok := d.p.FindFileFromName(r.Path); ok {
			model.RootRelationships = append(model.RootRelationships, r)
			model.Attachments = d.addAttachment(model.Attachments, att)
		}
	}
	if rootFile == nil {
		return nil, errors.New("go3mf: package does not have root model")
	}
	return rootFile, nil
}

func (d *Decoder) extractCoreAttachments(modelFile packageFile, model *Model, isRoot bool) {
	for _, rel := range modelFile.Relationships() {
		if file, ok := modelFile.FindFileFromName(rel.Path); ok {
			if isRoot {
				if rel.Type == RelType3DModel {
					d.nonRootModels = append(d.nonRootModels, file)
					if model.Childs == nil {
						model.Childs = make(map[string]*ChildModel)
					}
					model.Childs[file.Name()] = new(ChildModel)
				} else {
					model.Attachments = d.addAttachment(model.Attachments, file)
					model.Relationships = append(model.Relationships, rel)
				}
			} else if rel.Type != RelType3DModel {
				if child, ok := model.Childs[modelFile.Name()]; ok {
					model.Attachments = d.addAttachment(model.Attachments, file)
					child.Relationships = append(child.Relationships, rel)
				}
			}
		}
	}
}

func (d *Decoder) addAttachment(attachments []Attachment, file packageFile) []Attachment {
	for _, att := range attachments {
		if strings.EqualFold(att.Path, file.Name()) {
			return attachments
		}
	}
	if buff, err := copyFile(file); err == nil {
		return append(attachments, Attachment{
			Path:        file.Name(),
			Stream:      buff,
			ContentType: file.ContentType(),
		})
	}
	return attachments
}

func (d *Decoder) readChildModel(ctx context.Context, i int, model *Model) (*Scanner, error) {
	attachment := d.nonRootModels[i]
	file, err := attachment.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := decodeModelFile(ctx, d.tokenReader(file), model, attachment.Name(), false, d.Strict, d.extensionDecoder)
	return scanner, scanner.Err
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
	data []byte
}

func (f *fakePackageFile) Name() string                                { return DefaultModelPath }
func (f *fakePackageFile) ContentType() string                         { return ContentType3DModel }
func (f *fakePackageFile) FindFileFromName(string) (packageFile, bool) { return nil, false }
func (f *fakePackageFile) Relationships() []Relationship               { return nil }
func (f *fakePackageFile) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer(f.data)), nil
}
