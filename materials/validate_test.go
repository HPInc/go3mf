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
	rootPath := go3mf.DefaultModelPath
	type args struct {
		model *go3mf.Model
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{"empty", args{new(go3mf.Model)}, []error{}},
		{"child", args{&go3mf.Model{Childs: map[string]*go3mf.ChildModel{
			"/other.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&ColorGroup{ID: 1},
			}}},
			"/that.model": {Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&MultiProperties{ID: 2},
			}}},
		}}}, []error{
			fmt.Errorf("/other.model@Resources@ColorGroup#0: %v", specerr.ErrEmptyResourceProps),
			fmt.Errorf("/that.model@Resources@MultiProperties#0: %v", &specerr.MissingFieldError{Name: attrPIDs}),
			fmt.Errorf("/that.model@Resources@MultiProperties#0: %v", specerr.ErrMultiBlend),
			fmt.Errorf("/that.model@Resources@MultiProperties#0: %v", specerr.ErrEmptyResourceProps),
		}},
		{"multi", args{&go3mf.Model{
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&MultiProperties{ID: 4},
				&MultiProperties{ID: 5, Multis: []Multi{{PIndex: []uint32{}}}, PIDs: []uint32{4, 100}},
				&go3mf.BaseMaterials{ID: 1, Materials: []go3mf.Base{
					{Name: "a", Color: color.RGBA{R: 1}},
					{Name: "b", Color: color.RGBA{G: 1}},
				}},
				&ColorGroup{ID: 6, Colors: []color.RGBA{{R: 1}, {R: 2, G: 3, B: 4, A: 5}}},
				&CompositeMaterials{ID: 3, MaterialID: 1, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&MultiProperties{ID: 2, Multis: []Multi{{PIndex: []uint32{1, 0}}}, PIDs: []uint32{1, 6}},
				&MultiProperties{ID: 7, Multis: []Multi{{PIndex: []uint32{1, 3}}}, PIDs: []uint32{1, 6}},
				&MultiProperties{ID: 8, Multis: []Multi{{PIndex: []uint32{}}}, PIDs: []uint32{6, 1, 6}},
				&MultiProperties{ID: 9, Multis: []Multi{{PIndex: []uint32{}}}, PIDs: []uint32{1, 3}},
			}},
		}}, []error{
			fmt.Errorf("%s@Resources@MultiProperties#0: %v", rootPath, &specerr.MissingFieldError{Name: attrPIDs}),
			fmt.Errorf("%s@Resources@MultiProperties#0: %v", rootPath, specerr.ErrMultiBlend),
			fmt.Errorf("%s@Resources@MultiProperties#0: %v", rootPath, specerr.ErrEmptyResourceProps),
			fmt.Errorf("%s@Resources@MultiProperties#1: %v", rootPath, specerr.ErrMultiRefMulti),
			fmt.Errorf("%s@Resources@MultiProperties#1: %v", rootPath, specerr.ErrMissingResource),
			fmt.Errorf("%s@Resources@MultiProperties#6@Multi#0: %v", rootPath, specerr.ErrIndexOutOfBounds),
			fmt.Errorf("%s@Resources@MultiProperties#7: %v", rootPath, specerr.ErrMaterialMulti),
			fmt.Errorf("%s@Resources@MultiProperties#7: %v", rootPath, specerr.ErrMultiColors),
			fmt.Errorf("%s@Resources@MultiProperties#8: %v", rootPath, specerr.ErrMaterialMulti),
		}},
		{"missingTextPart", args{&go3mf.Model{
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&Texture2D{ID: 1},
				&Texture2D{ID: 2, ContentType: TextureTypePNG, Path: "/a.png"},
			}}},
		}, []error{
			fmt.Errorf("%s@Resources@Texture2D#0: %v", rootPath, &specerr.MissingFieldError{Name: attrPath}),
			fmt.Errorf("%s@Resources@Texture2D#0: %v", rootPath, &specerr.MissingFieldError{Name: attrContentType}),
			fmt.Errorf("%s@Resources@Texture2D#1: %v", rootPath, specerr.ErrMissingTexturePart),
		}},
		{"textureGroup", args{&go3mf.Model{
			Attachments: []go3mf.Attachment{{Path: "/a.png"}},
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&Texture2D{ID: 1, ContentType: TextureTypePNG, Path: "/a.png"},
				&Texture2DGroup{ID: 2},
				&Texture2DGroup{ID: 3, TextureID: 1, Coords: []TextureCoord{{}}},
				&Texture2DGroup{ID: 4, TextureID: 2, Coords: []TextureCoord{{}}},
				&Texture2DGroup{ID: 5, TextureID: 100, Coords: []TextureCoord{{}}},
			}}},
		}, []error{
			fmt.Errorf("%s@Resources@Texture2DGroup#1: %v", rootPath, &specerr.MissingFieldError{Name: attrTexID}),
			fmt.Errorf("%s@Resources@Texture2DGroup#1: %v", rootPath, specerr.ErrEmptyResourceProps),
			fmt.Errorf("%s@Resources@Texture2DGroup#3: %v", rootPath, specerr.ErrTextureReference),
			fmt.Errorf("%s@Resources@Texture2DGroup#4: %v", rootPath, specerr.ErrTextureReference),
		}},
		{"colorGroup", args{&go3mf.Model{
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&ColorGroup{ID: 1},
				&ColorGroup{ID: 2, Colors: []color.RGBA{{R: 1}, {R: 2, G: 3, B: 4, A: 5}}},
				&ColorGroup{ID: 3, Colors: []color.RGBA{{R: 1}, {}}},
			}}}}, []error{
			fmt.Errorf("%s@Resources@ColorGroup#0: %v", rootPath, specerr.ErrEmptyResourceProps),
			fmt.Errorf("%s@Resources@ColorGroup#2@RGBA#1: %v", rootPath, &specerr.MissingFieldError{Name: attrColor}),
		}},
		{"composite", args{&go3mf.Model{
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
			}}}}, []error{
			fmt.Errorf("%s@Resources@CompositeMaterials#1: %v", rootPath, &specerr.MissingFieldError{Name: attrMatID}),
			fmt.Errorf("%s@Resources@CompositeMaterials#1: %v", rootPath, &specerr.MissingFieldError{Name: attrMatIndices}),
			fmt.Errorf("%s@Resources@CompositeMaterials#1: %v", rootPath, specerr.ErrEmptyResourceProps),
			fmt.Errorf("%s@Resources@CompositeMaterials#3: %v", rootPath, specerr.ErrIndexOutOfBounds),
			fmt.Errorf("%s@Resources@CompositeMaterials#4: %v", rootPath, specerr.ErrCompositeBase),
			fmt.Errorf("%s@Resources@CompositeMaterials#5: %v", rootPath, specerr.ErrMissingResource),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.model.Validate()
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
