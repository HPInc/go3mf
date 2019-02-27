package model

import (
	"github.com/qmuntal/go3mf/internal/model"
	"bytes"
	"errors"
	"io"

	"github.com/qmuntal/go3mf/internal/progress"
)

// ReadError defines a error while reading a 3mf.
type ReadError struct {
	Level   WarningLevel
	Message string
	Code    int
}

func (e *ReadError) Error() string {
	return e.Message
}

type relationship interface {
	Type() string
	TargetURI() string
}

type packageFile interface {
	Name() string
	FindFileFromRel(string) packageFile
	FindFileFromName(string) packageFile
	Relationships() []relationship
	Open() (io.ReadCloser, error)
}

type packageReader interface {
	FindFileFromRel(string) packageFile
	FindFileFromName(string) packageFile
	Relationships() []relationship
}

// Decoder implements a 3mf file decoder.
type Decoder struct {
	Warnings            []error
	AttachmentRelations []string
	progress            *progress.Monitor
	r                   packageReader
}

// NewDecoder returns a new Decoder reading a 3mf file from r.
func NewDecoder(r io.ReaderAt, size int64) (*Decoder, error) {
	opcr, err := newOPCReader(r, size)
	if err != nil {
		return nil, err
	}
	return &Decoder{
		r:        opcr,
		progress: progress.NewMonitor(),
	}, nil
}

// SetProgressCallback specifies the callback to be executed on every step of the progress.
func (d *Decoder) SetProgressCallback(callback progress.ProgressCallback, userData interface{}) {
	d.progress.SetProgressCallback(callback, userData)
}

// Decode reads the 3mf file and unmarshall its content into the model.
func (d *Decoder) Decode(model *Model) error {
	if err := d.processOPC(model); err != nil {
		return err
	}
	return nil
}

func (d *Decoder) processOPC(model *Model) (io.ReadCloser, error) {
	rootFile := d.r.FindFileFromRel(relTypeModel3D)
	if rootFile == nil {
		return nil, errors.New("go3mf: package does not have root model")
	}

	model.RootPath = rootFile.Name()
	d.extractTexturesAttachments(model, rootFile)
	d.extractCustomAttachments(model, rootFile)
	d.extractModelAttachments(model, rootFile)
	for _, a := range model.ProductionAttachments {
		file := d.r.FindFileFromName(a.Path)
		d.extractCustomAttachments(model, file)
		d.extractTexturesAttachments(model, file)
	}
	thumbFile := rootFile.FindFileFromRel(relTypeThumbnail)
	if thumbFile != nil {
		if buff, err := copyFile(thumbFile); err == nil {
			model.SetThumbnail(buff)
		}
	}

	return rootFile.Open()
}

func (d *Decoder) extractTexturesAttachments(model *Model, rootFile packageFile) {
	for _, r := range rootFile.Relationships() {
		if r.Type() != relTypeTexture3D && r.Type() != relTypeThumbnail {
			continue
		}
		file := rootFile.FindFileFromRel(r.TargetURI())
		if file != nil {
			d.addAttachment(model.Attachments, file, r.Type())
		}		
	}
}

func (d *Decoder) extractCustomAttachments(model *Model, rootFile packageFile) {
	for _, r := range d.AttachmentRelations {
		file := rootFile.FindFileFromRel(r)
		if file != nil {
			d.addAttachment(model.Attachments, file, r)
		}	
	}
}

func (d *Decoder) extractModelAttachments(model *Model, rootFile packageFile) {
	for _, r := range rootFile.Relationships() {
		if r.Type() != relTypeModel3D {
			continue
		}
		file := rootFile.FindFileFromRel(r.TargetURI())
		if file != nil {
			d.addAttachment(model.ProductionAttachments, file, r.Type())
		}		
	}
}

func (d *Decoder) addAttachment(attachments []*Attachment, file packageFile, relType string) error {
	buff, err := copyFile(file)
	if err == nil {
		attachments = append(attachments, &Attachment{
			RelationshipType: relType,
			Path:              file.Name(),
			Stream:           buff,
		})
	}	
	return err
}

func copyFile(file packageFile) (io.Reader, error) {
	stream, err := file.Open()
	if err != nil {
		return err
	}
	buff := new(bytes.Buffer)
	_, err := io.Copy(buff, stream)
	stream.Close()
	return err
}