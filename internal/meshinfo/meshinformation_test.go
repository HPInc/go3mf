package meshinfo

import (
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_newBaseMeshInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockInvalidator := NewMockInvalidator(mockCtrl)
	mockContainer := NewMockMeshInformationContainer(mockCtrl)
	type args struct {
		container   MeshInformationContainer
		invalidator Invalidator
	}
	tests := []struct {
		name string
		args args
		want *baseMeshInformation
	}{
		{"new", args{mockContainer, mockInvalidator}, &baseMeshInformation{mockContainer, mockInvalidator, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newBaseMeshInformation(tt.args.container, tt.args.invalidator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newBaseMeshInformation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMeshInformation_ResetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockInvalidator := NewMockInvalidator(mockCtrl)
	mockContainer := NewMockMeshInformationContainer(mockCtrl)

	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name     string
		b        *baseMeshInformation
		args     args
		wantData bool
		wantErr  bool
	}{
		{"nil data", newBaseMeshInformation(mockContainer, mockInvalidator), args{1}, false, false},
		{"error", newBaseMeshInformation(mockContainer, mockInvalidator), args{2}, true, true},
		{"nil data and error", newBaseMeshInformation(mockContainer, mockInvalidator), args{3}, false, true},
		{"success", newBaseMeshInformation(mockContainer, mockInvalidator), args{4}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				data  FaceData
				err   error
				times int
			)
			if tt.wantData {
				data = &fakeFaceData{}
			}
			if tt.wantErr {
				err = errors.New("")
			}
			if tt.wantData && !tt.wantErr {
				times = 1
			}

			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(data, err)
			mockInvalidator.EXPECT().Invalidate(data).Times(times)
			tt.b.ResetFaceInformation(tt.args.faceIndex)
		})
	}
}

func Test_baseMeshInformation_Clear(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockInvalidator := NewMockInvalidator(mockCtrl)
	mockContainer := NewMockMeshInformationContainer(mockCtrl)
	tests := []struct {
		name    string
		b       *baseMeshInformation
		faceNum uint32
	}{
		{"empty", newBaseMeshInformation(mockContainer, mockInvalidator), 0},
		{"one", newBaseMeshInformation(mockContainer, mockInvalidator), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().GetCurrentFaceCount().Return(tt.faceNum)
			mockContainer.EXPECT().GetFaceData(gomock.Any()).Return(nil, nil).Times(int(tt.faceNum))
			tt.b.Clear()
		})
	}
}

func Test_baseMeshInformation_setInternalID(t *testing.T) {
	type args struct {
		internalID uint64
	}
	tests := []struct {
		name string
		b    *baseMeshInformation
		args args
	}{
		{"zero", newBaseMeshInformation(nil, nil), args{0}},
		{"one", newBaseMeshInformation(nil, nil), args{1}},
		{"two", newBaseMeshInformation(nil, nil), args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.setInternalID(tt.args.internalID)
			if got := tt.b.internalID; got != tt.args.internalID {
				t.Errorf("baseMeshInformation.setInternalID() = %v, want %v", got, tt.args.internalID)
			}
		})
	}
}

func Test_baseMeshInformation_getInternalID(t *testing.T) {
	tests := []struct {
		name string
		b    *baseMeshInformation
		want uint64
	}{
		{"new", newBaseMeshInformation(nil, nil), 0},
		{"one", &baseMeshInformation{nil, nil, 1}, 1},
		{"two", &baseMeshInformation{nil, nil, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getInternalID(); got != tt.want {
				t.Errorf("baseMeshInformation.getInternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}
