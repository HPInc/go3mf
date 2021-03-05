package slices

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	rootPath := go3mf.DefaultModelPath
	tests := []struct {
		name  string
		model *go3mf.Model
		want  []string
	}{
		{"extRequired", &go3mf.Model{
			AnyAttr:    go3mf.AnyAttr{&ObjectAttr{SliceStackID: 10}},
			Extensions: []go3mf.Extension{{Namespace: Namespace, LocalName: "s", IsRequired: true}}, Resources: go3mf.Resources{
				Assets: []go3mf.Asset{
					&SliceStack{ID: 1, Slices: []*Slice{{TopZ: 1}}},
				},
				Objects: []*go3mf.Object{
					{ID: 2, AnyAttr: go3mf.AnyAttr{&ObjectAttr{
						SliceStackID: 1, MeshResolution: ResolutionLow,
					}}},
				}},
		}, []string{
			fmt.Sprintf("Resources@Object#0: %v", errors.ErrInvalidObject),
		}},
		{"child", &go3mf.Model{Childs: map[string]*go3mf.ChildModel{
			"/other.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&SliceStack{ID: 1},
			}}},
			"/that.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&SliceStack{ID: 2},
			}}},
		}}, []string{
			fmt.Sprintf("/other.model@Resources@SliceStack#0: %v", ErrSlicesAndRefs),
			fmt.Sprintf("/that.model@Resources@SliceStack#0: %v", ErrSlicesAndRefs),
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
		}}, []string{
			fmt.Sprintf("Resources@SliceStack#0@Slice#0: %v", &errors.MissingFieldError{Name: attrZTop}),
			fmt.Sprintf("Resources@SliceStack#0@Slice#1: %v", ErrSliceSmallTopZ),
			fmt.Sprintf("Resources@SliceStack#0@Slice#1: %v", ErrSliceInsufficientVertices),
			fmt.Sprintf("Resources@SliceStack#0@Slice#1: %v", ErrSliceInsufficientPolygons),
			fmt.Sprintf("Resources@SliceStack#0@Slice#2@Polygon#0: %v", ErrSliceInsufficientSegments),
			fmt.Sprintf("Resources@SliceStack#0@Slice#3: %v", ErrSliceNoMonotonic),
			fmt.Sprintf("Resources@SliceStack#0@Slice#4: %v", ErrSliceNoMonotonic),
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
			}}, []string{
			fmt.Sprintf("Resources@SliceStack#1@SliceRef#0: %v", &errors.MissingFieldError{Name: attrSlicePath}),
			fmt.Sprintf("Resources@SliceStack#1@SliceRef#0: %v", &errors.MissingFieldError{Name: attrSliceRefID}),
			fmt.Sprintf("Resources@SliceStack#1@SliceRef#1: %v", ErrSliceRefSamePart),
			fmt.Sprintf("Resources@SliceStack#1@SliceRef#2: %v", errors.ErrMissingResource),
			fmt.Sprintf("Resources@SliceStack#1@SliceRef#3: %v", ErrSliceRefRef),
			fmt.Sprintf("Resources@SliceStack#1@SliceRef#4: %v", ErrNonSliceStack),
			fmt.Sprintf("Resources@SliceStack#1@SliceRef#6: %v", ErrSliceNoMonotonic),
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
					{ID: 1, Mesh: &go3mf.Mesh{}, AnyAttr: go3mf.AnyAttr{&ObjectAttr{
						SliceStackID: 1,
					}}},
					{ID: 2, Type: go3mf.ObjectTypeSupport, Components: []*go3mf.Component{{ObjectID: 1}}, AnyAttr: go3mf.AnyAttr{&ObjectAttr{
						SliceStackID: 3, MeshResolution: ResolutionLow,
					}}},
					{ID: 4, Components: []*go3mf.Component{{ObjectID: 10, Transform: go3mf.Matrix{2, 3, 0, 0, 1, 3, 0, 0, 0, 0, 2, 0, 2, 3, 4, 1}}},
						AnyAttr: go3mf.AnyAttr{&ObjectAttr{
							SliceStackID: 0,
						}}},
					{ID: 5, Components: []*go3mf.Component{{ObjectID: 12}}, AnyAttr: go3mf.AnyAttr{&ObjectAttr{
						SliceStackID: 6,
					}}},
					{ID: 7, Components: []*go3mf.Component{{ObjectID: 1}, {ObjectID: 4}}, AnyAttr: go3mf.AnyAttr{&ObjectAttr{
						SliceStackID: 9,
					}}},
					{ID: 10, Type: go3mf.ObjectTypeSolidSupport, Components: []*go3mf.Component{{ObjectID: 1}}, AnyAttr: go3mf.AnyAttr{&ObjectAttr{
						SliceStackID: 3,
					}}},
					{ID: 12, Components: []*go3mf.Component{
						{ObjectID: 7, Transform: go3mf.Matrix{2, 3, 0, 0, 1, 3, 1, 0, 0, 0, 1, 0, 2, 3, 4, 1}},
						{ObjectID: 12},
						{ObjectID: 5},
					}, AnyAttr: go3mf.AnyAttr{&ObjectAttr{
						SliceStackID: 11,
					}}},
				}}}, []string{
			fmt.Sprintf("Resources@Object#0@Mesh: %v", errors.ErrInsufficientVertices),
			fmt.Sprintf("Resources@Object#0@Mesh: %v", errors.ErrInsufficientTriangles),
			fmt.Sprintf("Resources@Object#0: %v", errors.ErrMissingResource),
			fmt.Sprintf("Resources@Object#1: %v", ErrSliceInvalidTranform),
			fmt.Sprintf("Resources@Object#1: %v", ErrSliceExtRequired),
			fmt.Sprintf("Resources@Object#2: %v", &errors.MissingFieldError{Name: attrSliceRefID}),
			fmt.Sprintf("Resources@Object#3: %v", ErrNonSliceStack),
			fmt.Sprintf("Resources@Object#4: %v", ErrSliceInvalidTranform),
			fmt.Sprintf("Resources@Object#4: %v", ErrSlicePolygonNotClosed),
			fmt.Sprintf("Resources@Object#5: %v", ErrSliceInvalidTranform),
			fmt.Sprintf("Resources@Object#5: %v", ErrSlicePolygonNotClosed),
			fmt.Sprintf("Resources@Object#6@Component#1: %v", errors.ErrRecursion),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.model.Extensions) == 0 {
				tt.model.Extensions = []go3mf.Extension{DefaultExtension}
			}
			err := tt.model.Validate()
			if err == nil {
				t.Fatal("error expected")
			}
			var errs []string
			for _, err := range err.(*errors.List).Errors {
				errs = append(errs, err.Error())
			}
			if diff := deep.Equal(errs, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
