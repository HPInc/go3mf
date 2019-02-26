package model

import (
	"reflect"
	"testing"
)

func newTexture2D() *Texture2DResource {
	r, _ := NewTexture2DResource(0, new(Model))
	return r
}

func TestTexture2DResource_Box(t *testing.T) {
	tests := []struct {
		name       string
		t          *Texture2DResource
		wantU      float32
		wantV      float32
		wantWidth  float32
		wantHeight float32
		wantHasBox bool
	}{
		{"base", newTexture2D(), 0, 0, 1, 1, false},
		{"set", newTexture2D().SetBox(1, 2, 3, 4), 1, 2, 3, 4, true},
		{"clear", newTexture2D().SetBox(1, 2, 3, 4).ClearBox(), 0, 0, 1, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, gotV, gotWidth, gotHeight, gotHasBox := tt.t.Box()
			if gotU != tt.wantU {
				t.Errorf("Texture2DResource.Box() gotU = %v, want %v", gotU, tt.wantU)
			}
			if gotV != tt.wantV {
				t.Errorf("Texture2DResource.Box() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("Texture2DResource.Box() gotWidth = %v, want %v", gotWidth, tt.wantWidth)
			}
			if gotHeight != tt.wantHeight {
				t.Errorf("Texture2DResource.Box() gotHeight = %v, want %v", gotHeight, tt.wantHeight)
			}
			if gotHasBox != tt.wantHasBox {
				t.Errorf("Texture2DResource.Box() gotHasBox = %v, want %v", gotHasBox, tt.wantHasBox)
			}
		})
	}
}

func TestTexture2DResource_Copy(t *testing.T) {
	type args struct {
		other *Texture2DResource
	}
	tests := []struct {
		name string
		t    *Texture2DResource
		args args
	}{
		{"equal", newTexture2D().SetBox(1, 2, 3, 4), args{newTexture2D().SetBox(1, 2, 3, 4)}},
		{"diff", newTexture2D(), args{newTexture2D().SetBox(1, 2, 3, 4)}},
		{"noBox", newTexture2D().SetBox(1, 2, 3, 4), args{newTexture2D()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Copy(tt.args.other)
			gotU, gotV, gotWidth, gotHeight, gotHasBox := tt.t.Box()
			wantU, wantV, wantWidth, wantHeight, wantHasBox := tt.args.other.Box()
			if gotU != wantU {
				t.Errorf("Texture2DResource.Copy() gotU = %v, want %v", gotU, wantU)
			}
			if gotV != wantV {
				t.Errorf("Texture2DResource.Copy() gotV = %v, want %v", gotV, wantV)
			}
			if gotWidth != wantWidth {
				t.Errorf("Texture2DResource.Copy() gotWidth = %v, want %v", gotWidth, wantWidth)
			}
			if gotHeight != wantHeight {
				t.Errorf("Texture2DResource.Copy() gotHeight = %v, want %v", gotHeight, wantHeight)
			}
			if gotHasBox != wantHasBox {
				t.Errorf("Texture2DResource.Copy() gotHasBox = %v, want %v", gotHasBox, wantHasBox)
			}
			if tt.t.Path != tt.args.other.Path {
				t.Errorf("Texture2DResource.Copy() gotPath = %v, want %v", tt.t.Path, tt.args.other.Path)
			}
			if tt.t.ContentType != tt.args.other.ContentType {
				t.Errorf("Texture2DResource.Copy() gotContentType = %v, want %v", tt.t.ContentType, tt.args.other.ContentType)
			}
		})
	}
}

func TestNewTexture2DResource(t *testing.T) {
	model := new(Model)
	type args struct {
		id    uint64
		model *Model
	}
	tests := []struct {
		name    string
		args    args
		want    *Texture2DResource
		wantErr bool
	}{
		{"base", args{0, model}, &Texture2DResource{
			Resource:    Resource{Model: model, ResourceID: &PackageResourceID{"", 0, 1}},
			ContentType: UnknownTexture,
			boxWidth:    1,
			boxHeight:   1,
			TileStyleU:  WrapTile,
			TileStyleV:  WrapTile,
			Filter:      AutoFilter,
		}, false},
		{"dup", args{0, model}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTexture2DResource(tt.args.id, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTexture2DResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTexture2DResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
