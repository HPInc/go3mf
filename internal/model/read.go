package model

import (
	"bytes"
	"errors"
	"io"

	"github.com/qmuntal/go3mf/internal/progress"
	"github.com/qmuntal/opc"
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

type packageReader interface {
	FindPartFromRel(string) *opc.File
	FindPart(string) *opc.File
}

type opcReader struct {
	r *opc.Reader
}

func (o *opcReader) FindPartFromRel(relName string) *opc.File {
	name := o.findPartURI(relName)
	if name == "" {
		return nil
	}

	return o.FindPart(name)
}

func (o *opcReader) FindPart(name string) *opc.File {
	for _, f := range o.r.Files {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (o *opcReader) findPartURI(relName string) string {
	for _, r := range o.r.Relationships {
		if r.Type == relName {
			return r.TargetURI
		}
	}
	return ""
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
	opcr, err := opc.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	return &Decoder{
		r:        &opcReader{opcr},
		progress: progress.NewMonitor(),
	}, nil
}

// SetProgressCallback specifies the callback to be executed on every step of the progress.
func (d *Decoder) SetProgressCallback(callback progress.ProgressCallback, userData interface{}) {
	d.progress.SetProgressCallback(callback, userData)
}

func (d *Decoder) Decode(model *Model) error {

	return nil
}

func (d *Decoder) processOPC(model *Model) error {
	rootPart := d.r.FindPart(relTypeRootModel)
	if rootPart == nil {
		return errors.New("go3mf: Package does not have root model.")
	}

	model.RootPath = rootPart.Name

	return nil
}

func (d *Decoder) extractTexturesFromRels(model *Model, rootPart *opc.Part) error {
	for _, r := range rootPart.Relationships {
		if r.Type == relTypeTexture3D || r.Type == relTypeThumbnail {
			part := d.r.FindPart(r.TargetURI)
			if part != nil {
				stream, err := part.Open()
				buff := new(bytes.Buffer)
				if err == nil {
					io.Copy(buff, stream)
					stream.Close()
					model.Attachments = append(model.Attachments, &Attachment{
						RelationshipType: r.Type,
						URI:              part.Name,
						Stream:           buff,
					})
				}
			}
		}
	}
	return nil
}
