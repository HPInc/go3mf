package materials

import (
	"encoding/xml"
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
		{"empty", args{new(go3mf.Model)}, nil},
		{"noNamespace", args{&go3mf.Model{Resources: go3mf.Resources{Assets: []go3mf.Asset{
			&ColorGroupResource{ID: 1},
		}}}}, nil},
		{"child", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Childs: map[string]*go3mf.ChildModel{
			"/other.model": &go3mf.ChildModel{Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&ColorGroupResource{ID: 1},
			}}},
			"/that.model": &go3mf.ChildModel{Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&MultiPropertiesResource{ID: 2},
			}}},
		}}}, []error{
			&specerr.AssetError{Path: "/other.model", Index: 0, Name: "ColorGroupResource", Err: specerr.ErrEmptyResourceProps},
			&specerr.AssetError{Path: "/that.model", Index: 0, Name: "MultiPropertiesResource", Err: &specerr.MissingFieldError{Name: attrPIDs}},
			&specerr.AssetError{Path: "/that.model", Index: 0, Name: "MultiPropertiesResource", Err: specerr.ErrMultiBlend},
			&specerr.AssetError{Path: "/that.model", Index: 0, Name: "MultiPropertiesResource", Err: specerr.ErrEmptyResourceProps},
		}},
		{"multi", args{&go3mf.Model{
			Namespaces: []xml.Name{{Space: ExtensionName}},
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&MultiPropertiesResource{ID: 4},
				&MultiPropertiesResource{ID: 5, Multis: []Multi{{PIndex: []uint32{}}}, PIDs: []uint32{4, 100}},
				&go3mf.BaseMaterialsResource{ID: 1, Materials: []go3mf.BaseMaterial{
					{Name: "a", Color: color.RGBA{R: 1}},
					{Name: "b", Color: color.RGBA{G: 1}},
				}},
				&ColorGroupResource{ID: 6, Colors: []color.RGBA{{R: 1}, {R: 2, G: 3, B: 4, A: 5}}},
				&CompositeMaterialsResource{ID: 3, MaterialID: 1, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&MultiPropertiesResource{ID: 2, Multis: []Multi{{PIndex: []uint32{1, 0}}}, PIDs: []uint32{1, 6}},
				&MultiPropertiesResource{ID: 7, Multis: []Multi{{PIndex: []uint32{1, 3}}}, PIDs: []uint32{1, 6}},
				&MultiPropertiesResource{ID: 8, Multis: []Multi{{PIndex: []uint32{}}}, PIDs: []uint32{6, 1, 6}},
				&MultiPropertiesResource{ID: 9, Multis: []Multi{{PIndex: []uint32{}}}, PIDs: []uint32{1, 3}},
			}},
		}}, []error{
			&specerr.AssetError{Path: rootPath, Index: 0, Name: "MultiPropertiesResource", Err: &specerr.MissingFieldError{Name: attrPIDs}},
			&specerr.AssetError{Path: rootPath, Index: 0, Name: "MultiPropertiesResource", Err: specerr.ErrMultiBlend},
			&specerr.AssetError{Path: rootPath, Index: 0, Name: "MultiPropertiesResource", Err: specerr.ErrEmptyResourceProps},
			&specerr.AssetError{Path: rootPath, Index: 1, Name: "MultiPropertiesResource", Err: specerr.ErrMultiRefMulti},
			&specerr.AssetError{Path: rootPath, Index: 1, Name: "MultiPropertiesResource", Err: specerr.ErrMissingResource},
			&specerr.AssetError{Path: rootPath, Index: 6, Name: "MultiPropertiesResource", Err: &specerr.ResourcePropertyError{
				Index: 0,
				Err:   specerr.ErrIndexOutOfBounds,
			}},
			&specerr.AssetError{Path: rootPath, Index: 7, Name: "MultiPropertiesResource", Err: specerr.ErrMaterialMulti},
			&specerr.AssetError{Path: rootPath, Index: 7, Name: "MultiPropertiesResource", Err: specerr.ErrMultiColors},
			&specerr.AssetError{Path: rootPath, Index: 8, Name: "MultiPropertiesResource", Err: specerr.ErrMaterialMulti},
		}},
		{"missingTextPart", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}},
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&Texture2DResource{ID: 1},
				&Texture2DResource{ID: 2, ContentType: TextureTypePNG, Path: "/a.png"},
			}}},
		}, []error{
			&specerr.AssetError{Path: rootPath, Index: 0, Name: "Texture2DResource", Err: &specerr.MissingFieldError{Name: attrPath}},
			&specerr.AssetError{Path: rootPath, Index: 0, Name: "Texture2DResource", Err: &specerr.MissingFieldError{Name: attrContentType}},
			&specerr.AssetError{Path: rootPath, Index: 1, Name: "Texture2DResource", Err: ErrMissingTexturePart},
		}},
		{"textureGroup", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}},
			Attachments: []go3mf.Attachment{{Path: "/a.png"}},
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&Texture2DResource{ID: 1, ContentType: TextureTypePNG, Path: "/a.png"},
				&Texture2DGroupResource{ID: 2},
				&Texture2DGroupResource{ID: 3, TextureID: 1, Coords: []TextureCoord{{}}},
				&Texture2DGroupResource{ID: 4, TextureID: 2, Coords: []TextureCoord{{}}},
				&Texture2DGroupResource{ID: 5, TextureID: 100, Coords: []TextureCoord{{}}},
			}}},
		}, []error{
			&specerr.AssetError{Path: rootPath, Index: 1, Name: "Texture2DGroupResource", Err: &specerr.MissingFieldError{Name: attrTexID}},
			&specerr.AssetError{Path: rootPath, Index: 1, Name: "Texture2DGroupResource", Err: specerr.ErrEmptyResourceProps},
			&specerr.AssetError{Path: rootPath, Index: 3, Name: "Texture2DGroupResource", Err: ErrTextureReference},
			&specerr.AssetError{Path: rootPath, Index: 4, Name: "Texture2DGroupResource", Err: ErrTextureReference},
		}},
		{"colorGroup", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}},
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&ColorGroupResource{ID: 1},
				&ColorGroupResource{ID: 2, Colors: []color.RGBA{{R: 1}, {R: 2, G: 3, B: 4, A: 5}}},
				&ColorGroupResource{ID: 3, Colors: []color.RGBA{{R: 1}, {}}},
			}}}}, []error{
			&specerr.AssetError{Path: rootPath, Index: 0, Name: "ColorGroupResource", Err: specerr.ErrEmptyResourceProps},
			&specerr.AssetError{Path: rootPath, Index: 2, Name: "ColorGroupResource", Err: &specerr.ResourcePropertyError{
				Index: 1,
				Err:   &specerr.MissingFieldError{Name: attrColor},
			}},
		}},
		{"composite", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}},
			Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&go3mf.BaseMaterialsResource{ID: 1, Materials: []go3mf.BaseMaterial{
					{Name: "a", Color: color.RGBA{R: 1}},
					{Name: "b", Color: color.RGBA{G: 1}},
				}},
				&CompositeMaterialsResource{ID: 2},
				&CompositeMaterialsResource{ID: 3, MaterialID: 1, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&CompositeMaterialsResource{ID: 4, MaterialID: 1, Indices: []uint32{100, 100}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&CompositeMaterialsResource{ID: 5, MaterialID: 2, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
				&CompositeMaterialsResource{ID: 6, MaterialID: 100, Indices: []uint32{0, 1}, Composites: []Composite{{Values: []float32{1, 2}}}},
			}}}}, []error{
			&specerr.AssetError{Path: rootPath, Index: 1, Name: "CompositeMaterialsResource", Err: &specerr.MissingFieldError{Name: attrMatID}},
			&specerr.AssetError{Path: rootPath, Index: 1, Name: "CompositeMaterialsResource", Err: &specerr.MissingFieldError{Name: attrMatIndices}},
			&specerr.AssetError{Path: rootPath, Index: 1, Name: "CompositeMaterialsResource", Err: specerr.ErrEmptyResourceProps},
			&specerr.AssetError{Path: rootPath, Index: 3, Name: "CompositeMaterialsResource", Err: specerr.ErrIndexOutOfBounds},
			&specerr.AssetError{Path: rootPath, Index: 4, Name: "CompositeMaterialsResource", Err: ErrCompositeBase},
			&specerr.AssetError{Path: rootPath, Index: 5, Name: "CompositeMaterialsResource", Err: specerr.ErrMissingResource},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validate(tt.args.model)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
