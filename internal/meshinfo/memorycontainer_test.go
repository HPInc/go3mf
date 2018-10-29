package meshinfo

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNewmemoryContainer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockFaceData(mockCtrl)
	type args struct {
		currentFaceCount uint32
		infoType         reflect.Type
	}
	tests := []struct {
		name string
		args args
	}{
		{"zero", args{0, reflect.TypeOf(*mockContainer)}},
		{"one", args{1, reflect.TypeOf(*mockContainer)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newmemoryContainer(tt.args.currentFaceCount, tt.args.infoType).(*memoryContainer)
			if got.GetCurrentFaceCount() != tt.args.currentFaceCount || got.infoType != tt.args.infoType {
				t.Error("newmemoryContainer() created an invalid container")
			}
		})
	}
}

func Test_memoryContainer_Clone(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockFaceData(mockCtrl)
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		m    *memoryContainer
		args args
		want *memoryContainer
	}{
		{"empty", newmemoryContainer(0, reflect.TypeOf(*mockContainer)).(*memoryContainer), args{2}, newmemoryContainer(1, reflect.TypeOf(*mockContainer)).(*memoryContainer)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Clone(tt.args.currentFaceCount); got.GetCurrentFaceCount() != tt.args.currentFaceCount {
				t.Errorf("memoryContainer.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryContainer_AddFaceData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFaceData := NewMockFaceData(mockCtrl)
	m := newmemoryContainer(0, reflect.TypeOf(*mockFaceData)).(*memoryContainer)
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
		{"invalid data type", &memoryContainer{nil, 0, reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(*mockFaceData)), 0, 0)}, args{0}, nil, true},
		{"invalid face number", m, args{0}, nil, true},
		{"valid face number", m, args{2}, mockFaceData, false},
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

func Test_memoryContainer_GetFaceData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFaceData := NewMockFaceData(mockCtrl)
	m := newmemoryContainer(0, reflect.TypeOf(*mockFaceData)).(*memoryContainer)
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
		{"valid index", m, args{0}, mockFaceData, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.m.GetFaceData(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryContainer.GetFaceData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got := gotVal
				if !(got != initial) {
					t.Errorf("memoryContainer.GetFaceData() = %v, want %v", got, initial)
				}
			}
		})
	}
}

func Test_memoryContainer_GetCurrentFaceCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFaceData := NewMockFaceData(mockCtrl)
	m := newmemoryContainer(0, reflect.TypeOf(*mockFaceData)).(*memoryContainer)
	mempty := newmemoryContainer(0, reflect.TypeOf(*mockFaceData)).(*memoryContainer)
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

func Test_memoryContainer_Clear(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFaceData := NewMockFaceData(mockCtrl)
	m := newmemoryContainer(0, reflect.TypeOf(*mockFaceData)).(*memoryContainer)
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

func Test_memoryContainer_InfoType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFaceData := NewMockFaceData(mockCtrl)
	tests := []struct {
		name string
		m    *memoryContainer
		want reflect.Type
	}{
		{"base", newmemoryContainer(0, reflect.TypeOf(*mockFaceData)).(*memoryContainer), reflect.TypeOf(*mockFaceData)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryContainer.InfoType() = %v, want %v", got, tt.want)
			}
		})
	}
}
