package meshinfo

import (
	"reflect"
	"testing"
)

func TestBaseMaterial_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseMaterial
	}{
		{"base", &BaseMaterial{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Invalidate()
			want := new(BaseMaterial)
			if !reflect.DeepEqual(tt.b, want) {
				t.Errorf("BaseMaterial.Invalidate() = %v, want %v", tt.b, want)
			}
		})
	}
}

func TestBaseMaterial_Copy(t *testing.T) {
	type args struct {
		from interface{}
	}
	tests := []struct {
		name string
		b    *BaseMaterial
		args args
		want *BaseMaterial
	}{
		{"nil", new(BaseMaterial), args{nil}, new(BaseMaterial)},
		{"othertype", new(BaseMaterial), args{new(NodeColor)}, new(BaseMaterial)},
		{"copied", new(BaseMaterial), args{&BaseMaterial{2, 3}}, &BaseMaterial{2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Copy(tt.args.from)
			if !reflect.DeepEqual(tt.b, tt.want) {
				t.Errorf("BaseMaterial.Copy() = %v, want %v", tt.b, tt.want)
			}
		})
	}
}

func TestBaseMaterial_HasData(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseMaterial
		want bool
	}{
		{"nodata", new(BaseMaterial), false},
		{"data", &BaseMaterial{2, 3}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.HasData(); got != tt.want {
				t.Errorf("BaseMaterial.HasData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterial_Permute(t *testing.T) {
	type args struct {
		index1 uint32
		index2 uint32
		index3 uint32
	}
	tests := []struct {
		name string
		b    *BaseMaterial
		args args
	}{
		{"notimplemented", new(BaseMaterial), args{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Permute(tt.args.index1, tt.args.index2, tt.args.index3)
		})
	}
}

func TestBaseMaterial_Merge(t *testing.T) {
	type args struct {
		other interface{}
	}
	tests := []struct {
		name string
		b    *BaseMaterial
		args args
	}{
		{"notimplemented", new(BaseMaterial), args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Merge(tt.args.other)
		})
	}
}

func Test_newbaseMaterialContainer(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		args args
		want *baseMaterialContainer
	}{
		{"zero", args{0}, &baseMaterialContainer{0, make([]*BaseMaterial, 0)}},
		{"one", args{1}, &baseMaterialContainer{1, []*BaseMaterial{new(BaseMaterial)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newbaseMaterialContainer(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newbaseMaterialContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMaterialContainer_clone(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		m    *baseMaterialContainer
		args args
		want Container
	}{
		{"zero", &baseMaterialContainer{0, make([]*BaseMaterial, 0)}, args{0}, &baseMaterialContainer{0, make([]*BaseMaterial, 0)}},
		{"one", &baseMaterialContainer{1, []*BaseMaterial{new(BaseMaterial)}}, args{1}, &baseMaterialContainer{1, []*BaseMaterial{new(BaseMaterial)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseMaterialContainer.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMaterialContainer_InfoType(t *testing.T) {
	tests := []struct {
		name string
		m    *baseMaterialContainer
		want dataType
	}{
		{"base", new(baseMaterialContainer), baseMaterialType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseMaterialContainer.InfoType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMaterialContainer_AddFaceData(t *testing.T) {
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name      string
		m         *baseMaterialContainer
		args      args
		want      FaceData
		wantPanic bool
	}{
		{"invalid face number", new(baseMaterialContainer), args{2}, new(BaseMaterial), true},
		{"valid face number", new(baseMaterialContainer), args{1}, new(BaseMaterial), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.wantPanic && r == nil {
					t.Error("memoryContainer.AddFaceData() want panic")
				}
			}()
			if got := tt.m.AddFaceData(tt.args.newFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseMaterialContainer.AddFaceData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMaterialContainer_FaceData(t *testing.T) {
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		m    *baseMaterialContainer
		args args
		want FaceData
	}{
		{"valid index", newbaseMaterialContainer(1), args{0}, new(BaseMaterial)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.FaceData(tt.args.faceIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseMaterialContainer.FaceData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMaterialContainer_FaceCount(t *testing.T) {
	tests := []struct {
		name string
		m    *baseMaterialContainer
		want uint32
	}{
		{"empty", new(baseMaterialContainer), 0},
		{"1", newbaseMaterialContainer(1), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.FaceCount(); got != tt.want {
				t.Errorf("baseMaterialContainer.FaceCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMaterialContainer_Clear(t *testing.T) {
	tests := []struct {
		name string
		m    *baseMaterialContainer
	}{
		{"base", new(baseMaterialContainer)},
		{"1", newbaseMaterialContainer(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Clear()
		})
	}
}
