package booleanoperations

import (
	"github.com/qmuntal/go3mf"
)

// Namespace is the canonical name of this extension.
const Namespace = "http://www.hp.com/schemas/3dmanufacturing/booleanoperations/2021/02"

var DefaultExtension = go3mf.Extension{
	Namespace:  Namespace,
	LocalName:  "bo",
	IsRequired: false,
}

func init() {
	go3mf.Register(Namespace, Spec{})
}

type Spec struct{}

// Association defines the Association for the Boolean Operation.
// logical  - association defines that the <components> element is an assembly of one or more object components,
// physical -  association defines that the componenents element produces a single part combining the different object componenents.
type Association uint8

// Supported Association.
const (
	Association_logical Association = iota + 1
	Association_physical
)

func newAssociation(s string) (c Association, ok bool) {
	c, ok = map[string]Association{
		"logical":  Association_logical,
		"physical": Association_physical,
	}[s]
	return
}

func (c Association) String() string {
	return map[Association]string{
		Association_logical:  "logical",
		Association_physical: "physical",
	}[c]
}

// A BooleanOperation is an enumerable for the different BooleaOperation.
type BooleanOperation uint8

/**
- union. The new object shape is defined as the merger of the shapes. The new object surface property is defined by the property of the surface property defining the outer surface.
If material and the volumetric property, if available, in the overlapped volume is defined by the added object.
union(a,b,c,d) = ((a Ս b) Ս c) Ս d
- difference. The new object shape is defined by the shape in the first object shape that is not in any other object shape.
The new object surface property, where overlaps, is defined by the object surface property of the substracting object(s).
difference(a,b,c,d) = ((a - b) - c) - d = a - union(b,c,d)
- intersection. The new object shape is defined as the common (clipping) shape in all objects.
The new object surface property is defined as the object surface property of the object clipping that surface.
intersection(a,b,c,d) = ((a Ո b) Ո c) Ո d
*/

// Supported BooleanOperation.
const (
	BooleanOperation_union BooleanOperation = iota + 1
	BooleanOperation_difference
	BooleanOperation_intersection
)

const (
	attrCompsBoolOperOperation   = "operation"
	attrCompsBoolOperAssociation = "association"
)

func newOperation(s string) (t BooleanOperation, ok bool) {
	t, ok = map[string]BooleanOperation{
		"union":        BooleanOperation_union,
		"difference":   BooleanOperation_difference,
		"intersection": BooleanOperation_intersection,
	}[s]
	return
}

func (b BooleanOperation) String() string {
	return map[BooleanOperation]string{
		BooleanOperation_union:        "union",
		BooleanOperation_difference:   "difference",
		BooleanOperation_intersection: "intersection",
	}[b]
}

type BooleanOperationAttr struct {
	Association Association
	Operation   BooleanOperation
}

func (b *BooleanOperationAttr) GetAssociation() Association {
	return b.Association
}
func (b *BooleanOperationAttr) GetBooleanOperation() BooleanOperation {
	return b.Operation
}

func GetBooleanOperationAttr(component *go3mf.Components) *BooleanOperationAttr {
	for _, a := range component.AnyAttr {
		if a, ok := a.(*BooleanOperationAttr); ok {
			return a
		}
	}
	return nil
}
