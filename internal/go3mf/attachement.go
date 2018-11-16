package go3mf

import "io"

// Attachement defines the Model Attachment.
type Attachement struct {
	Stream           io.Reader
	RelationshipType string
	uri              string
}

// NewAttachement creates a new attachement.
func NewAttachement(stream io.Reader, relType, uri string) *Attachement {
	return &Attachement{
		Stream:           stream,
		RelationshipType: relType,
		uri:              uri,
	}
}

// URI returns the attachement uri.
func (a *Attachement) URI() string {
	return a.uri
}
