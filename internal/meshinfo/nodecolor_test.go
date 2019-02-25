package meshinfo

import (
	"image/color"
	"reflect"
	"testing"
)

func TestNodeColor_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		n    *NodeColor
	}{
		{"base", &NodeColor{[3]color.RGBA{{}, {}, {}}}},
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
	targetColor := [3]color.RGBA{{1, 2, 3, 4}, {3, 1, 2, 3}, {1, 2, 34, 3}}
	type args struct {
		from FaceData
	}
	tests := []struct {
		name string
		n    *NodeColor
		args args
		want *NodeColor
	}{
		{"nil", new(NodeColor), args{nil}, new(NodeColor)},
		{"othertype", new(NodeColor), args{new(BaseMaterial)}, new(NodeColor)},
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
		{"data1", &NodeColor{[3]color.RGBA{{1, 2, 3, 4}, {}, {}}}, true},
		{"data2", &NodeColor{[3]color.RGBA{{}, {1, 2, 3, 4}, {}}}, true},
		{"data3", &NodeColor{[3]color.RGBA{{}, {}, {1, 2, 3, 4}}}, true},
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
		other FaceData
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

func Test_newnodeColorContainer(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		args args
		want *nodeColorContainer
	}{
		{"zero", args{0}, &nodeColorContainer{0, make([]*NodeColor, 0)}},
		{"one", args{1}, &nodeColorContainer{1, []*NodeColor{new(NodeColor)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newnodeColorContainer(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newnodeColorContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeColorContainer_clone(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		m    *nodeColorContainer
		args args
		want Container
	}{
		{"zero", &nodeColorContainer{0, make([]*NodeColor, 0)}, args{0}, &nodeColorContainer{0, make([]*NodeColor, 0)}},
		{"one", &nodeColorContainer{1, []*NodeColor{new(NodeColor)}}, args{1}, &nodeColorContainer{1, []*NodeColor{new(NodeColor)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeColorContainer.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeColorContainer_InfoType(t *testing.T) {
	tests := []struct {
		name string
		m    *nodeColorContainer
		want DataType
	}{
		{"base", new(nodeColorContainer), NodeColorType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeColorContainer.InfoType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeColorContainer_AddFaceData(t *testing.T) {
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name      string
		m         *nodeColorContainer
		args      args
		want      FaceData
		wantPanic bool
	}{
		{"invalid face number", new(nodeColorContainer), args{2}, new(NodeColor), true},
		{"valid face number", new(nodeColorContainer), args{1}, new(NodeColor), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.wantPanic && r == nil {
					t.Error("memoryContainer.AddFaceData() want panic")
				}
			}()
			if got := tt.m.AddFaceData(tt.args.newFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeColorContainer.AddFaceData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeColorContainer_FaceData(t *testing.T) {
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		m    *nodeColorContainer
		args args
		want FaceData
	}{
		{"valid index", newnodeColorContainer(1), args{0}, new(NodeColor)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.FaceData(tt.args.faceIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeColorContainer.FaceData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeColorContainer_FaceCount(t *testing.T) {
	tests := []struct {
		name string
		m    *nodeColorContainer
		want uint32
	}{
		{"empty", new(nodeColorContainer), 0},
		{"1", newnodeColorContainer(1), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.FaceCount(); got != tt.want {
				t.Errorf("nodeColorContainer.FaceCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeColorContainer_Clear(t *testing.T) {
	tests := []struct {
		name string
		m    *nodeColorContainer
	}{
		{"base", new(nodeColorContainer)},
		{"1", newnodeColorContainer(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Clear()
			if got := tt.m.FaceCount(); got != 0 {
				t.Errorf("nodeColorContainer.Clear() = %v, want 0", got)
			}
		})
	}
}
