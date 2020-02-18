package materials

import (
	"encoding/xml"
	"image/color"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestMarshalModel(t *testing.T) {
	baseTexture := &Texture2DResource{ID: 6, Path: "/3D/Texture/msLogo.png", ContentType: TextureTypePNG, TileStyleU: TileWrap, TileStyleV: TileMirror, Filter: TextureFilterAuto}
	colorGroup := &ColorGroupResource{ID: 1, Colors: []color.RGBA{{R: 255, G: 255, B: 255, A: 255}, {R: 0, G: 0, B: 0, A: 255}, {R: 26, G: 181, B: 103, A: 255}, {R: 223, G: 4, B: 90, A: 255}}}
	texGroup := &Texture2DGroupResource{ID: 2, TextureID: 6, Coords: []TextureCoord{{0.3, 0.5}, {0.3, 0.8}, {0.5, 0.8}, {0.5, 0.5}}}
	compositeGroup := &CompositeMaterialsResource{ID: 4, MaterialID: 5, Indices: []uint32{1, 2}, Composites: []Composite{{Values: []float32{0.5, 0.5}}, {Values: []float32{0.2, 0.8}}}}
	multiGroup := &MultiPropertiesResource{ID: 9, BlendMethods: []BlendMethod{BlendMultiply}, PIDs: []uint32{5, 2}, Multis: []Multi{{PIndex: []uint32{0, 0}}, {PIndex: []uint32{1, 0}}, {PIndex: []uint32{2, 3}}}}
	m := &go3mf.Model{Path: "/3D/3dmodel.model", Namespaces: []xml.Name{{Space: ExtensionName, Local: "m"}}}
	m.Resources.Assets = append(m.Resources.Assets, baseTexture, colorGroup, texGroup, compositeGroup, multiGroup)

	t.Run("base", func(t *testing.T) {
		b, err := go3mf.MarshalModel(m)
		if err != nil {
			t.Errorf("materials.MarshalModel() error = %v", err)
			return
		}
		d := go3mf.NewDecoder(nil, 0)
		RegisterExtension(d)
		newModel := new(go3mf.Model)
		newModel.Path = m.Path
		if err := d.UnmarshalModel(b, newModel); err != nil {
			t.Errorf("materials.MarshalModel() error decoding = %v, s = %s", err, string(b))
			return
		}
		if diff := deep.Equal(m, newModel); diff != nil {
			t.Errorf("materials.MarshalModel() = %v, s = %s", diff, string(b))
		}
	})
}
