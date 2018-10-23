package meshinfo

import (
	"reflect"
	"testing"
)

type fakeFaceData struct {
	a int
}

func TestNewInMemoryMeshInformationContainer(t *testing.T) {
	type args struct {
		currentFaceCount uint32
		elemExample      FaceData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil", args{0, nil}, true},
		{"zero", args{0, fakeFaceData{}}, false},
		{"one", args{1, fakeFaceData{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newInMemoryMeshInformationContainer(tt.args.currentFaceCount, tt.args.elemExample)
			if (err != nil) != tt.wantErr {
				t.Errorf("newInMemoryMeshInformationContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && (got.getCurrentFaceCount() != tt.args.currentFaceCount || got.elemType != reflect.TypeOf(tt.args.elemExample)) {
				t.Error("newInMemoryMeshInformationContainer() created an invalid container")
			}
		})
	}
}

func TestInMemoryMeshInformationContainer_AddFaceData(t *testing.T) {
	m, _ := newInMemoryMeshInformationContainer(0, fakeFaceData{})
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name    string
		m       *inMemoryMeshInformationContainer
		args    args
		wantVal FaceData
		wantErr bool
	}{
		{"invalid element type", &inMemoryMeshInformationContainer{nil, 0, reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(fakeFaceData{})), 0, 0)}, args{0}, nil, true},
		{"invalid face number", m, args{0}, nil, true},
		{"valid face number", m, args{2}, fakeFaceData{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.m.addFaceData(tt.args.newFaceCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemoryMeshInformationContainer.addFaceData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reflect.TypeOf(gotVal) == reflect.TypeOf(&tt.wantVal) {
				t.Errorf("inMemoryMeshInformationContainer.addFaceData() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestInMemoryMeshInformationContainer_GetFaceData(t *testing.T) {
	m, _ := newInMemoryMeshInformationContainer(0, fakeFaceData{})
	initial, _ := m.addFaceData(1)
	type args struct {
		index uint32
	}
	tests := []struct {
		name    string
		m       *inMemoryMeshInformationContainer
		args    args
		wantVal FaceData
		wantErr bool
	}{
		{"invalid index", m, args{1}, nil, true},
		{"valid index", m, args{0}, fakeFaceData{4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.m.getFaceData(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemoryMeshInformationContainer.getFaceData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got := gotVal.(*fakeFaceData)
				if !(got != initial) {
					t.Errorf("inMemoryMeshInformationContainer.getFaceData() = %v, want %v", got, initial)
				}
				got.a = tt.wantVal.(fakeFaceData).a
				newVal, _ := tt.m.getFaceData(tt.args.index)
				newGot := newVal.(*fakeFaceData)
				if !reflect.DeepEqual(*newGot, tt.wantVal) {
					t.Errorf("inMemoryMeshInformationContainer.getFaceData() = %v, want %v", newGot, tt.wantVal)
				}
			}
		})
	}
}

func TestInMemoryMeshInformationContainer_GetCurrentFaceCount(t *testing.T) {
	m, _ := newInMemoryMeshInformationContainer(0, fakeFaceData{})
	mempty, _ := newInMemoryMeshInformationContainer(0, fakeFaceData{})
	m.addFaceData(1)
	tests := []struct {
		name string
		m    *inMemoryMeshInformationContainer
		want uint32
	}{
		{"empty", mempty, 0},
		{"one", m, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getCurrentFaceCount(); got != tt.want {
				t.Errorf("inMemoryMeshInformationContainer.getCurrentFaceCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryMeshInformationContainer_Clear(t *testing.T) {
	m, _ := newInMemoryMeshInformationContainer(0, fakeFaceData{})
	m.addFaceData(1)
	tests := []struct {
		name string
		m    *inMemoryMeshInformationContainer
	}{
		{"base", m},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.clear()
			if got := tt.m.getCurrentFaceCount(); got != 0 {
				t.Errorf("inMemoryMeshInformationContainer.clear() = %v, want %v", got, 0)
			}
		})
	}
}
