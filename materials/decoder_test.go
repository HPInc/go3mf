package materials

import (
	"context"
	"image/color"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestDecode(t *testing.T) {
	baseTexture := &Texture2DResource{ID: 6, ModelPath: "/3d/3dmodel.model", Path: "/3D/Texture/msLogo.png", ContentType: TextureTypePNG, TileStyleU: TileWrap, TileStyleV: TileMirror, Filter: TextureFilterAuto}
	meshRes := &go3mf.Mesh{
		ObjectResource: go3mf.ObjectResource{
			ID: 8, Name: "Box 1", ModelPath: "/3d/3dmodel.model", Thumbnail: "/a.png", DefaultPropertyID: 5, PartNumber: "11111111-1111-1111-1111-111111111111"},
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
		{NodeIndices: [3]uint32{3, 2, 1}, Resource: 5},
		{NodeIndices: [3]uint32{1, 0, 3}, Resource: 5},
		{NodeIndices: [3]uint32{4, 5, 6}, Resource: 5, ResourceIndices: [3]uint32{1, 1, 1}},
		{NodeIndices: [3]uint32{6, 7, 4}, Resource: 5, ResourceIndices: [3]uint32{1, 1, 1}},
		{NodeIndices: [3]uint32{0, 1, 5}, Resource: 2, ResourceIndices: [3]uint32{0, 1, 2}},
		{NodeIndices: [3]uint32{5, 4, 0}, Resource: 2, ResourceIndices: [3]uint32{3, 0, 2}},
		{NodeIndices: [3]uint32{1, 2, 6}, Resource: 1, ResourceIndices: [3]uint32{0, 1, 2}},
		{NodeIndices: [3]uint32{6, 5, 1}, Resource: 1, ResourceIndices: [3]uint32{2, 1, 3}},
		{NodeIndices: [3]uint32{2, 3, 7}, Resource: 5},
		{NodeIndices: [3]uint32{7, 6, 2}, Resource: 5},
		{NodeIndices: [3]uint32{3, 0, 4}, Resource: 5},
		{NodeIndices: [3]uint32{4, 7, 3}, Resource: 5},
	}...)

	want := &go3mf.Model{Path: "/3d/3dmodel.model"}
	colorGroup := &ColorGroupResource{ID: 1, ModelPath: "/3d/3dmodel.model", Colors: []color.RGBA{{R: 255, G: 255, B: 255, A: 255}, {R: 0, G: 0, B: 0, A: 255}, {R: 26, G: 181, B: 103, A: 255}, {R: 223, G: 4, B: 90, A: 255}}}
	texGroup := &Texture2DGroupResource{ID: 2, ModelPath: "/3d/3dmodel.model", TextureID: 6, Coords: []TextureCoord{{0.3, 0.5}, {0.3, 0.8}, {0.5, 0.8}, {0.5, 0.5}}}
	compositeGroup := &CompositeMaterialsResource{ID: 4, ModelPath: "/3d/3dmodel.model", MaterialID: 5, Indices: []uint32{1, 2}, Composites: []Composite{{Values: []float32{0.5, 0.5}}, {Values: []float32{0.2, 0.8}}}}
	multiGroup := &MultiPropertiesResource{ID: 9, ModelPath: "/3d/3dmodel.model", BlendMethods: []BlendMethod{BlendMultiply}, Resources: []uint32{5, 2}, Multis: []Multi{{ResourceIndices: []uint32{0, 0}}, {ResourceIndices: []uint32{1, 0}}, {ResourceIndices: []uint32{2, 3}}}}
	want.Resources = append(want.Resources, baseTexture, colorGroup, texGroup, compositeGroup, multiGroup, meshRes)
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	rootFile := `
	<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:m="http://schemas.microsoft.com/3dmanufacturing/material/2015/02">
		<resources>
			<m:texture2d id="6" path="/3D/Texture/msLogo.png" contenttype="image/png" tilestyleu="wrap" tilestylev="mirror" filter="auto" />
			<m:colorgroup id="1">
				<m:color color="#FFFFFF" /> <m:color color="#000000" /> <m:color color="#1AB567" /> <m:color color="#DF045A" />
			</m:colorgroup>
			<m:texture2dgroup id="2" texid="6">
				<m:tex2coord u="0.3" v="0.5" /> <m:tex2coord u="0.3" v="0.8" />	<m:tex2coord u="0.5" v="0.8" />	<m:tex2coord u="0.5" v="0.5" />
			</m:texture2dgroup>
			<m:compositematerials id="4" matid="5" matindices="1 2">
				<m:composite values="0.5 0.5"/>
				<m:composite values="0.2 0.8"/>
			</m:compositematerials>
			<m:multiproperties id="9" pids="5 2" blendmethods="multiply">
				<m:multi pindices="0 0" />
				<m:multi pindices="1 0" />
				<m:multi pindices="2 3" />
			</m:multiproperties>
			<object id="8" name="Box 1" pid="5" pindex="0" thumbnail="/a.png" partnumber="11111111-1111-1111-1111-111111111111" type="model">
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
						<triangle v1="4" v2="5" v3="6" p1="1" />
						<triangle v1="6" v2="7" v3="4" pid="5" p1="1" />
						<triangle v1="0" v2="1" v3="5" pid="2" p1="0" p2="1" p3="2"/>
						<triangle v1="5" v2="4" v3="0" pid="2" p1="3" p2="0" p3="2"/>
						<triangle v1="1" v2="2" v3="6" pid="1" p1="0" p2="1" p3="2"/>
						<triangle v1="6" v2="5" v3="1" pid="1" p1="2" p2="1" p3="3"/>
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
	</model>`
	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		RegisterExtension(d)
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
		go3mf.ParsePropertyError{ResourceID: 0, Element: "texture2d", Name: "id", Value: "b", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "texture2d", ModelPath: "/3d/3dmodel.model", Name: "path"},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "texture2d", ModelPath: "/3d/3dmodel.model", Name: "id"},
		go3mf.ParsePropertyError{ResourceID: 1, Element: "color", Name: "color", Value: "#FFFFF", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 2, Element: "tex2coord", Name: "u", Value: "b", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 2, Element: "tex2coord", Name: "v", Value: "c", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 4, Element: "compositematerials", ModelPath: "/3d/3dmodel.model", Name: "matid"},
		go3mf.MissingPropertyError{ResourceID: 4, Element: "compositematerials", ModelPath: "/3d/3dmodel.model", Name: "matindices"},
		go3mf.MissingPropertyError{ResourceID: 4, Element: "composite", ModelPath: "/3d/3dmodel.model", Name: "values"},
		go3mf.ParsePropertyError{ResourceID: 4, Element: "composite", Name: "values", Value: "a", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 9, Element: "multiproperties", ModelPath: "/3d/3dmodel.model", Name: "pids", Value: "a", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 9, Element: "multi", ModelPath: "/3d/3dmodel.model", Name: "pindices"},
		go3mf.MissingPropertyError{ResourceID: 19, Element: "multiproperties", ModelPath: "/3d/3dmodel.model", Name: "pids"},
	}
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	rootFile := `
	<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:m="http://schemas.microsoft.com/3dmanufacturing/material/2015/02">
		<resources>
			<m:texture2d id="6" qm:mq="other" path="/3D/Texture/msLogo.png" contenttype="image/png" tilestyleu="wrap" tilestylev="mirror" filter="auto" />
			<m:texture2d id="b" contenttype="image/png" tilestyleu="wrap" tilestylev="mirror" filter="auto" />
			<m:colorgroup id="1">
				<m:color color="#FFFFF" /> <m:color color="#000000" /> <m:color color="#1AB567" /> <m:color color="#DF045A" />
			</m:colorgroup>
			<m:texture2dgroup qm:mq="other" id="2" texid="6">
				<m:tex2coord qm:mq="other" u="b" v="0.5" /> <m:tex2coord u="0.3" v="c" />	<m:tex2coord u="0.5" v="0.8" />	<m:tex2coord u="0.5" v="0.5" />
			</m:texture2dgroup>
			<m:compositematerials id="4" qm:mq="other">
				<m:composite/>
				<m:composite values="a 0.8"/>
			</m:compositematerials>
			<m:multiproperties id="9" qm:mq="other" pids="a 2">
				<m:multi />
			</m:multiproperties>
			<m:multiproperties id="19" />
			<object id="8" name="Box 1" pid="5" pindex="0" type="model">
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
						<triangle v1="4" v2="5" v3="6" p1="1" />
						<triangle v1="6" v2="7" v3="4" pid="5" p1="1" />
						<triangle v1="0" v2="1" v3="5" pid="2" p1="0" p2="1" p3="2"/>
						<triangle v1="5" v2="4" v3="0" pid="2" p1="3" p2="0" p3="2"/>
						<triangle v1="1" v2="2" v3="6" pid="1" p1="0" p2="1" p3="2"/>
						<triangle v1="6" v2="5" v3="1" pid="1" p1="2" p2="1" p3="3"/>
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
	</model>`
	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		RegisterExtension(d)
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

func Test_fileFilter(t *testing.T) {
	type args struct {
		relType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"accepted", args{RelTypeTexture3D}, true},
		{"rejected", args{"other"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileFilter(tt.args.relType, true); got != tt.want {
				t.Errorf("fileFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
