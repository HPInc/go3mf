package model

import (
	"bytes"
	"errors"
	"io"

	"github.com/qmuntal/go3mf/internal/progress"
)

// ErrUserAborted defines a user function abort.
var ErrUserAborted = errors.New("go3mf: the called function was aborted by the user")

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
	FindFileFromRel(string) (packageFile, bool)
	Relationships() []relationship
	Open() (io.ReadCloser, error)
}

type packageReader interface {
	FindFileFromRel(string) (packageFile, bool)
	FindFileFromName(string) (packageFile, bool)
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
	d.progress.ResetLevels()
	if !d.progress.Progress(0.05, progress.StageExtractOPCPackage) {
		return ErrUserAborted
	}
	_, err := d.processOPC(model)
	if err != nil {
		return err
	}
	if !d.progress.Progress(0.1, progress.StageReadNonRootModels) {
		return ErrUserAborted
	}
	progressNonRoot := 0.6
	if len(model.ProductionAttachments) == 0 {
		progressNonRoot = 0.1
	}
	d.progress.PushLevel(0.1, progressNonRoot)
	// read production attachments
	d.progress.PopLevel()
	if !d.progress.Progress(progressNonRoot, progress.StageReadRootModel) {
		return ErrUserAborted
	}

	return nil
}

func (d *Decoder) processOPC(model *Model) (io.ReadCloser, error) {
	rootFile, ok := d.r.FindFileFromRel(relTypeModel3D)
	if !ok {
		return nil, errors.New("go3mf: package does not have root model")
	}

	model.RootPath = rootFile.Name()
	d.extractTexturesAttachments(model, rootFile)
	d.extractCustomAttachments(model, rootFile)
	d.extractModelAttachments(model, rootFile)
	for _, a := range model.ProductionAttachments {
		file, _ := d.r.FindFileFromName(a.Path)
		d.extractCustomAttachments(model, file)
		d.extractTexturesAttachments(model, file)
	}
	thumbFile, ok := rootFile.FindFileFromRel(relTypeThumbnail)
	if ok {
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

		if file, ok := rootFile.FindFileFromRel(r.TargetURI()); ok {
			model.Attachments = d.addAttachment(model.Attachments, file, r.Type())
		}
	}
}

func (d *Decoder) extractCustomAttachments(model *Model, rootFile packageFile) {
	for _, r := range d.AttachmentRelations {
		if file, ok := rootFile.FindFileFromRel(r); ok {
			model.Attachments = d.addAttachment(model.Attachments, file, r)
		}
	}
}

func (d *Decoder) extractModelAttachments(model *Model, rootFile packageFile) {
	for _, r := range rootFile.Relationships() {
		if r.Type() != relTypeModel3D {
			continue
		}

		if file, ok := rootFile.FindFileFromRel(r.TargetURI()); ok {
			model.ProductionAttachments = d.addAttachment(model.ProductionAttachments, file, r.Type())
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
