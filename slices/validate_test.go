package slices

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	rootPath := go3mf.DefaultModelPath
	tests := []struct {
		name  string
		model *go3mf.Model
		want  []error
	}{
		{"extRequired", &go3mf.Model{
			AnyAttr: go3mf.AnyAttr{&SliceStackInfo{SliceStackID: 10}},
			Specs:   map[string]go3mf.Spec{ExtensionSpace: &Extension{IsRequired: true}}, Resources: go3mf.Resources{
				Assets: []go3mf.Asset{
					&SliceStack{ID: 1, Slices: []*Slice{{TopZ: 1}}},
				},
				Objects: []*go3mf.Object{
					{ID: 2, AnyAttr: go3mf.AnyAttr{&SliceStackInfo{
						SliceStackID: 1, SliceResolution: ResolutionLow,
					}}},
				}},
		}, []error{
			fmt.Errorf("Resources@Object#0: %v", specerr.ErrInvalidObject),
		}},
		{"child", &go3mf.Model{Childs: map[string]*go3mf.ChildModel{
			"/other.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&SliceStack{ID: 1},
			}}},
			"/that.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&SliceStack{ID: 2},
			}}},
		}}, []error{
			fmt.Errorf("/other.model@Resources@SliceStack#0: %v", specerr.ErrSlicesAndRefs),
			fmt.Errorf("/that.model@Resources@SliceStack#0: %v", specerr.ErrSlicesAndRefs),
		}},
		{"slicestack", &go3mf.Model{Resources: go3mf.Resources{
			Assets: []go3mf.Asset{&SliceStack{
				ID: 1, BottomZ: 1, Slices: []*Slice{
					{},
					{TopZ: 0.5, Vertices: make([]go3mf.Point2D, 1)},
					{TopZ: 1.5, Vertices: make([]go3mf.Point2D, 2), Polygons: []Polygon{
						{Segments: []Segment{}},
						{Segments: []Segment{{}}},
					}},
					{TopZ: 1.5},
					{TopZ: 1.4},
				},
			}},
		}}, []error{
			fmt.Errorf("Resources@SliceStack#0@Slice#0: %v", &specerr.MissingFieldError{Name: attrZTop}),
			fmt.Errorf("Resources@SliceStack#0@Slice#1: %v", specerr.ErrSliceSmallTopZ),
			fmt.Errorf("Resources@SliceStack#0@Slice#1: %v", specerr.ErrSliceInsufficientVertices),
			fmt.Errorf("Resources@SliceStack#0@Slice#1: %v", specerr.ErrSliceInsufficientPolygons),
			fmt.Errorf("Resources@SliceStack#0@Slice#2@Polygon#0: %v", specerr.ErrSliceInsufficientSegments),
			fmt.Errorf("Resources@SliceStack#0@Slice#3: %v", specerr.ErrSliceNoMonotonic),
			fmt.Errorf("Resources@SliceStack#0@Slice#4: %v", specerr.ErrSliceNoMonotonic),
		}},
		{"sliceref", &go3mf.Model{
			Childs: map[string]*go3mf.ChildModel{
				"/that.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
					&SliceStack{ID: 1, Slices: []*Slice{{TopZ: 1}, {TopZ: 2}}},
					&SliceStack{ID: 2, Refs: []SliceRef{{SliceStackID: 1, Path: rootPath}}},
					&go3mf.BaseMaterials{ID: 3, Materials: []go3mf.Base{{Name: "a", Color: color.RGBA{R: 1}}}},
					&SliceStack{ID: 4, Slices: []*Slice{{TopZ: 1.5}}},
				}}},
			},
			Resources: go3mf.Resources{
				Assets: []go3mf.Asset{
					&SliceStack{ID: 1, Slices: []*Slice{{TopZ: 1}, {TopZ: 2}}},
					&SliceStack{ID: 3, Refs: []SliceRef{
						{},
						{SliceStackID: 1, Path: rootPath},
						{SliceStackID: 1, Path: "/other.model"},
						{SliceStackID: 2, Path: "/that.model"},
						{SliceStackID: 3, Path: "/that.model"},
						{SliceStackID: 1, Path: "/that.model"},
						{SliceStackID: 4, Path: "/that.model"},
					},
					}},
			}}, []error{
			fmt.Errorf("Resources@SliceStack#1@SliceRef#0: %v", &specerr.MissingFieldError{Name: attrSlicePath}),
			fmt.Errorf("Resources@SliceStack#1@SliceRef#0: %v", &specerr.MissingFieldError{Name: attrSliceRefID}),
			fmt.Errorf("Resources@SliceStack#1@SliceRef#1: %v", specerr.ErrSliceRefSamePart),
			fmt.Errorf("Resources@SliceStack#1@SliceRef#2: %v", specerr.ErrMissingResource),
			fmt.Errorf("Resources@SliceStack#1@SliceRef#3: %v", specerr.ErrSliceRefRef),
			fmt.Errorf("Resources@SliceStack#1@SliceRef#4: %v", specerr.ErrNonSliceStack),
			fmt.Errorf("Resources@SliceStack#1@SliceRef#6: %v", specerr.ErrSliceNoMonotonic),
		}},
		{"info", &go3mf.Model{Build: go3mf.Build{Items: []*go3mf.Item{
			{ObjectID: 7},
			{ObjectID: 2},
			{ObjectID: 2, Transform: go3mf.Identity()},
			{ObjectID: 2, Transform: go3mf.Matrix{2, 3, 0, 0, 1, 3, 0, 0, 0, 0, 1, 0, 2, 3, 4, 1}},
			{ObjectID: 2, Transform: go3mf.Matrix{2, 3, 1, 0, 1, 3, 0, 0, 0, 0, 1, 0, 2, 3, 4, 1}},
			{ObjectID: 4},
			{ObjectID: 12},
		}},
			Childs: map[string]*go3mf.ChildModel{
				"/that.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
					&SliceStack{ID: 1, Slices: []*Slice{{TopZ: 1, Vertices: []go3mf.Point2D{{}, {}, {}}, Polygons: []Polygon{
						{StartV: 1, Segments: []Segment{{V2: 2}}},
					}}}},
				}}},
			}, Resources: go3mf.Resources{
				Assets: []go3mf.Asset{
					&SliceStack{ID: 3, Slices: []*Slice{{TopZ: 1, Vertices: []go3mf.Point2D{{}, {}, {}}, Polygons: []Polygon{
						{StartV: 1, Segments: []Segment{{V2: 2}}},
					}}}},
					&go3mf.BaseMaterials{ID: 6, Materials: []go3mf.Base{{Name: "a", Color: color.RGBA{R: 1}}}},
					&SliceStack{ID: 9, Refs: []SliceRef{{SliceStackID: 1, Path: "/that.model"}}},
					&SliceStack{ID: 11, Slices: []*Slice{{TopZ: 1, Vertices: []go3mf.Point2D{{}, {}, {}}, Polygons: []Polygon{
						{StartV: 1, Segments: []Segment{{V2: 2}, {V2: 1}}},
					}}}},
				},
				Objects: []*go3mf.Object{
					{ID: 1, Mesh: &go3mf.Mesh{}, AnyAttr: go3mf.AnyAttr{&SliceStackInfo{
						SliceStackID: 1,
					}}},
					{ID: 2, ObjectType: go3mf.ObjectTypeSupport, Components: []*go3mf.Component{{ObjectID: 1}}, AnyAttr: go3mf.AnyAttr{&SliceStackInfo{
						SliceStackID: 3, SliceResolution: ResolutionLow,
					}}},
					{ID: 4, Components: []*go3mf.Component{{ObjectID: 10, Transform: go3mf.Matrix{2, 3, 0, 0, 1, 3, 0, 0, 0, 0, 2, 0, 2, 3, 4, 1}}},
						AnyAttr: go3mf.AnyAttr{&SliceStackInfo{
							SliceStackID: 0,
						}}},
					{ID: 5, Components: []*go3mf.Component{{ObjectID: 12}}, AnyAttr: go3mf.AnyAttr{&SliceStackInfo{
						SliceStackID: 6,
					}}},
					{ID: 7, Components: []*go3mf.Component{{ObjectID: 1}, {ObjectID: 4}}, AnyAttr: go3mf.AnyAttr{&SliceStackInfo{
						SliceStackID: 9,
					}}},
					{ID: 10, ObjectType: go3mf.ObjectTypeSolidSupport, Components: []*go3mf.Component{{ObjectID: 1}}, AnyAttr: go3mf.AnyAttr{&SliceStackInfo{
						SliceStackID: 3,
					}}},
					{ID: 12, Components: []*go3mf.Component{
						{ObjectID: 7, Transform: go3mf.Matrix{2, 3, 0, 0, 1, 3, 1, 0, 0, 0, 1, 0, 2, 3, 4, 1}},
						{ObjectID: 12},
						{ObjectID: 5},
					}, AnyAttr: go3mf.AnyAttr{&SliceStackInfo{
						SliceStackID: 11,
					}}},
				}}}, []error{
			fmt.Errorf("Resources@Object#0@Mesh: %v", specerr.ErrInsufficientVertices),
			fmt.Errorf("Resources@Object#0@Mesh: %v", specerr.ErrInsufficientTriangles),
			fmt.Errorf("Resources@Object#0: %v", specerr.ErrMissingResource),
			fmt.Errorf("Resources@Object#1: %v", specerr.ErrSliceInvalidTranform),
			fmt.Errorf("Resources@Object#1: %v", specerr.ErrSliceExtRequired),
			fmt.Errorf("Resources@Object#2: %v", &specerr.MissingFieldError{Name: attrSliceRefID}),
			fmt.Errorf("Resources@Object#3: %v", specerr.ErrNonSliceStack),
			fmt.Errorf("Resources@Object#4: %v", specerr.ErrSliceInvalidTranform),
			fmt.Errorf("Resources@Object#4: %v", specerr.ErrSlicePolygonNotClosed),
			fmt.Errorf("Resources@Object#5: %v", specerr.ErrSliceInvalidTranform),
			fmt.Errorf("Resources@Object#5: %v", specerr.ErrSlicePolygonNotClosed),
			fmt.Errorf("Resources@Object#6@Component#1: %v", specerr.ErrRecursion),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.model.Specs) == 0 {
				tt.model.WithExtension(&Extension{})
			}
			got := tt.model.Validate()
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
