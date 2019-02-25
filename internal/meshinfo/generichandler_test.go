package meshinfo

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func Test_newgenericHandler(t *testing.T) {
	tests := []struct {
		name string
		want *genericHandler
	}{
		{"new", &genericHandler{
			internalIDCounter: 1,
			lookup:            map[dataType]Handleable{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newgenericHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newgenericHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_addInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockHandleable := NewMockHandleable(mockCtrl)
	h := newgenericHandler()
	herr := newgenericHandler()
	herr.internalIDCounter = maxInternalID
	type args struct {
		info *MockHandleable
	}
	tests := []struct {
		name               string
		h                  *genericHandler
		args               args
		wantPanic          bool
		expectedInternalID uint64
	}{
		{"1", h, args{mockHandleable}, false, 1},
		{"2", h, args{mockHandleable}, false, 2},
		{"3", h, args{mockHandleable}, false, 3},
		{"max", herr, args{mockHandleable}, true, maxInternalID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.wantPanic && r == nil {
					t.Error("genericHandler.addInformation() want panic")
				}
			}()
			tt.args.info.EXPECT().InfoType().Return(nodeColorType)
			tt.args.info.EXPECT().setInternalID(tt.expectedInternalID)
			tt.h.addInformation(tt.args.info)
		})
	}
}

func sameTypeSlice(x, y []dataType) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[dataType]int, len(x))
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

func TestHandler_infoTypes(t *testing.T) {
	h := newgenericHandler()
	h.lookup[nodeColorType] = nil
	h.lookup[textureCoordsType] = nil
	h.lookup[baseMaterialType] = nil
	tests := []struct {
		name string
		h    *genericHandler
		want []dataType
	}{
		{"types", h, []dataType{nodeColorType, textureCoordsType, baseMaterialType}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.infoTypes(); !sameTypeSlice(got, tt.want) {
				t.Errorf("genericHandler.infoTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_AddFace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	h := newgenericHandler()
	handleable := NewMockHandleable(mockCtrl)
	h.lookup[nodeColorType] = handleable
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name string
		h    *genericHandler
		args args
		data *MockFaceData
	}{
		{"success", h, args{3}, NewMockFaceData(mockCtrl)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleable.EXPECT().AddFaceData(tt.args.newFaceCount).Return(tt.data)
			tt.data.EXPECT().Invalidate().Return()
			tt.h.AddFace(tt.args.newFaceCount)
		})
	}
}

func TestHandler_informationByType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	h := newgenericHandler()
	handleable1 := NewMockHandleable(mockCtrl)
	handleable2 := NewMockHandleable(mockCtrl)
	h.lookup[nodeColorType] = handleable1
	h.lookup[baseMaterialType] = handleable2
	type args struct {
		infoType dataType
	}
	tests := []struct {
		name  string
		h     *genericHandler
		args  args
		want  Handleable
		want1 bool
	}{
		{"valid1", h, args{nodeColorType}, handleable1, true},
		{"valid1", h, args{baseMaterialType}, handleable2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.h.informationByType(tt.args.infoType)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genericHandler.informationByType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("genericHandler.informationByType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHandler_InformationCount(t *testing.T) {
	h := newgenericHandler()
	h.lookup[nodeColorType] = nil
	h.lookup[textureCoordsType] = nil
	tests := []struct {
		name string
		h    *genericHandler
		want uint32
	}{
		{"empty", new(genericHandler), 0},
		{"withdata", h, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.InformationCount(); got != tt.want {
				t.Errorf("genericHandler.InformationCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_AddInfoFrom(t *testing.T) {
	types := []dataType{nodeColorType, textureCoordsType, baseMaterialType}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	otherHandleable := NewMockHandleable(mockCtrl)
	ownHandleable := NewMockHandleable(mockCtrl)
	h := newgenericHandler()
	h.lookup[textureCoordsType] = ownHandleable
	h.lookup[baseMaterialType] = ownHandleable
	type args struct {
		otherHandler     *MockTypedInformer
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		h    *genericHandler
		args args
	}{
		{"added", h, args{NewMockTypedInformer(mockCtrl), 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.otherHandler.EXPECT().infoTypes().Return(types)
			tt.args.otherHandler.EXPECT().informationByType(gomock.Any()).Return(otherHandleable, true).MaxTimes(3)
			otherHandleable.EXPECT().clone(tt.args.currentFaceCount).Return(ownHandleable)
			ownHandleable.EXPECT().InfoType().Return(nodeColorType)
			ownHandleable.EXPECT().setInternalID(tt.h.internalIDCounter)
			tt.h.AddInfoFrom(tt.args.otherHandler, tt.args.currentFaceCount)
		})
	}
}

func TestHandler_CopyFaceInfosFrom(t *testing.T) {
	types := []dataType{nodeColorType, textureCoordsType, baseMaterialType}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	otherHandleable := NewMockHandleable(mockCtrl)
	ownHandleable := NewMockHandleable(mockCtrl)
	h := newgenericHandler()
	h.lookup[textureCoordsType] = ownHandleable
	h.lookup[baseMaterialType] = ownHandleable
	type args struct {
		faceIndex      uint32
		otherHandler   *MockTypedInformer
		otherFaceIndex uint32
	}
	tests := []struct {
		name string
		h    *genericHandler
		args args
	}{
		{"base", h, args{2, NewMockTypedInformer(mockCtrl), 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.otherHandler.EXPECT().infoTypes().Return(types)
			tt.args.otherHandler.EXPECT().informationByType(gomock.Any()).Return(otherHandleable, true).MaxTimes(3)
			ownHandleable.EXPECT().copyFaceInfosFrom(tt.args.faceIndex, ownHandleable, tt.args.otherFaceIndex).MaxTimes(2)
			tt.h.CopyFaceInfosFrom(tt.args.faceIndex, tt.args.otherHandler, tt.args.otherFaceIndex)
		})
	}
}

func TestHandler_ResetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	handleable := NewMockHandleable(mockCtrl)
	h := newgenericHandler()
	h.lookup[textureCoordsType] = handleable
	h.lookup[baseMaterialType] = handleable
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		h    *genericHandler
		args args
	}{
		{"base", h, args{2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleable.EXPECT().resetFaceInformation(tt.args.faceIndex).MaxTimes(2)
			tt.h.ResetFaceInformation(tt.args.faceIndex)
		})
	}
}

func TestHandler_removeInformation(t *testing.T) {
	h := newgenericHandler()
	h.lookup[textureCoordsType] = nil
	h.lookup[baseMaterialType] = nil
	type args struct {
		infoType dataType
	}
	tests := []struct {
		name string
		h    *genericHandler
		args args
		want int
	}{
		{"other", h, args{nodeColorType}, 2},
		{"1", h, args{baseMaterialType}, 1},
		{"0", h, args{textureCoordsType}, 0},
		{"empty", h, args{textureCoordsType}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.removeInformation(tt.args.infoType)
			if got := len(tt.h.lookup); got != tt.want {
				t.Errorf("genericHandler.removeInformation() want = %v, got %v", tt.want, got)
			}
		})
	}
}

func TestHandler_PermuteNodeInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	handleable := NewMockHandleable(mockCtrl)
	h := newgenericHandler()
	h.lookup[textureCoordsType] = handleable
	h.lookup[baseMaterialType] = handleable
	type args struct {
		faceIndex  uint32
		nodeIndex1 uint32
		nodeIndex2 uint32
		nodeIndex3 uint32
	}
	tests := []struct {
		name string
		h    *genericHandler
		args args
	}{
		{"base", h, args{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleable.EXPECT().permuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3).MaxTimes(2)
			tt.h.PermuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
		})
	}
}

func Test_genericHandler_RemoveAllInformations(t *testing.T) {
	h := newgenericHandler()
	h.lookup[textureCoordsType] = nil
	h.lookup[baseMaterialType] = nil
	tests := []struct {
		name string
		h    *genericHandler
	}{
		{"empty", newgenericHandler()},
		{"notempty", h},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.RemoveAllInformations()
			if got := len(tt.h.lookup); got != 0 {
				t.Errorf("genericHandler.removeInformation() want = %v, got %v", 0, got)
			}
		})
	}
}
