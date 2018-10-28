package meshinfo

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_newbaseMeshInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	type args struct {
		container Container
	}
	tests := []struct {
		name string
		args args
		want *baseMeshInfo
	}{
		{"new", args{mockContainer}, &baseMeshInfo{mockContainer, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newbaseMeshInfo(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newbaseMeshInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseMeshInfo_resetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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
		{"error", newbaseMeshInfo(mockContainer), args{2}, true},
		{"success", newbaseMeshInfo(mockContainer), args{4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInvalidator := NewMockInvalidator(mockCtrl)
			var (
				err   error
				times int
			)
			if tt.wantErr {
				err = errors.New("")
			} else {
				times = 1
			}

			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(mockInvalidator, err)
			mockInvalidator.EXPECT().Invalidate().Times(times)
			tt.b.resetFaceInformation(tt.args.faceIndex)
		})
	}
}

func Test_baseMeshInfo_Clear(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	tests := []struct {
		name    string
		b       *baseMeshInfo
		faceNum uint32
	}{
		{"empty", newbaseMeshInfo(mockContainer), 0},
		{"one", newbaseMeshInfo(mockContainer), 1},
		{"two", newbaseMeshInfo(mockContainer), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInvalidator := NewMockInvalidator(mockCtrl)
			mockContainer.EXPECT().GetCurrentFaceCount().Return(tt.faceNum)
			mockContainer.EXPECT().GetFaceData(gomock.Any()).Return(mockInvalidator, nil).Times(int(tt.faceNum))
			mockInvalidator.EXPECT().Invalidate().Times(int(tt.faceNum))
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
		{"zero", newbaseMeshInfo(nil), args{0}},
		{"one", newbaseMeshInfo(nil), args{1}},
		{"two", newbaseMeshInfo(nil), args{3}},
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
		{"new", newbaseMeshInfo(nil), 0},
		{"one", &baseMeshInfo{nil, 1}, 1},
		{"two", &baseMeshInfo{nil, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getInternalID(); got != tt.want {
				t.Errorf("baseMeshInfo.getInternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}
