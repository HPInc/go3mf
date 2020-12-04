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
	ErrMeshConsistency        = errors.New("mesh has non-manifold edges without consistent triangle orientation")
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

func NewMissingFieldError(name string) error {
	return &MissingFieldError{Name: name}
}

type MissingFieldError struct {
	Name string
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("required field '%s' is not set", e.Name)
}

type ParseAttrError struct {
	Name     string
	Required bool
}

func NewParseAttrError(name string, required bool) *ParseAttrError {
	return &ParseAttrError{name, required}
}

func (e *ParseAttrError) Error() string {
	req := "required"
	if !e.Required {
		req = "optional"
	}
	return fmt.Sprintf("error parsing %s attribute '%s'", req, e.Name)
}
