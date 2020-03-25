package materials

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		model *go3mf.Model
		want  []error
	}{
		{"child", &go3mf.Model{Childs: map[string]*go3mf.ChildModel{
			"/other.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&ColorGroup{ID: 1},
			}}},
			"/that.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&MultiProperties{ID: 2},
			}}},
		}}, []error{
			fmt.Errorf("/other.model@Resources@ColorGroup#0: %v", specerr.ErrEmptyResourceProps),
			fmt.Errorf("/that.model@Resources@MultiProperties#0: %v", &specerr.MissingFieldError{Name: attrPIDs}),
			fmt.Errorf("/that.model@Resources@MultiProperties#0: %v", specerr.ErrMultiBlend),
			fmt.Errorf("/that.model@Resources@MultiProperties#0: %v", specerr.ErrEmptyResourceProps),
		}},
		{"multi", &go3mf.Model{
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&MultiProperties{ID: 4},
				&MultiProperties{ID: 5, Multis: []Multi{{PIndices: []uint32{}}}, PIDs: []uint32{4, 100}},
				&go3mf.BaseMaterials{ID: 1, Materials: []go3mf.Base{
					{Name: "a", Color: color.RGBA{R: 1}},
					{Name: "b", Color: color.RGBA{G: 1}},
				}},
				&ColorGroup{ID: 6, Colors: []color.RGBA{{R: 1}, {R: 2, G: 3, B: 4, A: 5}}},
				&CompositeMaterials{ID: 3, MaterialID: 1, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&MultiProperties{ID: 2, Multis: []Multi{{PIndices: []uint32{1, 0}}}, PIDs: []uint32{1, 6}},
				&MultiProperties{ID: 7, Multis: []Multi{{PIndices: []uint32{1, 3}}}, PIDs: []uint32{1, 6}},
				&MultiProperties{ID: 8, Multis: []Multi{{PIndices: []uint32{}}}, PIDs: []uint32{6, 1, 6}},
				&MultiProperties{ID: 9, Multis: []Multi{{PIndices: []uint32{}}}, PIDs: []uint32{1, 3}},
			}},
		}, []error{
			fmt.Errorf("Resources@MultiProperties#0: %v", &specerr.MissingFieldError{Name: attrPIDs}),
			fmt.Errorf("Resources@MultiProperties#0: %v", specerr.ErrMultiBlend),
			fmt.Errorf("Resources@MultiProperties#0: %v", specerr.ErrEmptyResourceProps),
			fmt.Errorf("Resources@MultiProperties#1: %v", specerr.ErrMultiRefMulti),
			fmt.Errorf("Resources@MultiProperties#1: %v", specerr.ErrMissingResource),
			fmt.Errorf("Resources@MultiProperties#6@Multi#0: %v", specerr.ErrIndexOutOfBounds),
			fmt.Errorf("Resources@MultiProperties#7: %v", specerr.ErrMaterialMulti),
			fmt.Errorf("Resources@MultiProperties#7: %v", specerr.ErrMultiColors),
			fmt.Errorf("Resources@MultiProperties#8: %v", specerr.ErrMaterialMulti),
		}},
		{"missingTextPart", &go3mf.Model{
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&Texture2D{ID: 1},
				&Texture2D{ID: 2, ContentType: TextureTypePNG, Path: "/a.png"},
			}},
		}, []error{
			fmt.Errorf("Resources@Texture2D#0: %v", &specerr.MissingFieldError{Name: attrPath}),
			fmt.Errorf("Resources@Texture2D#0: %v", &specerr.MissingFieldError{Name: attrContentType}),
			fmt.Errorf("Resources@Texture2D#1: %v", specerr.ErrMissingTexturePart),
		}},
		{"textureGroup", &go3mf.Model{
			Attachments: []go3mf.Attachment{{Path: "/a.png"}},
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&Texture2D{ID: 1, ContentType: TextureTypePNG, Path: "/a.png"},
				&Texture2DGroup{ID: 2},
				&Texture2DGroup{ID: 3, TextureID: 1, Coords: []TextureCoord{{}}},
				&Texture2DGroup{ID: 4, TextureID: 2, Coords: []TextureCoord{{}}},
				&Texture2DGroup{ID: 5, TextureID: 100, Coords: []TextureCoord{{}}},
			}},
		}, []error{
			fmt.Errorf("Resources@Texture2DGroup#1: %v", &specerr.MissingFieldError{Name: attrTexID}),
			fmt.Errorf("Resources@Texture2DGroup#1: %v", specerr.ErrEmptyResourceProps),
			fmt.Errorf("Resources@Texture2DGroup#3: %v", specerr.ErrTextureReference),
			fmt.Errorf("Resources@Texture2DGroup#4: %v", specerr.ErrTextureReference),
		}},
		{"colorGroup", &go3mf.Model{
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&ColorGroup{ID: 1},
				&ColorGroup{ID: 2, Colors: []color.RGBA{{R: 1}, {R: 2, G: 3, B: 4, A: 5}}},
				&ColorGroup{ID: 3, Colors: []color.RGBA{{R: 1}, {}}},
			}},
		}, []error{
			fmt.Errorf("Resources@ColorGroup#0: %v", specerr.ErrEmptyResourceProps),
			fmt.Errorf("Resources@ColorGroup#2@RGBA#1: %v", &specerr.MissingFieldError{Name: attrColor}),
		}},
		{"composite", &go3mf.Model{
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&go3mf.BaseMaterials{ID: 1, Materials: []go3mf.Base{
					{Name: "a", Color: color.RGBA{R: 1}},
					{Name: "b", Color: color.RGBA{G: 1}},
				}},
				&CompositeMaterials{ID: 2},
				&CompositeMaterials{ID: 3, MaterialID: 1, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&CompositeMaterials{ID: 4, MaterialID: 1, Indices: []uint32{100, 100}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&CompositeMaterials{ID: 5, MaterialID: 2, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&CompositeMaterials{ID: 6, MaterialID: 100, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
			}}}, []error{
			fmt.Errorf("Resources@CompositeMaterials#1: %v", &specerr.MissingFieldError{Name: attrMatID}),
			fmt.Errorf("Resources@CompositeMaterials#1: %v", &specerr.MissingFieldError{Name: attrMatIndices}),
			fmt.Errorf("Resources@CompositeMaterials#1: %v", specerr.ErrEmptyResourceProps),
			fmt.Errorf("Resources@CompositeMaterials#3: %v", specerr.ErrIndexOutOfBounds),
			fmt.Errorf("Resources@CompositeMaterials#4: %v", specerr.ErrCompositeBase),
			fmt.Errorf("Resources@CompositeMaterials#5: %v", specerr.ErrMissingResource),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.WithSpec(&Spec{})
			got := tt.model.Validate()
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
