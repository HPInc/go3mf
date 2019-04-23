package go3mf

import "testing"

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

func TestTexture2DResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		t     *Texture2DResource
		want  string
		want1 uint32
	}{
		{"base", &Texture2DResource{ID: 1, ModelPath: "3d/3dmodel.model"}, "3d/3dmodel.model", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.t.Identify()
			if got != tt.want {
				t.Errorf("Texture2DResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Texture2DResource.Identify() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTextureCoord_U(t *testing.T) {
	tests := []struct {
		name string
		t    TextureCoord
		want float32
	}{
		{"base", TextureCoord{1, 2}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.U(); got != tt.want {
				t.Errorf("TextureCoord.U() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureCoord_V(t *testing.T) {
	tests := []struct {
		name string
		t    TextureCoord
		want float32
	}{
		{"base", TextureCoord{1, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.V(); got != tt.want {
				t.Errorf("TextureCoord.V() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTexture2DGroupResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		t     *Texture2DGroupResource
		want  string
		want1 uint32
	}{
		{"base", &Texture2DGroupResource{ID: 1, ModelPath: "3d/3dmodel"}, "3d/3dmodel", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.t.Identify()
			if got != tt.want {
				t.Errorf("Texture2DGroupResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Texture2DGroupResource.Identify() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestColorGroupResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		c     *ColorGroupResource
		want  string
		want1 uint32
	}{
		{"base", &ColorGroupResource{ID: 1, ModelPath: "3d/3dmodel"}, "3d/3dmodel", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Identify()
			if got != tt.want {
				t.Errorf("ColorGroupResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ColorGroupResource.Identify() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCompositeMaterialsResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		c     *CompositeMaterialsResource
		want  string
		want1 uint32
	}{
		{"base", &CompositeMaterialsResource{ID: 1, ModelPath: "3d/3dmodel.model"}, "3d/3dmodel.model", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Identify()
			if got != tt.want {
				t.Errorf("CompositeMaterialsResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CompositeMaterialsResource.Identify() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
