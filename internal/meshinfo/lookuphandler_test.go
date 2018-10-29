package meshinfo

import (
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNewLookupHandler(t *testing.T) {
	tests := []struct {
		name string
		want *lookupHandler
	}{
		{"new", &lookupHandler{
			internalIDCounter: 1,
			lookup:            map[reflect.Type]MeshInfo{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLookupHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLookupHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lookupHandler_AddInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMesh := NewMockMeshInfo(mockCtrl)
	h := NewLookupHandler().(*lookupHandler)
	herr := NewLookupHandler().(*lookupHandler)
	herr.internalIDCounter = maxInternalID
	type args struct {
		info MeshInfo
	}
	tests := []struct {
		name               string
		h                  *lookupHandler
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
				t.Errorf("lookupHandler.AddInformation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lookupHandler_InfoTypes(t *testing.T) {
	h := NewLookupHandler().(*lookupHandler)
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = nil
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = nil
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = nil
	tests := []struct {
		name string
		h    *lookupHandler
		want []reflect.Type
	}{
		{"types", h, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*float32)(nil)).Elem(), reflect.TypeOf((*float64)(nil)).Elem()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.InfoTypes(); !sameTypeSlice(got, tt.want) {
				t.Errorf("lookupHandler.InfoTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lookupHandler_AddFace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	h := NewLookupHandler().(*lookupHandler)
	meshInfo := NewMockMeshInfo(mockCtrl)
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = meshInfo
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name    string
		h       *lookupHandler
		args    args
		data    *MockFaceData
		err     error
		wantErr bool
	}{
		{"err1", h, args{3}, NewMockFaceData(mockCtrl), errors.New(""), true},
		{"succsess", h, args{3}, NewMockFaceData(mockCtrl), nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meshInfo.EXPECT().AddFaceData(tt.args.newFaceCount).Return(tt.data, tt.err)
			if tt.err == nil {
				tt.data.EXPECT().Invalidate().Return()
			}
			if err := tt.h.AddFace(tt.args.newFaceCount); (err != nil) != tt.wantErr {
				t.Errorf("lookupHandler.AddFace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lookupHandler_GetInformationByType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	h := NewLookupHandler().(*lookupHandler)
	meshInfo1 := NewMockMeshInfo(mockCtrl)
	meshInfo2 := NewMockMeshInfo(mockCtrl)
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = meshInfo1
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = meshInfo2
	type args struct {
		infoType reflect.Type
	}
	tests := []struct {
		name  string
		h     *lookupHandler
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
				t.Errorf("lookupHandler.GetInformationByType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("lookupHandler.GetInformationByType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_lookupHandler_GetInformationCount(t *testing.T) {
	h := NewLookupHandler().(*lookupHandler)
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = nil
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = nil
	tests := []struct {
		name string
		h    *lookupHandler
		want uint32
	}{
		{"empty", new(lookupHandler), 0},
		{"withdata", h, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.GetInformationCount(); got != tt.want {
				t.Errorf("lookupHandler.GetInformationCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lookupHandler_AddInfoFromTable(t *testing.T) {
	types := []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*float32)(nil)).Elem(), reflect.TypeOf((*float64)(nil)).Elem()}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	otherMeshInfo := NewMockMeshInfo(mockCtrl)
	ownMeshInfo := NewMockMeshInfo(mockCtrl)
	herr := NewLookupHandler().(*lookupHandler)
	herr.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = ownMeshInfo
	herr.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = ownMeshInfo
	h := NewLookupHandler().(*lookupHandler)
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = ownMeshInfo
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = ownMeshInfo
	type args struct {
		otherHandler     *MockHandler
		currentFaceCount uint32
	}
	tests := []struct {
		name    string
		h       *lookupHandler
		args    args
		wantErr bool
	}{
		{"error", herr, args{NewMockHandler(mockCtrl), 3}, true},
		{"added", h, args{NewMockHandler(mockCtrl), 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.otherHandler.EXPECT().InfoTypes().Return(types)
			tt.args.otherHandler.EXPECT().GetInformationByType(gomock.Any()).Return(otherMeshInfo, true).MaxTimes(3)
			if tt.wantErr {
				tt.h.internalIDCounter = maxInternalID
			} else {
				ownMeshInfo.EXPECT().mergeInformationFrom(otherMeshInfo).MaxTimes(3)
			}
			otherMeshInfo.EXPECT().Clone(tt.args.currentFaceCount).Return(ownMeshInfo)
			ownMeshInfo.EXPECT().InfoType().Return(reflect.TypeOf((*string)(nil)).Elem())
			ownMeshInfo.EXPECT().setInternalID(tt.h.internalIDCounter)
			if err := tt.h.AddInfoFromTable(tt.args.otherHandler, tt.args.currentFaceCount); (err != nil) != tt.wantErr {
				t.Errorf("lookupHandler.AddInfoFromTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lookupHandler_CloneFaceInfosFrom(t *testing.T) {
	types := []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*float32)(nil)).Elem(), reflect.TypeOf((*float64)(nil)).Elem()}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	otherMeshInfo := NewMockMeshInfo(mockCtrl)
	ownMeshInfo := NewMockMeshInfo(mockCtrl)
	h := NewLookupHandler().(*lookupHandler)
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = ownMeshInfo
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = ownMeshInfo
	type args struct {
		faceIndex      uint32
		otherHandler   *MockHandler
		otherFaceIndex uint32
	}
	tests := []struct {
		name string
		h    *lookupHandler
		args args
	}{
		{"base", h, args{2, NewMockHandler(mockCtrl), 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.otherHandler.EXPECT().InfoTypes().Return(types)
			tt.args.otherHandler.EXPECT().GetInformationByType(gomock.Any()).Return(otherMeshInfo, true).MaxTimes(3)
			ownMeshInfo.EXPECT().cloneFaceInfosFrom(tt.args.faceIndex, ownMeshInfo, tt.args.otherFaceIndex).MaxTimes(2)
			tt.h.CloneFaceInfosFrom(tt.args.faceIndex, tt.args.otherHandler, tt.args.otherFaceIndex)
		})
	}
}

func Test_lookupHandler_ResetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	meshInfo := NewMockMeshInfo(mockCtrl)
	h := NewLookupHandler().(*lookupHandler)
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = meshInfo
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = meshInfo
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		h    *lookupHandler
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

func Test_lookupHandler_RemoveInformation(t *testing.T) {
	h := NewLookupHandler().(*lookupHandler)
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = nil
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = nil
	type args struct {
		infoType reflect.Type
	}
	tests := []struct {
		name string
		h    *lookupHandler
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
				t.Errorf("lookupHandler.RemoveInformation() want = %v, got %v", tt.want, got)
			}
		})
	}
}

func Test_lookupHandler_PermuteNodeInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	meshInfo := NewMockMeshInfo(mockCtrl)
	h := NewLookupHandler().(*lookupHandler)
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
		h    *lookupHandler
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

func sameTypeSlice(x, y []reflect.Type) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[reflect.Type]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}
	for _, _y := range y {
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	if len(diff) == 0 {
		return true
	}
	return false
}
