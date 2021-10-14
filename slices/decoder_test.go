// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package slices

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/hpinc/go3mf"
	specerr "github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
)

func TestDecode(t *testing.T) {
	sliceStack := &SliceStack{ID: 3, BottomZ: 1,
		Slices: []Slice{
			{
				TopZ:     0,
				Vertices: Vertices{Vertex: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 1, PID: 1, P1: 2, P2: 3}, {V2: 2, PID: 1, P1: 2, P2: 2}, {V2: 3}, {V2: 0}}}},
			},
			{
				TopZ:     0.1,
				Vertices: Vertices{Vertex: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 2}, {V2: 1}, {V2: 3}, {V2: 0}}}},
			},
		},
	}
	sliceStackRef := &SliceStack{ID: 7, BottomZ: 1.1, Refs: []SliceRef{{SliceStackID: 10, Path: "/2D/2Dmodel.model"}}}
	meshRes := &go3mf.Object{
		Mesh: new(go3mf.Mesh),
		ID:   8, Name: "Box 1",
		AnyAttr: spec.AnyAttr{&ObjectAttr{SliceStackID: 3, MeshResolution: ResolutionLow}},
	}

	want := &go3mf.Model{
		Path:       "/3D/3dmodel.model",
		Extensions: []go3mf.Extension{DefaultExtension},
		Resources: go3mf.Resources{
			Assets: []go3mf.Asset{sliceStack, sliceStackRef}, Objects: []*go3mf.Object{meshRes},
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
		if err := go3mf.UnmarshalModel([]byte(rootFile), got); err != nil {
			t.Errorf("DecodeRawModel() unexpected error = %v", err)
			return
		}
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("DecodeRawModell() = %v", diff)
			return
		}
	})
}

func TestDecode_warns(t *testing.T) {
	want := []string{
		fmt.Sprintf("go3mf: XPath: /model/resources/slicestack[0]: %v", specerr.NewParseAttrError("id", true)),
		fmt.Sprintf("go3mf: XPath: /model/resources/slicestack[0]: %v", specerr.NewParseAttrError("zbottom", false)),
		fmt.Sprintf("go3mf: XPath: /model/resources/slicestack[0]/slice[0]/vertices/vertex[0]: %v", specerr.NewParseAttrError("x", true)),
		fmt.Sprintf("go3mf: XPath: /model/resources/slicestack[0]/slice[0]/vertices/vertex[1]: %v", specerr.NewParseAttrError("y", true)),
		fmt.Sprintf("go3mf: XPath: /model/resources/slicestack[0]/slice[1]: %v", specerr.NewParseAttrError("ztop", true)),
		fmt.Sprintf("go3mf: XPath: /model/resources/slicestack[0]/slice[1]/polygon[0]: %v", specerr.NewParseAttrError("startv", true)),
		fmt.Sprintf("go3mf: XPath: /model/resources/slicestack[0]/slice[1]/polygon[0]/segment[1]: %v", specerr.NewParseAttrError("v2", true)),
		fmt.Sprintf("go3mf: XPath: /model/resources/slicestack[0]/sliceref[0]: %v", specerr.NewParseAttrError("slicestackid", true)),
		fmt.Sprintf("go3mf: XPath: /model/resources/object[0]: %v", specerr.NewParseAttrError("meshresolution", false)),
		fmt.Sprintf("go3mf: XPath: /model/resources/object[0]: %v", specerr.NewParseAttrError("slicestackid", true)),
	}
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:s="http://schemas.microsoft.com/3dmanufacturing/slice/2015/07">
		<resources>
			<s:slicestack id="a" zbottom="a">
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
						<s:segment v2="1"></s:segment> <s:segment v2="a"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
					</s:polygon>
				</s:slice>
				<s:sliceref slicestackid="a" slicepath="/3D/3dmodel.model" />
			</s:slicestack>
			<s:slicestack id="7" zbottom="1.1">
				<s:sliceref slicepath="/2D/2Dmodel.model" />
			</s:slicestack>
			<object id="8" name="Box 1" s:meshresolution="invalid" s:slicestackid="a">
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
		err := go3mf.UnmarshalModel([]byte(rootFile), got)
		if err == nil {
			t.Fatal("error expected")
		}
		var errs []string
		for _, err := range err.(*specerr.List).Errors {
			errs = append(errs, err.Error())
		}
		if diff := deep.Equal(errs, want); diff != nil {
			t.Errorf("UnmarshalModel_warn() = %v", diff)
			return
		}
	})
}
