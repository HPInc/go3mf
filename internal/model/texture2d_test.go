package model

import (
	"reflect"
	"testing"
)

func newTexture2D() *Texture2DResource {
	r, _ := NewTexture2DResource(0, new(Model))
	return r
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
		{"equal", &Texture2DResource{Path: "/a.png", ContentType: PNGTexture}, args{&Texture2DResource{Path: "/a.png", ContentType: PNGTexture}}},
		{"diff", &Texture2DResource{Path: "/b.png", ContentType: PNGTexture}, args{&Texture2DResource{Path: "/a.png", ContentType: JPEGTexture}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Copy(tt.args.other)
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
			Resource:    Resource{Model: model, ResourceID: &ResourceID{"", 0, 1}},
			ContentType: PNGTexture,
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
