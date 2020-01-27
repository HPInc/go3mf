package slices

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestDecode(t *testing.T) {
	otherSlices := SliceStack{
		BottomZ: 2,
		Slices: []*Slice{
			{
				TopZ:     1.2,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 1, 2, 3, 0}},
			},
		},
	}
	sliceStack := &SliceStackResource{ID: 3, ModelPath: "/3d/3dmodel.model", Stack: SliceStack{
		BottomZ: 1,
		Slices: []*Slice{
			{
				TopZ:     0,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 1, 2, 3, 0}},
			},
			{
				TopZ:     0.1,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 2, 1, 3, 0}},
			},
		},
	}}
	sliceStackRef := &SliceStackResource{ID: 7, ModelPath: "/3d/3dmodel.model", Stack: SliceStack{BottomZ: 1.1, Refs: []SliceRef{{SliceStackID: 10, Path: "/2D/2Dmodel.model"}}}}
	meshRes := &go3mf.Mesh{
		ObjectResource: go3mf.ObjectResource{
			ID: 8, Name: "Box 1", ModelPath: "/3d/3dmodel.model",
			Extensions: map[string]interface{}{ExtensionName: &ObjectAttr{SliceStackID: 3, SliceResolution: ResolutionLow}}},
	}
	meshRes.Nodes = append(meshRes.Nodes, []go3mf.Point3D{
		{0, 0, 0},
		{100, 0, 0},
		{100, 100, 0},
		{0, 100, 0},
		{0, 0, 100},
		{100, 0, 100},
		{100, 100, 100},
		{0, 100, 100},
	}...)
	meshRes.Faces = append(meshRes.Faces, []go3mf.Face{
		{NodeIndices: [3]uint32{3, 2, 1}},
		{NodeIndices: [3]uint32{1, 0, 3}},
		{NodeIndices: [3]uint32{4, 5, 6}},
		{NodeIndices: [3]uint32{6, 7, 4}},
		{NodeIndices: [3]uint32{0, 1, 5}},
		{NodeIndices: [3]uint32{5, 4, 0}},
		{NodeIndices: [3]uint32{1, 2, 6}},
		{NodeIndices: [3]uint32{6, 5, 1}},
		{NodeIndices: [3]uint32{2, 3, 7}},
		{NodeIndices: [3]uint32{7, 6, 2}},
		{NodeIndices: [3]uint32{3, 0, 4}},
		{NodeIndices: [3]uint32{4, 7, 3}},
	}...)

	want := &go3mf.Model{Path: "/3d/3dmodel.model"}
	want.Resources = append(want.Resources, &SliceStackResource{ID: 10, ModelPath: "/2D/2Dmodel.model", Stack: otherSlices})
	want.Resources = append(want.Resources, sliceStack, sliceStackRef, meshRes)
	want.BuildItems = append(want.BuildItems, &go3mf.BuildItem{Object: meshRes})
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	got.Resources = append(got.Resources, &SliceStackResource{ID: 10, ModelPath: "/2D/2Dmodel.model", Stack: otherSlices})
	rootFile := `
	<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:s="http://schemas.microsoft.com/3dmanufacturing/slice/2015/07">
		<resources>
			<s:slicestack id="3" zbottom="1">
				<s:slice ztop="0">
					<s:vertices>
						<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
					</s:vertices>
					<s:polygon startv="0">
						<s:segment v2="1"></s:segment> <s:segment v2="2"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
					</s:polygon>
				</s:slice>
				<s:slice ztop="0.1">
					<s:vertices>
						<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
					</s:vertices>
					<s:polygon startv="0"> 
						<s:segment v2="2"></s:segment> <s:segment v2="1"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
					</s:polygon>
				</s:slice>
			</s:slicestack>
			<s:slicestack id="7" zbottom="1.1">
				<s:sliceref slicestackid="10" slicepath="/2D/2Dmodel.model" />
			</s:slicestack>
			<object id="8" name="Box 1" s:meshresolution="lowres" s:slicestackid="3" type="model">
				<mesh>
					<vertices>
						<vertex x="0" y="0" z="0" />
						<vertex x="100.00000" y="0" z="0" />
						<vertex x="100.00000" y="100.00000" z="0" />
						<vertex x="0" y="100.00000" z="0" />
						<vertex x="0" y="0" z="100.00000" />
						<vertex x="100.00000" y="0" z="100.00000" />
						<vertex x="100.00000" y="100.00000" z="100.00000" />
						<vertex x="0" y="100.00000" z="100.00000" />
					</vertices>
					<triangles>
						<triangle v1="3" v2="2" v3="1" />
						<triangle v1="1" v2="0" v3="3" />
						<triangle v1="4" v2="5" v3="6" />
						<triangle v1="6" v2="7" v3="4" />
						<triangle v1="0" v2="1" v3="5" />
						<triangle v1="5" v2="4" v3="0" />
						<triangle v1="1" v2="2" v3="6" />
						<triangle v1="6" v2="5" v3="1" />
						<triangle v1="2" v2="3" v3="7" />
						<triangle v1="7" v2="6" v3="2" />
						<triangle v1="3" v2="0" v3="4" />
						<triangle v1="4" v2="7" v3="3" />
					</triangles>
				</mesh>
			</object>
		</resources>
		<build>
			<item objectid="8"/>
		</build>
	</model>`

	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		d.Strict = true
		if err := d.DecodeRawModel(context.Background(), got, rootFile); err != nil {
			t.Errorf("DecodeRawModel() unexpected error = %v", err)
			return
		}
		deep.CompareUnexportedFields = true
		deep.MaxDepth = 20
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("DecodeRawModell() = %v", diff)
			return
		}
	})
}

func TestDecode_warns(t *testing.T) {
	want := []error{
		go3mf.MissingPropertyError{ResourceID: 3, Element: "slice", ModelPath: "/3d/3dmodel.model", Name: "ztop"},
		go3mf.ParsePropertyError{ResourceID: 3, Element: "vertex", Name: "x", Value: "a", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 3, Element: "vertex", Name: "y", Value: "b", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.GenericError{ResourceID: 3, Element: "polygon", ModelPath: "/3d/3dmodel.model", Message: "invalid slice segment index"},
		go3mf.GenericError{ResourceID: 3, Element: "segment", ModelPath: "/3d/3dmodel.model", Message: "invalid slice segment index"},
		go3mf.GenericError{ResourceID: 3, Element: "polygon", ModelPath: "/3d/3dmodel.model", Message: "a closed slice polygon is actually a line"},
		go3mf.GenericError{ResourceID: 3, Element: "sliceref", ModelPath: "/3d/3dmodel.model", Message: "a slicepath is invalid"},
		go3mf.GenericError{ResourceID: 3, Element: "sliceref", ModelPath: "/3d/3dmodel.model", Message: "non-existent referenced resource"},
		go3mf.GenericError{ResourceID: 3, Element: "slicestack", ModelPath: "/3d/3dmodel.model", Message: "slicestack contains slices and slicerefs"},
		go3mf.MissingPropertyError{ResourceID: 7, Element: "sliceref", ModelPath: "/3d/3dmodel.model", Name: "slicestackid"},
		go3mf.GenericError{ResourceID: 7, Element: "sliceref", ModelPath: "/3d/3dmodel.model", Message: "non-existent referenced resource"},
		go3mf.ParsePropertyError{ResourceID: 8, Element: "object", ModelPath: "/3d/3dmodel.model", Name: "meshresolution", Value: "invalid", Type: go3mf.PropertyOptional},
	}
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:s="http://schemas.microsoft.com/3dmanufacturing/slice/2015/07">
		<resources>
			<s:slicestack id="3" zbottom="1">
				<s:slice>
					<s:vertices>
						<s:vertex x="a" y="1.02" /> <s:vertex x="9.03" y="b" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
					</s:vertices>
					<s:polygon startv="50">
						<s:segment v2="1"/>
						<s:segment v2="100"/>
					</s:polygon>
				</s:slice>
				<s:slice ztop="0.1">
					<s:vertices>
						<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
					</s:vertices>
					<s:polygon startv="0"> 
						<s:segment v2="2"></s:segment> <s:segment v2="1"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
					</s:polygon>
				</s:slice>
				<s:sliceref slicestackid="10" slicepath="/3d/3dmodel.model" />
			</s:slicestack>
			<s:slicestack id="7" zbottom="1.1">
				<s:sliceref slicepath="/2D/2Dmodel.model" />
			</s:slicestack>
			<object id="8" name="Box 1"s:meshresolution="invalid" s:slicestackid="3">
				<mesh>
					<vertices>
						<vertex x="0" y="0" z="0" />
						<vertex x="100.00000" y="0" z="0" />
						<vertex x="100.00000" y="100.00000" z="0" />
						<vertex x="0" y="100.00000" z="0" />
						<vertex x="0" y="0" z="100.00000" />
						<vertex x="100.00000" y="0" z="100.00000" />
						<vertex x="100.00000" y="100.00000" z="100.00000" />
						<vertex x="0" y="100.00000" z="100.00000" />
					</vertices>
					<triangles>
						<triangle v1="2" v2="3" v3="1" />
						<triangle v1="3" v2="2" v3="1" />
						<triangle v1="3" v2="2" v3="1" />
						<triangle v1="1" v2="0" v3="3" />
						<triangle v1="4" v2="5" v3="6" />
						<triangle v1="6" v2="7" v3="4" />
						<triangle v1="0" v2="1" v3="5" />
						<triangle v1="5" v2="4" v3="0" />
						<triangle v1="1" v2="2" v3="6" />
						<triangle v1="6" v2="5" v3="1" />
						<triangle v1="2" v2="3" v3="7" />
						<triangle v1="7" v2="6" v3="2" />
						<triangle v1="3" v2="0" v3="4" />
						<triangle v1="4" v2="7" v3="3" />
					</triangles>
				</mesh>
			</object>
		</resources>
		<build>
			<item objectid="8" p:path="/3d/other.model"/>
		</build>
		</model>
		`

	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		d.Strict = false
		if err := d.DecodeRawModel(context.Background(), got, rootFile); err != nil {
			t.Errorf("DecodeRawModel_warn() unexpected error = %v", err)
			return
		}
		deep.MaxDiff = 1
		if diff := deep.Equal(d.Warnings, want); diff != nil {
			t.Errorf("DecodeRawModel_warn() = %v", diff)
			return
		}
	})
}
