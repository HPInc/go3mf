package meshinfo

import (
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_newbaseMeshInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockInvalidator := NewMockInvalidator(mockCtrl)
	mockContainer := NewMockContainer(mockCtrl)
	type args struct {
		container   Container
		invalidator Invalidator
	}
	tests := []struct {
		name string
		args args
		want *baseMeshInfo
	}{
		{"new", args{mockContainer, mockInvalidator}, &baseMeshInfo{mockContainer, mockInvalidator, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newbaseMeshInfo(tt.args.container, tt.args.invalidator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newbaseMeshInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMeshInfo_resetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockInvalidator := NewMockInvalidator(mockCtrl)
	mockContainer := NewMockContainer(mockCtrl)

	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name    string
		b       *baseMeshInfo
		args    args
		wantErr bool
	}{
		{"error", newbaseMeshInfo(mockContainer, mockInvalidator), args{2}, true},
		{"success", newbaseMeshInfo(mockContainer, mockInvalidator), args{4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &fakeFaceData{}
			var (
				err   error
				times int
			)
			if tt.wantErr {
				err = errors.New("")
			} else {
				times = 1
			}

			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(data, err)
			mockInvalidator.EXPECT().Invalidate(data).Times(times)
			tt.b.resetFaceInformation(tt.args.faceIndex)
		})
	}
}

func Test_baseMeshInfo_Clear(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockInvalidator := NewMockInvalidator(mockCtrl)
	mockContainer := NewMockContainer(mockCtrl)
	tests := []struct {
		name    string
		b       *baseMeshInfo
		faceNum uint32
	}{
		{"empty", newbaseMeshInfo(mockContainer, mockInvalidator), 0},
		{"one", newbaseMeshInfo(mockContainer, mockInvalidator), 1},
		{"two", newbaseMeshInfo(mockContainer, mockInvalidator), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &fakeFaceData{}
			mockContainer.EXPECT().GetCurrentFaceCount().Return(tt.faceNum)
			mockContainer.EXPECT().GetFaceData(gomock.Any()).Return(data, nil).Times(int(tt.faceNum))
			mockInvalidator.EXPECT().Invalidate(data).Times(int(tt.faceNum))
			tt.b.Clear()
		})
	}
}

func Test_baseMeshInfo_setInternalID(t *testing.T) {
	type args struct {
		internalID uint64
	}
	tests := []struct {
		name string
		b    *baseMeshInfo
		args args
	}{
		{"zero", newbaseMeshInfo(nil, nil), args{0}},
		{"one", newbaseMeshInfo(nil, nil), args{1}},
		{"two", newbaseMeshInfo(nil, nil), args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.setInternalID(tt.args.internalID)
			if got := tt.b.internalID; got != tt.args.internalID {
				t.Errorf("baseMeshInfo.setInternalID() = %v, want %v", got, tt.args.internalID)
			}
		})
	}
}

func Test_baseMeshInfo_getInternalID(t *testing.T) {
	tests := []struct {
		name string
		b    *baseMeshInfo
		want uint64
	}{
		{"new", newbaseMeshInfo(nil, nil), 0},
		{"one", &baseMeshInfo{nil, nil, 1}, 1},
		{"two", &baseMeshInfo{nil, nil, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getInternalID(); got != tt.want {
				t.Errorf("baseMeshInfo.getInternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}
