package errors

import (
	"errors"
	"fmt"
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
	ErrOPCDuplicatedModelName = errors.New("model part names MUST be unique")
	ErrMetadataName           = errors.New("names without a namespace MUST be restricted to predefined values")
	ErrMetadataNamespace      = errors.New("namespace MUST be declared on the model")
	ErrMetadataDuplicated     = errors.New("names MUST NOT be duplicated")
	ErrOtherItem              = errors.New("MUST NOT reference objects of type other")
	ErrNonObject              = errors.New("MUST NOT reference non-object resources")
	ErrRequiredExt            = errors.New("unsupported required extension")
	ErrEmptyResourceProps     = errors.New("resource properties MUST NOT be empty")
	ErrRecursion              = errors.New("MUST NOT contain recursive references")
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
	ErrProdExtRequired  = errors.New("a 3MF package which uses referenced objects MUST enlist the production extension as required")
	ErrProdRefInNonRoot = errors.New("non-root model file components MUST only reference objects in the same model file")
	// slices
	ErrSliceExtRequired          = errors.New("a 3MF package which uses low resolution objects MUST enlist the slice extension as required")
	ErrNonSliceStack             = errors.New("slicestackid MUST reference a slice stack resource")
	ErrSlicesAndRefs             = errors.New("may either contain slices or refs, but they MUST NOT contain both element types")
	ErrSliceRefSamePart          = errors.New("the path of the referenced slice stack MUST be different than the path of the original slice stack")
	ErrSliceRefRef               = errors.New("a referenced slice stack MUST NOT contain any further sliceref elements")
	ErrSliceSmallTopZ            = errors.New("slice ztop is smaller than stack zbottom")
	ErrSliceNoMonotonic          = errors.New("the first ztop in the next slicestack MUST be greater than the last ztop in the previous slicestack")
	ErrSliceInsufficientVertices = errors.New("slice MUST contain at least 2 vertices")
	ErrSliceInsufficientPolygons = errors.New("slice MUST contain at least 1 polygon")
	ErrSliceInsufficientSegments = errors.New("slice polygon MUST contain at least 1 segment")
	ErrSlicePolygonNotClosed     = errors.New("objects with type 'model' and 'solidsupport' MUST not reference slices with open polygons")
	ErrSliceInvalidTranform      = errors.New("any transform applied to an object that references a slice stack MUST be planar")
	// beamlattice
	ErrLatticeObjType       = errors.New("MUST only be added to a mesh object of type model or solidsupport")
	ErrLatticeClippedNoMesh = errors.New("if clipping mode is not equal to none, a clippingmesh resource MUST be specified")
	ErrLatticeInvalidMesh   = errors.New("the clippingmesh and representationmesh MUST be a mesh object of type model and MUST NOT contain a beamlattice")
	ErrLatticeSameVertex    = errors.New("a beam MUST consist of two distinct vertex indices")
	ErrLatticeBeamR2        = errors.New("r2 MUST not be defined, if r1 is not defined")
)

type Level struct {
	Element interface{}
	Index   int // -1 if not needed
}

func (l *Level) String() string {
	name := fmt.Sprintf("%T", l.Element)
	s := strings.Split(name, ".")
	if len(s) > 0 {
		name = s[len(s)-1] // remove package name
	}
	name = strings.Replace(name, "*", "", -1)
	if l.Index == -1 {
		return name
	}
	return fmt.Sprintf("%s#%d", name, l.Index)
}

type Error struct {
	Target []Level
	Err    error
	Path   string
}

func Wrap(err error, element interface{}) error {
	return WrapIndex(err, element, -1)
}

func WrapIndex(err error, element interface{}, index int) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		e.Target = append(e.Target, Level{element, index})
		return e
	}
	if e, ok := err.(*List); ok {
		for i, e1 := range e.Errors {
			e.Errors[i] = WrapIndex(e1, element, index)
		}
		return e
	}
	return &Error{Target: []Level{{element, index}}, Err: err}
}

func WrapPath(err error, element interface{}, path string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		e.Path = path
		e.Target = append(e.Target, Level{element, -1})
		return e
	}
	if e, ok := err.(*List); ok {
		for i, e1 := range e.Errors {
			e.Errors[i] = WrapPath(e1, element, path)
		}
		return e
	}
	return &Error{Target: []Level{{element, -1}}, Err: err, Path: path}
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	levels := make([]string, len(e.Target)+1)
	levels[0] = e.Path
	for i, l := range e.Target {
		levels[len(e.Target)-i] = l.String()
	}
	if e.Path == "" {
		levels = levels[1:]
	}
	return fmt.Sprintf("%s: %v", strings.Join(levels, "@"), e.Err)
}

type MissingFieldError struct {
	Name string
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("required field '%s' is not set", e.Name)
}

// A &specerr.ParseFieldError represents an error while decoding a required or an optional property.
// If ResourceID is 0 means that the error took place while parsing the resource property before the ID appeared.
// When Element is 'item' the ResourceID is the objectID property of a build item.
// Field value is not reported to avoid leaking confidential information.
type ParseFieldError struct {
	Context    string
	ResourceID uint32
	Name       string
	Required   bool
}

func (e *ParseFieldError) Error() string {
	req := "required"
	if !e.Required {
		req = "optional"
	}
	return fmt.Sprintf("%s#%d: error parsing %s attribute '%s'", e.Context, e.ResourceID, req, e.Name)
}
