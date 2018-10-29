package meshinfo

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNewBaseMaterial(t *testing.T) {
	type args struct {
		groupID uint32
		index   uint32
	}
	tests := []struct {
		name string
		args args
		want *BaseMaterial
	}{
		{"base", args{1, 2}, &BaseMaterial{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBaseMaterial(tt.args.groupID, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseMaterial() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockData := NewMockFaceData(mockCtrl)
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
		{"othertype", new(BaseMaterial), args{mockData}, new(BaseMaterial)},
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
