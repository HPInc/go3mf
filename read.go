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
	"unsafe"

	xml3mf "github.com/qmuntal/go3mf/internal/xml"
)

var checkEveryTokens = 1000

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
	modelName := xml.Name{Space: Namespace, Local: attrModel}
	if name == modelName {
		child = &modelDecoder{model: d.model}
	}
	return
}

func decodeModelFile(ctx context.Context, r io.Reader, model *Model, path string, isRoot, strict bool) (*Scanner, error) {
	x := xml3mf.NewDecoder(r)
	scanner := Scanner{
		extensionDecoder: make(map[string]SpecDecoder),
		IsRoot:           isRoot,
		ModelPath:        path,
	}
	for _, ext := range model.Specs {
		if ext, ok := ext.(SpecDecoder); ok {
			scanner.extensionDecoder[ext.Namespace()] = ext
		}
	}
	state, names := make([]NodeDecoder, 0, 10), make([]xml.Name, 0, 10)

	var (
		currentDecoder, tmpDecoder NodeDecoder
		currentName                xml.Name
	)
	currentDecoder = &topLevelDecoder{isRoot: isRoot, model: model}
	currentDecoder.SetScanner(&scanner)
	var err error
	x.OnStart = func(tp xml3mf.StartElement) {
		tmpDecoder = currentDecoder.Child(tp.Name)
		if tmpDecoder != nil {
			tmpDecoder.SetScanner(&scanner)
			state = append(state, currentDecoder)
			names = append(names, currentName)
			scanner.contex = append(names, tp.Name)
			currentName = tp.Name
			currentDecoder = tmpDecoder
			currentDecoder.Start(*(*[]XMLAttr)(unsafe.Pointer(&tp.Attr)))
		}
	}
	x.OnEnd = func(tp xml.EndElement) {
		if currentName == tp.Name {
			currentDecoder.End()
			currentDecoder, state = state[len(state)-1], state[:len(state)-1]
			currentName, names = names[len(names)-1], names[:len(names)-1]
		}
	}
	x.OnChar = func(tp xml.CharData) {
		currentDecoder.Text(tp)
	}
	var i int
	for {
		err = x.RawToken()
		if err != nil || (strict && scanner.Err.Len() != 0) {
			break
		}
		if i%checkEveryTokens == 0 {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			default: // Default is must to avoid blocking
			}
			if err != nil {
				break
			}
		}
		i++
	}
	if err == io.EOF {
		err = nil
	}
	if err == nil && scanner.Err.Len() != 0 {
		if strict || scanner.Err.Len() == 1 {
			err = scanner.Err.Unwrap()
		} else {
			err = &scanner.Err
		}
	}
	return &scanner, err
}

// Decoder implements a 3mf file decoder.
type Decoder struct {
	Strict        bool
	p             packageReader
	flate         func(r io.Reader) io.ReadCloser
	nonRootModels []packageFile
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

// UnmarshalModel fills a model with the data of a root model file
// using not strict mode.
func UnmarshalModel(data []byte, model *Model) error {
	d := NewDecoder(nil, 0)
	d.Strict = false
	return d.processRootModel(context.Background(), &fakePackageFile{data: data}, model)
}

func (d *Decoder) processRootModel(ctx context.Context, rootFile packageFile, model *Model) error {
	f, err := rootFile.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	scanner, err := decodeModelFile(ctx, f, model, rootFile.Name(), true, d.Strict)
	if err != nil {
		return err
	}
	d.addModelFile(scanner, model)
	for _, ext := range scanner.extensionDecoder {
		ext.OnDecoded(model)
	}
	return nil
}

func (d *Decoder) addChildModelFile(p *Scanner, model *Model) {
	model.Childs[p.ModelPath].Resources = p.Resources
}

func (d *Decoder) addModelFile(p *Scanner, model *Model) {
	for _, bi := range p.BuildItems {
		model.Build.Items = append(model.Build.Items, bi)
	}
	model.Resources = p.Resources
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
				return nil, errors.New("package root model points to an unexisting file")
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
		return nil, errors.New("package does not have root model")
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
	scanner, err := decodeModelFile(ctx, file, model, attachment.Name(), false, d.Strict)
	return scanner, err
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
