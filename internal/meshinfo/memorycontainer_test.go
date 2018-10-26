package meshinfo

import (
	"reflect"
	"testing"
)

type fakeFaceData struct {
	a int
}

func TestNewMemoryContainer(t *testing.T) {
	type args struct {
		currentFaceCount uint32
		elemType         reflect.Type
	}
	tests := []struct {
		name string
		args args
	}{
		{"zero", args{0, reflect.TypeOf(fakeFaceData{})}},
		{"one", args{1, reflect.TypeOf(fakeFaceData{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newMemoryContainer(tt.args.currentFaceCount, tt.args.elemType)
			if got.GetCurrentFaceCount() != tt.args.currentFaceCount || got.elemType != tt.args.elemType {
				t.Error("newMemoryContainer() created an invalid container")
			}
		})
	}
}

func Test_memoryContainer_Clone(t *testing.T) {
	p := newMemoryContainer(0, reflect.TypeOf(fakeFaceData{}))
	p.AddFaceData(1)
	tests := []struct {
		name string
		m    *memoryContainer
		want *memoryContainer
	}{
		{"empty", p, newMemoryContainer(1, reflect.TypeOf(fakeFaceData{}))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Clone(); got.GetCurrentFaceCount() != tt.want.faceCount {
				t.Errorf("memoryContainer.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryContainer_AddFaceData(t *testing.T) {
	m := newMemoryContainer(0, reflect.TypeOf(fakeFaceData{}))
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name    string
		m       *memoryContainer
		args    args
		wantVal FaceData
		wantErr bool
	}{
		{"invalid element type", &memoryContainer{nil, 0, reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(fakeFaceData{})), 0, 0)}, args{0}, nil, true},
		{"invalid face number", m, args{0}, nil, true},
		{"valid face number", m, args{2}, fakeFaceData{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.m.AddFaceData(tt.args.newFaceCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryContainer.AddFaceData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reflect.TypeOf(gotVal) == reflect.TypeOf(&tt.wantVal) {
				t.Errorf("memoryContainer.AddFaceData() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestMemoryContainer_GetFaceData(t *testing.T) {
	m := newMemoryContainer(0, reflect.TypeOf(fakeFaceData{}))
	initial, _ := m.AddFaceData(1)
	type args struct {
		index uint32
	}
	tests := []struct {
		name    string
		m       *memoryContainer
		args    args
		wantVal FaceData
		wantErr bool
	}{
		{"invalid index", m, args{1}, nil, true},
		{"valid index", m, args{0}, fakeFaceData{4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.m.GetFaceData(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryContainer.GetFaceData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got := gotVal.(*fakeFaceData)
				if !(got != initial) {
					t.Errorf("memoryContainer.GetFaceData() = %v, want %v", got, initial)
				}
				got.a = tt.wantVal.(fakeFaceData).a
				newVal, _ := tt.m.GetFaceData(tt.args.index)
				newGot := newVal.(*fakeFaceData)
				if !reflect.DeepEqual(*newGot, tt.wantVal) {
					t.Errorf("memoryContainer.GetFaceData() = %v, want %v", newGot, tt.wantVal)
				}
			}
		})
	}
}

func TestMemoryContainer_GetCurrentFaceCount(t *testing.T) {
	m := newMemoryContainer(0, reflect.TypeOf(fakeFaceData{}))
	mempty := newMemoryContainer(0, reflect.TypeOf(fakeFaceData{}))
	m.AddFaceData(1)
	tests := []struct {
		name string
		m    *memoryContainer
		want uint32
	}{
		{"empty", mempty, 0},
		{"one", m, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetCurrentFaceCount(); got != tt.want {
				t.Errorf("memoryContainer.GetCurrentFaceCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryContainer_Clear(t *testing.T) {
	m := newMemoryContainer(0, reflect.TypeOf(fakeFaceData{}))
	m.AddFaceData(1)
	tests := []struct {
		name string
		m    *memoryContainer
	}{
		{"base", m},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Clear()
			if got := tt.m.GetCurrentFaceCount(); got != 0 {
				t.Errorf("memoryContainer.Clear() = %v, want %v", got, 0)
			}
		})
	}
}
