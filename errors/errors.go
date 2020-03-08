package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Error guards.
var (
	// core
	ErrMissingID              = errors.New("resource ID MUST be greater than zero")
	ErrDuplicatedID           = errors.New("IDs MUST be unique among all resources under same Model")
	ErrMissingResource        = errors.New("resource MUST be defined prior to referencing")
	ErrDuplicatedIndices      = errors.New("indices v1, v2 and v3 MUST be distinct")
	ErrIndexOutOfBounds       = errors.New("index is bigger than referenced slice")
	ErrInsufficientVertices   = errors.New("mesh MUST contain at least 3 vertices to form a solid body")
	ErrInsufficientTriangles  = errors.New("mesh MUST contain at least 4 triangles to form a solid body")
	ErrComponentsPID          = errors.New("MUST NOT assign pid to objects that contain components")
	ErrOPCPartName            = errors.New("part name MUST conform to the syntax specified in the OPC specification")
	ErrOPCRelTarget           = errors.New("relationship target part MUST be included in the 3MF document")
	ErrOPCDuplicatedRel       = errors.New("there MUST NOT be more than one relationship of a given type from one part to a second part")
	ErrOPCContentType         = errors.New("part MUST use an appropriate content type specified")
	ErrOPCDuplicatedTicket    = errors.New("each model part MUST attach no more than one PrintTicket")
	ErrOPCDuplicatedModelName = errors.New("go3mf: model part names MUST be unique")
	ErrMetadataName           = errors.New("names without a namespace MUST be restricted to predefined values")
	ErrMetadataNamespace      = errors.New("namespace MUST be declared on the model")
	ErrMetadataDuplicated     = errors.New("names MUST NOT be duplicated")
	ErrOtherItem              = errors.New("MUST NOT reference objects of type other")
	ErrNonObject              = errors.New("MUST NOT reference non-object resources")
	ErrRequiredExt            = errors.New("go3mf: unsupported required extension")
	ErrEmptyResourceProps     = errors.New("resource properties MUST NOT be empty")
	ErrRecursiveComponent     = errors.New("MUST NOT contain recursive references")
	ErrInvalidObject          = errors.New("MUST contain a mesh or components")
	// materials
	ErrMultiBlend         = errors.New("there MUST NOT be more blendmethods than layers – 1")
	ErrMaterialMulti      = errors.New("a material, if included, MUST be positioned as the first layer")
	ErrMultiRefMulti      = errors.New("the pids list MUST NOT contain any references to a multiproperties")
	ErrMultiColors        = errors.New("the pids list MUST NOT contain more than one reference to a colorgroup")
	ErrTextureReference   = errors.New("MUST reference to a texture resource")
	ErrCompositeBase      = errors.New("MUST reference to a basematerials group")
	ErrMissingTexturePart = errors.New("texture part MUST be added as an attachment")
	// production
	ErrUUID             = errors.New("UUID MUST be any of the four UUID variants described in IETF RFC 4122")
	ErrProdExtRequired  = errors.New("go3mf: a 3MF package which uses referenced objects MUST enlist the production extension as required")
	ErrProdRefInNonRoot = errors.New("non-root model file components MUST only reference objects in the same model file")
	// slices
	ErrNonSliceStack             = errors.New("slicestackid MUST reference a slice stack resource")
	ErrSlicesAndRefs             = errors.New("may either contain slices or refs, but they MUST NOT contain both element types")
	ErrSliceRefSamePart          = errors.New("the path of the referenced slice stack MUST be different than the path of the original slice stack")
	ErrSliceRefRef               = errors.New("a referenced slice stack MUST NOT contain any further sliceref elements")
	ErrSliceSmallTopZ            = errors.New("slice ztop is smaller than stack zbottom")
	ErrSliceNoMonotonic          = errors.New("the first ztop in the next slicestack MUST be greater than the last ztop in the previous slicestack")
	ErrSliceInsufficientVertices = errors.New("slice MUST contain at least 2 vertices")
	ErrSliceInsufficientPolygons = errors.New("slice MUST contain at least 1 polygon")
	ErrSliceInsufficientSegments = errors.New("slice polygon MUST contain at least 1 segment")
	ErrSlicePolygonNotClosed     = errors.New("all polygons of all slices MUST be closed")
	ErrSliceInvalidTranform      = errors.New("slices transforms MUST be planar")
)

type IndexedError struct {
	Name  string
	Index int
	Err   error
}

func (e *IndexedError) Unwrap() error {
	return e.Err
}

func (e *IndexedError) Error() string {
	return fmt.Sprintf("%s: %v", e.Name, e.Err)
}

type BuildError struct {
	Err error
}

func (e *BuildError) Unwrap() error {
	return e.Err
}

func (e *BuildError) Error() string {
	return fmt.Sprintf("go3mf: build: %v", e.Err)
}

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
	Name  string
	Err   error
}

func NewAsset(path string, index int, asset interface{}, err error) error {
	return &AssetError{Path: path, Index: index, Name: reflect.TypeOf(asset).Elem().Name(), Err: err}
}

func (e *AssetError) Unwrap() error {
	return e.Err
}

func (e *AssetError) Error() string {
	return fmt.Sprintf("go3mf: %s %s#%d: %v", strings.ToLower(e.Name), e.Path, e.Index, e.Err)
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

// A &specerr.ParseFieldError represents an error while decoding a required or an optional property.
// If ResourceID is 0 means that the error took place while parsing the resource property before the ID appeared.
// When Element is 'item' the ResourceID is the objectID property of a build item.
type ParseFieldError struct {
	ResourceID uint32
	ModelPath  string
	Element    string
	Name       string
	Value      string
	Required   bool
}

func (e *ParseFieldError) Error() string {
	req := "required"
	if !e.Required {
		req = "optional"
	}
	return fmt.Sprintf("go3mf: [%s] error parsing property '%s = %s' of element '%s' in resource '%s:%d'", req, e.Name, e.Value, e.Element, e.ModelPath, e.ResourceID)
}
