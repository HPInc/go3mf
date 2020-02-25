package errors

import (
	"errors"
	"fmt"
)

// Error guards.
var (
	ErrMissingID              = errors.New("resource ID MUST be a positive integer")
	ErrDuplicatedID           = errors.New("IDs MUST be unique among all resources under same Model")
	ErrMissingResource        = errors.New("resources MUST be defined prior to referencing")
	ErrDuplicatedIndices      = errors.New("indices v1, v2 and v3 MUST be distinct")
	ErrIndexOutOfBounds       = errors.New("index is bigger than referenced slice")
	ErrInsufficientVertices   = errors.New("mesh has to contain at least 3 vertices to form a solid body")
	ErrInsufficientTriangles  = errors.New("mesh has to contain at least 4 triangles to form a solid body")
	ErrComponentsPID          = errors.New("MUST NOT assign pid to objects that contain components")
	ErrOPCPartName            = errors.New("part name MUST conform to the syntax specified in the OPC specification")
	ErrOPCRelTarget           = errors.New("relationship target part MUST be included in the 3MF document")
	ErrOPCDuplicatedRel       = errors.New("there MUST NOT be more than one relationship of a given type from one part to a second part")
	ErrOPCContentType         = errors.New("part MUST use an appropriate content type specified")
	ErrOPCDuplicatedTicket    = errors.New("each model part MUST attach no more than one PrintTicket")
	ErrOPCDuplicatedModelName = errors.New("go3mf: model part names MUST be unique")
	ErrBaseMaterialGradient   = errors.New("triangle with base material MUST NOT form gradients")
	ErrMetadataName           = errors.New("names without a namespace MUST be restricted to predefined values")
	ErrMetadataNamespace      = errors.New("namespace MUST be declared on the model")
	ErrMetadataDuplicated     = errors.New("names MUST NOT be duplicated")
	ErrOtherItem              = errors.New("MUST NOT reference objects of type other")
	ErrNonObject              = errors.New("MUST NOT reference non-object resources")
	ErrRequiredExt            = errors.New("go3mf: unsupported required extension")
	ErrEmptySlice             = errors.New("slice MUST NOT be empty")
	ErrRecursiveComponent     = errors.New("MUST NOT contain recursive references")
	ErrInvalidObject          = errors.New("MUST contain a mesh or components")
)

type ItemError struct {
	Index int
	Err   error
}

func NewItem(index int, err error) error {
	return &ItemError{Index: index, Err: err}
}

func (e *ItemError) Unwrap() error {
	return e.Err
}

func (e *ItemError) Error() string {
	return fmt.Sprintf("go3mf: build item %d: %v", e.Index, e.Err)
}

type AssetError struct {
	Path  string
	Index int
	Err   error
}

func NewAsset(path string, index int, err error) error {
	return &AssetError{Path: path, Index: index, Err: err}
}

func (e *AssetError) Unwrap() error {
	return e.Err
}

func (e *AssetError) Error() string {
	return fmt.Sprintf("go3mf: asset %s#%d: %v", e.Path, e.Index, e.Err)
}

type ObjectError struct {
	Path  string
	Index int
	Err   error
}

func NewObject(path string, index int, err error) error {
	return &ObjectError{Path: path, Index: index, Err: err}
}

func (e *ObjectError) Unwrap() error {
	return e.Err
}

func (e *ObjectError) Error() string {
	return fmt.Sprintf("go3mf: object %s#%d: %v", e.Path, e.Index, e.Err)
}

type RelationshipError struct {
	Path  string
	Index int
	Err   error
}

func (e *RelationshipError) Unwrap() error {
	return e.Err
}

func (e *RelationshipError) Error() string {
	return fmt.Sprintf("go3mf: relationship %s#%d: %v", e.Path, e.Index, e.Err)
}

type MissingFieldError struct {
	Name string
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("required field %s is not set", e.Name)
}

type ComponentError struct {
	Index int
	Err   error
}

func (e *ComponentError) Unwrap() error {
	return e.Err
}

func (e *ComponentError) Error() string {
	return fmt.Sprintf("component %d: %v", e.Index, e.Err)
}

type TriangleError struct {
	Index int
	Err   error
}

func (e *TriangleError) Unwrap() error {
	return e.Err
}

func (e *TriangleError) Error() string {
	return fmt.Sprintf("triangle %d: %v", e.Index, e.Err)
}

type BaseError struct {
	Index int
	Err   error
}

func (e *BaseError) Unwrap() error {
	return e.Err
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("base %d: %v", e.Index, e.Err)
}

type UnusupportedExtensionErr struct {
	Name string
}

func (e *UnusupportedExtensionErr) Error() string {
	return "go3mf: unupported extension %s" + e.Name
}

type MetadataError struct {
	Index int
	Err   error
}

func (e *MetadataError) Unwrap() error {
	return e.Err
}

func (e *MetadataError) Error() string {
	return fmt.Sprintf("metadata %d: %v", e.Index, e.Err)
}
