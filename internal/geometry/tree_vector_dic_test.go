package geometry

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewTreeVectorDic(t *testing.T) {
	tests := []struct {
		name string
		want *TreeVectorDic
	}{
		{"new", &TreeVectorDic{0.001, map[Vec3I]uint32{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTreeVectorDic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTreeVectorDic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTreeVectorDicWithUnits(t *testing.T) {
	type args struct {
		units float32
	}
	tests := []struct {
		name    string
		args    args
		want    *TreeVectorDic
		wantErr bool
	}{
		{"zero", args{0.0}, &TreeVectorDic{0.0, map[Vec3I]uint32{}}, true},
		{"one", args{1.0}, &TreeVectorDic{1.0, map[Vec3I]uint32{}}, false},
		{"big", args{1001.0}, &TreeVectorDic{0.0, map[Vec3I]uint32{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTreeVectorDicWithUnits(tt.args.units)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTreeVectorDicWithUnits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTreeVectorDicWithUnits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTreeVectorDic_Units(t *testing.T) {
	p := NewTreeVectorDic()
	if got := p.Units(); got != VectorDefaultUnits {
		t.Errorf("TreeVectorDic.Units() = %v, want %v", got, VectorDefaultUnits)
	}
	p.SetUnits(1.0)
	if got := p.Units(); got != 1.0 {
		t.Errorf("TreeVectorDic.Units() = %v, want %v", got, VectorDefaultUnits)
	}
}

func TestTreeVectorDic_SetUnits(t *testing.T) {
	p := NewTreeVectorDic()
	p.AddVector(mgl32.Vec3{}, 1)
	type args struct {
		units float32
	}
	tests := []struct {
		name    string
		t       *TreeVectorDic
		args    args
		wantErr bool
	}{
		{"zero", NewTreeVectorDic(), args{0.0}, true},
		{"minunitsfail", NewTreeVectorDic(), args{0.000009}, true},
		{"minunits", NewTreeVectorDic(), args{VectorMinUnits}, false},
		{"one", NewTreeVectorDic(), args{1.0}, false},
		{"maxunits", NewTreeVectorDic(), args{VectorMaxUnits}, false},
		{"maxunitsfail", NewTreeVectorDic(), args{1001.0}, true},
		{"notempty", p, args{1.0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.SetUnits(tt.args.units); (err != nil) != tt.wantErr {
				t.Errorf("TreeVectorDic.SetUnits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want := tt.args.units
			if tt.wantErr {
				want = VectorDefaultUnits
			}
			if got := tt.t.Units(); got != want {
				t.Errorf("TreeVectorDic.Units() = %v, want %v", got, want)
			}
		})
	}
}

func TestTreeVectorDic_AddFindVector(t *testing.T) {
	p := NewTreeVectorDic()
	type args struct {
		vec   mgl32.Vec3
		value uint32
	}
	tests := []struct {
		name string
		t    *TreeVectorDic
		args args
	}{
		{"new", p, args{mgl32.Vec3{10000.3, 20000.2, 1}, 2}},
		{"old", p, args{mgl32.Vec3{10000.3, 20000.2, 1}, 4}},
		{"new2", p, args{mgl32.Vec3{2, 1, 3.4}, 5}},
		{"old2", p, args{mgl32.Vec3{2, 1, 3.4}, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.AddVector(tt.args.vec, tt.args.value)
		})
		got, ok := p.FindVector(tt.args.vec)
		if !ok {
			t.Error("TreeVectorDic.AddMatch() haven't added the match")
			return
		}
		if got != tt.args.value {
			t.Errorf("TreeVectorDic.FindVector() = %v, want %v", got, tt.args.value)
		}
	}
}

func TestTreeVectorDic_RemoveVector(t *testing.T) {
	p := NewTreeVectorDic()
	p.AddVector(mgl32.Vec3{1, 2, 5.3}, 1)
	type args struct {
		vec mgl32.Vec3
	}
	tests := []struct {
		name string
		t    *TreeVectorDic
		args args
	}{
		{"nil", p, args{mgl32.Vec3{2, 3, 4}}},
		{"old", p, args{mgl32.Vec3{1, 2, 5.3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.RemoveVector(tt.args.vec)
		})
	}
}

func Test_newVec3IFromVec3(t *testing.T) {
	type args struct {
		vec   mgl32.Vec3
		units float32
	}
	tests := []struct {
		name string
		args args
		want Vec3I
	}{
		{"1", args{mgl32.Vec3{1.2, 2.3, 3.4}, 1.0}, Vec3I{1, 2, 3}},
		{"2", args{mgl32.Vec3{1.2, 2.3, 3.4}, 0.001}, Vec3I{1200, 2299, 3400}},
		{"3", args{mgl32.Vec3{1.2, 2.3, 3.4}, 1000.0}, Vec3I{0, 0, 0}},
		{"4", args{mgl32.Vec3{1000.2, 2000.3, 3000.4}, 1000.0}, Vec3I{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newVec3IFromVec3(tt.args.vec, tt.args.units); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newVec3IFromVec3() = %v, want %v", got, tt.want)
			}
		})
	}
}
