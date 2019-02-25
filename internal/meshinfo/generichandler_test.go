package meshinfo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
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
	h := newgenericHandler()
	herr := newgenericHandler()
	herr.internalIDCounter = maxInternalID
	tests := []struct {
		name               string
		h                  *genericHandler
		wantPanic          bool
		expectedInternalID uint64
	}{
		{"1", h, false, 1},
		{"2", h, false, 2},
		{"3", h, false, 3},
		{"max", herr, true, maxInternalID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.wantPanic && r == nil {
					t.Error("genericHandler.addInformation() want panic")
				}
			}()
			mockHandleable := new(MockHandleable)
			mockHandleable.On("InfoType").Return(nodeColorType)
			mockHandleable.On("setInternalID", tt.expectedInternalID)
			tt.h.addInformation(mockHandleable)
			mockHandleable.AssertExpectations(t)
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
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name string
		h    *genericHandler
		args args
		data *MockFaceData
	}{
		{"success", newgenericHandler(), args{3}, new(MockFaceData)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleable := new(MockHandleable)
			tt.h.lookup[nodeColorType] = handleable
			handleable.On("AddFaceData", tt.args.newFaceCount).Return(tt.data)
			tt.data.On("Invalidate")
			tt.h.AddFace(tt.args.newFaceCount)
			tt.data.AssertExpectations(t)
			handleable.AssertExpectations(t)
		})
	}
}

func TestHandler_informationByType(t *testing.T) {
	h := newgenericHandler()
	handleable1 := new(MockHandleable)
	handleable2 := new(MockHandleable)
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
	handleable1.AssertExpectations(t)
	handleable2.AssertExpectations(t)
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
	type args struct {
		otherHandler     *MockTypedInformer
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		h    *genericHandler
		args args
	}{
		{"added", newgenericHandler(), args{new(MockTypedInformer), 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otherHandleable := new(MockHandleable)
			ownHandleable := new(MockHandleable)
			tt.h.lookup[textureCoordsType] = ownHandleable
			tt.h.lookup[baseMaterialType] = ownHandleable
			tt.args.otherHandler.On("infoTypes").Return(types)
			tt.args.otherHandler.On("informationByType", mock.Anything).Return(otherHandleable, true).Times(3)
			otherHandleable.On("clone", tt.args.currentFaceCount).Return(ownHandleable)
			ownHandleable.On("InfoType").Return(nodeColorType)
			ownHandleable.On("setInternalID", tt.h.internalIDCounter)
			tt.h.AddInfoFrom(tt.args.otherHandler, tt.args.currentFaceCount)
			tt.args.otherHandler.AssertExpectations(t)
			otherHandleable.AssertExpectations(t)
			ownHandleable.AssertExpectations(t)
		})
	}
}

func TestHandler_CopyFaceInfosFrom(t *testing.T) {
	types := []dataType{nodeColorType, textureCoordsType, baseMaterialType}
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
		{"base", newgenericHandler(), args{2, new(MockTypedInformer), 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otherHandleable := new(MockHandleable)
			ownHandleable := new(MockHandleable)
			tt.h.lookup[textureCoordsType] = ownHandleable
			tt.h.lookup[baseMaterialType] = ownHandleable
			tt.args.otherHandler.On("infoTypes").Return(types)
			tt.args.otherHandler.On("informationByType", mock.Anything).Return(otherHandleable, true).Times(3)
			ownHandleable.On("copyFaceInfosFrom", tt.args.faceIndex, otherHandleable, tt.args.otherFaceIndex).Times(2)
			tt.h.CopyFaceInfosFrom(tt.args.faceIndex, tt.args.otherHandler, tt.args.otherFaceIndex)
			tt.args.otherHandler.AssertExpectations(t)
			otherHandleable.AssertExpectations(t)
			ownHandleable.AssertExpectations(t)
		})
	}
}

func TestHandler_ResetFaceInformation(t *testing.T) {
	handleable := new(MockHandleable)
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
			handleable.On("resetFaceInformation", tt.args.faceIndex).Times(2)
			tt.h.ResetFaceInformation(tt.args.faceIndex)
		})
	}
	handleable.AssertExpectations(t)
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
	handleable := new(MockHandleable)
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
			handleable.On("permuteNodeInformation", tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3).Times(2)
			tt.h.PermuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
		})
	}
	handleable.AssertExpectations(t)
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
