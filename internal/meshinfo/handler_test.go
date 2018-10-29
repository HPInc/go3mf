package meshinfo

import (
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name string
		want *Handler
	}{
		{"new", &Handler{
			internalIDCounter: 1,
			lookup:            map[reflect.Type]MeshInfo{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_AddInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMesh := NewMockMeshInfo(mockCtrl)
	h := NewHandler()
	herr := NewHandler()
	herr.internalIDCounter = maxInternalID
	type args struct {
		info MeshInfo
	}
	tests := []struct {
		name               string
		h                  *Handler
		args               args
		wantErr            bool
		expectedInternalID uint64
	}{
		{"1", h, args{mockMesh}, false, 1},
		{"2", h, args{mockMesh}, false, 2},
		{"3", h, args{mockMesh}, false, 3},
		{"max", herr, args{mockMesh}, true, maxInternalID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.info.(*MockMeshInfo).EXPECT().InfoType().Return(reflect.TypeOf(""))
			tt.args.info.(*MockMeshInfo).EXPECT().setInternalID(tt.expectedInternalID)
			if err := tt.h.AddInformation(tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("Handler.AddInformation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_AddFace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	h := NewHandler()
	meshInfo := NewMockMeshInfo(mockCtrl)
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = meshInfo
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		data    *MockFaceData
		err     error
		wantErr bool
	}{
		{"err1", h, args{3}, NewMockFaceData(mockCtrl), errors.New(""), true},
		{"success", h, args{3}, NewMockFaceData(mockCtrl), nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meshInfo.EXPECT().AddFaceData(tt.args.newFaceCount).Return(tt.data, tt.err)
			if tt.err == nil {
				tt.data.EXPECT().Invalidate().Return()
			}
			if err := tt.h.AddFace(tt.args.newFaceCount); (err != nil) != tt.wantErr {
				t.Errorf("Handler.AddFace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_GetInformationByType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	h := NewHandler()
	meshInfo1 := NewMockMeshInfo(mockCtrl)
	meshInfo2 := NewMockMeshInfo(mockCtrl)
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = meshInfo1
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = meshInfo2
	type args struct {
		infoType reflect.Type
	}
	tests := []struct {
		name  string
		h     *Handler
		args  args
		want  MeshInfo
		want1 bool
	}{
		{"nil", h, args{nil}, nil, false},
		{"valid1", h, args{reflect.TypeOf((*string)(nil)).Elem()}, meshInfo1, true},
		{"valid1", h, args{reflect.TypeOf((*float32)(nil)).Elem()}, meshInfo2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.h.GetInformationByType(tt.args.infoType)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.GetInformationByType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Handler.GetInformationByType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHandler_GetInformationCount(t *testing.T) {
	h := NewHandler()
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = nil
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = nil
	tests := []struct {
		name string
		h    *Handler
		want uint32
	}{
		{"empty", new(Handler), 0},
		{"withdata", h, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.GetInformationCount(); got != tt.want {
				t.Errorf("Handler.GetInformationCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_AddInfoFromTable(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	otherMeshInfo := NewMockMeshInfo(mockCtrl)
	ownMeshInfo := NewMockMeshInfo(mockCtrl)
	herr := NewHandler()
	herr.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = ownMeshInfo
	herr.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = ownMeshInfo
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = ownMeshInfo
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = ownMeshInfo
	type args struct {
		otherHandler     *Handler
		currentFaceCount uint32
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
		{"error", herr, args{NewHandler(), 3}, true},
		{"added", h, args{NewHandler(), 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.otherHandler.lookup[reflect.TypeOf((*string)(nil)).Elem()] = otherMeshInfo
			tt.args.otherHandler.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = otherMeshInfo
			tt.args.otherHandler.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = otherMeshInfo
			if tt.wantErr {
				tt.h.internalIDCounter = maxInternalID
			} else {
				ownMeshInfo.EXPECT().mergeInformationFrom(otherMeshInfo).MaxTimes(3)
			}
			otherMeshInfo.EXPECT().Clone(tt.args.currentFaceCount).Return(ownMeshInfo)
			ownMeshInfo.EXPECT().InfoType().Return(reflect.TypeOf((*string)(nil)).Elem())
			ownMeshInfo.EXPECT().setInternalID(tt.h.internalIDCounter)
			if err := tt.h.AddInfoFromTable(tt.args.otherHandler, tt.args.currentFaceCount); (err != nil) != tt.wantErr {
				t.Errorf("Handler.AddInfoFromTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_CloneFaceInfosFrom(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	otherMeshInfo := NewMockMeshInfo(mockCtrl)
	ownMeshInfo := NewMockMeshInfo(mockCtrl)
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = ownMeshInfo
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = ownMeshInfo
	type args struct {
		faceIndex      uint32
		otherHandler   *Handler
		otherFaceIndex uint32
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{"base", h, args{2, NewHandler(), 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.otherHandler.lookup[reflect.TypeOf((*string)(nil)).Elem()] = otherMeshInfo
			tt.args.otherHandler.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = otherMeshInfo
			tt.args.otherHandler.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = otherMeshInfo
			ownMeshInfo.EXPECT().cloneFaceInfosFrom(tt.args.faceIndex, ownMeshInfo, tt.args.otherFaceIndex).MaxTimes(2)
			tt.h.CloneFaceInfosFrom(tt.args.faceIndex, tt.args.otherHandler, tt.args.otherFaceIndex)
		})
	}
}

func TestHandler_ResetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	meshInfo := NewMockMeshInfo(mockCtrl)
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = meshInfo
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = meshInfo
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{"base", h, args{2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meshInfo.EXPECT().resetFaceInformation(tt.args.faceIndex).MaxTimes(2)
			tt.h.ResetFaceInformation(tt.args.faceIndex)
		})
	}
}

func TestHandler_RemoveInformation(t *testing.T) {
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = nil
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = nil
	type args struct {
		infoType reflect.Type
	}
	tests := []struct {
		name string
		h    *Handler
		args args
		want int
	}{
		{"other", h, args{reflect.TypeOf((*string)(nil)).Elem()}, 2},
		{"1", h, args{reflect.TypeOf((*float64)(nil)).Elem()}, 1},
		{"0", h, args{reflect.TypeOf((*float32)(nil)).Elem()}, 0},
		{"empty", h, args{reflect.TypeOf((*float32)(nil)).Elem()}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.RemoveInformation(tt.args.infoType)
			if got := len(tt.h.lookup); got != tt.want {
				t.Errorf("Handler.RemoveInformation() want = %v, got %v", tt.want, got)
			}
		})
	}
}

func TestHandler_PermuteNodeInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	meshInfo := NewMockMeshInfo(mockCtrl)
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = meshInfo
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = meshInfo
	type args struct {
		faceIndex  uint32
		nodeIndex1 uint32
		nodeIndex2 uint32
		nodeIndex3 uint32
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{"base", h, args{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meshInfo.EXPECT().permuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3).MaxTimes(2)
			tt.h.PermuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
		})
	}
}
