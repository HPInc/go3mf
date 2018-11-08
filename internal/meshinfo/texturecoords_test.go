package meshinfo

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	gomock "github.com/golang/mock/gomock"
)

func TestTextureCoords_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		t    *TextureCoords
	}{
		{"base", new(TextureCoords)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.TextureID = 2
			tt.t.Coords[0] = mgl32.Vec2{1.0, 2.0}
			tt.t.Coords[1] = mgl32.Vec2{5.0, 3.0}
			tt.t.Coords[2] = mgl32.Vec2{6.0, 4.0}
			tt.t.Invalidate()
			want := new(TextureCoords)
			if !reflect.DeepEqual(tt.t, want) {
				t.Errorf("TextureCoords.Invalidate() = %v, want %v", tt.t, want)
			}
		})
	}
}

func TestTextureCoords_Copy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockData := NewMockFaceData(mockCtrl)
	test := &TextureCoords{2, [3]mgl32.Vec2{{1.0, 2.0}, {5.0, 3.0}, {6.0, 4.0}}}
	type args struct {
		from interface{}
	}
	tests := []struct {
		name string
		t    *TextureCoords
		args args
		want *TextureCoords
	}{
		{"nil", new(TextureCoords), args{nil}, new(TextureCoords)},
		{"othertype", new(TextureCoords), args{mockData}, new(TextureCoords)},
		{"copied", new(TextureCoords), args{test}, test},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Copy(tt.args.from)
		})
		if !reflect.DeepEqual(tt.t, tt.want) {
			t.Errorf("TextureCoords.Copy() = %v, want %v", tt.t, tt.want)
		}
	}
}

func TestTextureCoords_HasData(t *testing.T) {
	tests := []struct {
		name string
		t    *TextureCoords
		want bool
	}{
		{"nodata", new(TextureCoords), false},
		{"data", &TextureCoords{2, [3]mgl32.Vec2{{1.0, 2.0}, {5.0, 3.0}, {6.0, 4.0}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.HasData(); got != tt.want {
				t.Errorf("TextureCoords.HasData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureCoords_Permute(t *testing.T) {
	test := &TextureCoords{2, [3]mgl32.Vec2{{1.0, 2.0}, {5.0, 3.0}, {6.0, 4.0}}}
	type args struct {
		index1 uint32
		index2 uint32
		index3 uint32
	}
	tests := []struct {
		name string
		t    *TextureCoords
		args args
		want *TextureCoords
	}{
		{"big1", test, args{3, 1, 0}, test},
		{"big2", test, args{2, 3, 0}, test},
		{"big3", test, args{2, 1, 3}, test},
		{"success1", test, args{2, 1, 0}, &TextureCoords{2, [3]mgl32.Vec2{{6.0, 4.0}, {5.0, 3.0}, {1.0, 2.0}}}},
		{"success2", test, args{1, 2, 0}, &TextureCoords{2, [3]mgl32.Vec2{{5.0, 3.0}, {6.0, 4.0}, {1.0, 2.0}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Permute(tt.args.index1, tt.args.index2, tt.args.index3)
		})
	}
}

func TestTextureCoords_Merge(t *testing.T) {
	type args struct {
		other interface{}
	}
	tests := []struct {
		name string
		t    *TextureCoords
		args args
	}{
		{"notimplemented", new(TextureCoords), args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Merge(tt.args.other)
		})
	}
}
