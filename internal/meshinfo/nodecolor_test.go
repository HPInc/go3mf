package meshinfo

import (
	"image/color"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNodeColor_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		n    *NodeColor
	}{
		{"base", &NodeColor{[3]color.RGBA{color.RGBA{}, color.RGBA{}, color.RGBA{}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Invalidate()
			want := new(NodeColor)
			if !reflect.DeepEqual(tt.n, want) {
				t.Errorf("NodeColor.Invalidate() = %v, want %v", tt.n, want)
			}
		})
	}
}

func TestNodeColor_Copy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockData := NewMockFaceData(mockCtrl)
	targetColor := [3]color.RGBA{color.RGBA{1, 2, 3, 4}, color.RGBA{3, 1, 2, 3}, color.RGBA{1, 2, 34, 3}}
	type args struct {
		from interface{}
	}
	tests := []struct {
		name string
		n    *NodeColor
		args args
		want *NodeColor
	}{
		{"nil", new(NodeColor), args{nil}, new(NodeColor)},
		{"othertype", new(NodeColor), args{mockData}, new(NodeColor)},
		{"copied", new(NodeColor), args{&NodeColor{targetColor}}, &NodeColor{targetColor}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Copy(tt.args.from)
			if !reflect.DeepEqual(tt.n, tt.want) {
				t.Errorf("NodeColor.Copy() = %v, want %v", tt.n, tt.want)
			}
		})
	}
}

func TestNodeColor_HasData(t *testing.T) {
	tests := []struct {
		name string
		n    *NodeColor
		want bool
	}{
		{"nodata", new(NodeColor), false},
		{"data1", &NodeColor{[3]color.RGBA{color.RGBA{1, 2, 3, 4}, color.RGBA{}, color.RGBA{}}}, true},
		{"data2", &NodeColor{[3]color.RGBA{color.RGBA{}, color.RGBA{1, 2, 3, 4}, color.RGBA{}}}, true},
		{"data3", &NodeColor{[3]color.RGBA{color.RGBA{}, color.RGBA{}, color.RGBA{1, 2, 3, 4}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.HasData(); got != tt.want {
				t.Errorf("NodeColor.HasData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func nr(c uint8) color.RGBA {
	return color.RGBA{c, c, c, c}
}

func TestNodeColor_Permute(t *testing.T) {
	type args struct {
		index1 uint32
		index2 uint32
		index3 uint32
	}
	tests := []struct {
		name string
		n    *NodeColor
		args args
		want *NodeColor
	}{
		{"big1", &NodeColor{[3]color.RGBA{nr(1), nr(2), nr(3)}}, args{3, 1, 0}, &NodeColor{[3]color.RGBA{nr(1), nr(2), nr(3)}}},
		{"big2", &NodeColor{[3]color.RGBA{nr(1), nr(2), nr(3)}}, args{2, 3, 0}, &NodeColor{[3]color.RGBA{nr(1), nr(2), nr(3)}}},
		{"big3", &NodeColor{[3]color.RGBA{nr(1), nr(2), nr(3)}}, args{2, 1, 3}, &NodeColor{[3]color.RGBA{nr(1), nr(2), nr(3)}}},
		{"success1", &NodeColor{[3]color.RGBA{nr(1), nr(2), nr(3)}}, args{2, 1, 0}, &NodeColor{[3]color.RGBA{nr(3), nr(2), nr(1)}}},
		{"success2", &NodeColor{[3]color.RGBA{nr(1), nr(2), nr(3)}}, args{1, 2, 0}, &NodeColor{[3]color.RGBA{nr(2), nr(3), nr(1)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Permute(tt.args.index1, tt.args.index2, tt.args.index3)
			if !reflect.DeepEqual(tt.n, tt.want) {
				t.Errorf("NodeColor.Permute() = %v, want %v", tt.n, tt.want)
			}
		})
	}
}

func TestNodeColor_Merge(t *testing.T) {
	type args struct {
		other interface{}
	}
	tests := []struct {
		name string
		n    *NodeColor
		args args
	}{
		{"notimplemented", new(NodeColor), args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Merge(tt.args.other)
		})
	}
}
