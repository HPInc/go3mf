package model

import "io"

// Attachement defines the Model Attachment.
type Attachement struct {
	Model            *Model
	Stream           io.Reader
	RelationshipType string
	URI              string
}
