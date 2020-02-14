package slices

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestDecode(t *testing.T) {
	sliceStack := &SliceStackResource{ID: 3, BottomZ: 1,
		Slices: []*Slice{
			{
				TopZ:     0,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 1, PID: 1, P1: 2, P2: 3}, {V2: 2, PID: 1, P1: 2, P2: 2}, {V2: 3}, {V2: 0}}}},
			},
			{
				TopZ:     0.1,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 2}, {V2: 1}, {V2: 3}, {V2: 0}}}},
			},
		},
	}
	sliceStackRef := &SliceStackResource{ID: 7, BottomZ: 1.1, Refs: []SliceRef{{SliceStackID: 10, Path: "/2D/2Dmodel.model"}}}
	meshRes := &go3mf.ObjectResource{
		Mesh: new(go3mf.Mesh),
		ID:   8, Name: "Box 1",
		ExtensionAttr: go3mf.ExtensionAttr{ExtensionName: &SliceStackInfo{SliceStackID: 3, SliceResolution: ResolutionLow}},
	}

	want := &go3mf.Model{Path: "/3D/3dmodel.model", Namespaces: []xml.Name{{Space: ExtensionName, Local: "s"}}, Resources: go3mf.Resources{
		Assets: []go3mf.Resource{sliceStack, sliceStackRef}, Objects: []*go3mf.ObjectResource{meshRes},
	}}
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
	<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:s="http://schemas.microsoft.com/3dmanufacturing/slice/2015/07">
		<resources>
			<s:other />
			<s:slicestack id="3" zbottom="1">
				<s:slice ztop="0">
					<s:vertices>
						<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
					</s:vertices>
					<s:polygon startv="0">
						<s:segment v2="1" pid="1" p1="2" p2="3"></s:segment> <s:segment v2="2" pid="1" p1="2" p2="2"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
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
					</vertices>
					<triangles>
					</triangles>
				</mesh>
			</object>
		</resources>
		<build>
		</build>
	</model>`

	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		RegisterExtension(d)
		d.Strict = true
		if err := d.UnmarshalModel([]byte(rootFile), got); err != nil {
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
		go3mf.ParsePropertyError{ResourceID: 3, Element: "slicestack", Name: "zbottom", Value: "a", ModelPath: "/3D/3dmodel.model", Type: go3mf.PropertyOptional},
		go3mf.MissingPropertyError{ResourceID: 3, Element: "slice", ModelPath: "/3D/3dmodel.model", Name: "ztop"},
		go3mf.ParsePropertyError{ResourceID: 3, Element: "vertex", Name: "x", Value: "a", ModelPath: "/3D/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 3, Element: "vertex", Name: "y", Value: "b", ModelPath: "/3D/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 3, Element: "slice", Name: "ztop", Value: "a", ModelPath: "/3D/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 3, Element: "polygon", Name: "startv", Value: "a", ModelPath: "/3D/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 3, Element: "segment", Name: "v2", Value: "a", ModelPath: "/3D/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 3, Element: "sliceref", Name: "slicestackid", Value: "a", ModelPath: "/3D/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 3, Element: "sliceref", ModelPath: "/3D/3dmodel.model", Name: "slicestackid"},
		//go3mf.GenericError{ResourceID: 3, Element: "sliceref", ModelPath: "/3D/3dmodel.model", Message: "non-existent referenced resource"},
		go3mf.GenericError{ResourceID: 3, Element: "slicestack", ModelPath: "/3D/3dmodel.model", Message: "slicestack contains slices and slicerefs"},
		go3mf.MissingPropertyError{ResourceID: 7, Element: "sliceref", ModelPath: "/3D/3dmodel.model", Name: "slicestackid"},
		//go3mf.GenericError{ResourceID: 7, Element: "sliceref", ModelPath: "/3D/3dmodel.model", Message: "non-existent referenced resource"},
		go3mf.ParsePropertyError{ResourceID: 8, Element: "object", ModelPath: "/3D/3dmodel.model", Name: "meshresolution", Value: "invalid", Type: go3mf.PropertyOptional},
		go3mf.ParsePropertyError{ResourceID: 8, Element: "object", ModelPath: "/3D/3dmodel.model", Name: "slicestackid", Value: "a", Type: go3mf.PropertyRequired},
	}
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:s="http://schemas.microsoft.com/3dmanufacturing/slice/2015/07">
		<resources>
			<s:slicestack id="3" zbottom="a">
				<s:slice>
					<s:vertices>
						<s:vertex x="a" y="1.02" /> <s:vertex x="9.03" y="b" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
					</s:vertices>
					<s:polygon startv="50">
						<s:segment v2="1"/>
						<s:segment v2="100"/>
					</s:polygon>
				</s:slice>
				<s:slice ztop="a">
					<s:vertices>
						<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
					</s:vertices>
					<s:polygon startv="a"> 
						<s:segment v2="a"></s:segment> <s:segment v2="1"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
					</s:polygon>
				</s:slice>
				<s:sliceref slicestackid="a" slicepath="/3D/3dmodel.model" />
			</s:slicestack>
			<s:slicestack id="7" zbottom="1.1">
				<s:sliceref slicepath="/2D/2Dmodel.model" />
			</s:slicestack>
			<object id="8" name="Box 1"s:meshresolution="invalid" s:slicestackid="a">
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
		</build>
		</model>
		`

	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		RegisterExtension(d)
		d.Strict = false
		if err := d.UnmarshalModel([]byte(rootFile), got); err != nil {
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

func Test_baseDecoder_Child(t *testing.T) {
	type args struct {
		in0 xml.Name
	}
	tests := []struct {
		name string
		d    *baseDecoder
		args args
		want go3mf.NodeDecoder
	}{
		{"base", new(baseDecoder), args{xml.Name{}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Child(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDecoder.Child() = %v, want %v", got, tt.want)
			}
		})
	}
}
