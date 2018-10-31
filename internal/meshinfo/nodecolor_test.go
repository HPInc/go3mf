package meshinfo

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNodeColor_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		n    *NodeColor
	}{
		{"base", &NodeColor{[3]Color{1, 2, 3}}},
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
		{"copied", new(NodeColor), args{&NodeColor{[3]Color{1, 2, 3}}}, &NodeColor{[3]Color{1, 2, 3}}},
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
		{"data1", &NodeColor{[3]Color{1, 0, 0}}, true},
		{"data2", &NodeColor{[3]Color{0, 1, 0}}, true},
		{"data3", &NodeColor{[3]Color{0, 0, 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.HasData(); got != tt.want {
				t.Errorf("NodeColor.HasData() = %v, want %v", got, tt.want)
			}
		})
	}
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
		{"big1", &NodeColor{[3]Color{1, 2, 3}}, args{3, 1, 0}, &NodeColor{[3]Color{1, 2, 3}}},
		{"big2", &NodeColor{[3]Color{1, 2, 3}}, args{2, 3, 0}, &NodeColor{[3]Color{1, 2, 3}}},
		{"big3", &NodeColor{[3]Color{1, 2, 3}}, args{2, 1, 3}, &NodeColor{[3]Color{1, 2, 3}}},
		{"success1", &NodeColor{[3]Color{1, 2, 3}}, args{2, 1, 0}, &NodeColor{[3]Color{3, 2, 1}}},
		{"success2", &NodeColor{[3]Color{1, 2, 3}}, args{1, 2, 0}, &NodeColor{[3]Color{2, 3, 1}}},
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
