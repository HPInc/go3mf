package meshinfo

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func Test_NewHandler(t *testing.T) {
	tests := []struct {
		name string
		want *Handler
	}{
		{"new", &Handler{
			internalIDCounter: 1,
			lookup:            map[reflect.Type]Handleable{},
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

func TestHandler_addInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockHandleable := NewMockHandleable(mockCtrl)
	h := NewHandler()
	herr := NewHandler()
	herr.internalIDCounter = maxInternalID
	type args struct {
		info *MockHandleable
	}
	tests := []struct {
		name               string
		h                  *Handler
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
					t.Error("Handler.addInformation() want panic")
				}
			}()
			tt.args.info.EXPECT().InfoType().Return(reflect.TypeOf(""))
			tt.args.info.EXPECT().setInternalID(tt.expectedInternalID)
			tt.h.addInformation(tt.args.info)
		})
	}
}

func TestHandler_InfoTypes(t *testing.T) {
	h := NewHandler()
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = nil
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = nil
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = nil
	tests := []struct {
		name string
		h    *Handler
		want []reflect.Type
	}{
		{"types", h, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*float32)(nil)).Elem(), reflect.TypeOf((*float64)(nil)).Elem()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.InfoTypes(); !sameTypeSlice(got, tt.want) {
				t.Errorf("Handler.InfoTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_AddFace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	h := NewHandler()
	handleable := NewMockHandleable(mockCtrl)
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = handleable
	type args struct {
		newFaceCount uint32
	}
	tests := []struct {
		name string
		h    *Handler
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

func TestHandler_getInformationByType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	h := NewHandler()
	handleable1 := NewMockHandleable(mockCtrl)
	handleable2 := NewMockHandleable(mockCtrl)
	h.lookup[reflect.TypeOf((*string)(nil)).Elem()] = handleable1
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = handleable2
	type args struct {
		infoType reflect.Type
	}
	tests := []struct {
		name  string
		h     *Handler
		args  args
		want  Handleable
		want1 bool
	}{
		{"nil", h, args{nil}, nil, false},
		{"valid1", h, args{reflect.TypeOf((*string)(nil)).Elem()}, handleable1, true},
		{"valid1", h, args{reflect.TypeOf((*float32)(nil)).Elem()}, handleable2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.h.getInformationByType(tt.args.infoType)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.getInformationByType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Handler.getInformationByType() got1 = %v, want %v", got1, tt.want1)
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

func TestHandler_AddInfoFrom(t *testing.T) {
	types := []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*float32)(nil)).Elem(), reflect.TypeOf((*float64)(nil)).Elem()}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	otherHandleable := NewMockHandleable(mockCtrl)
	ownHandleable := NewMockHandleable(mockCtrl)
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = ownHandleable
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = ownHandleable
	type args struct {
		otherHandler     *MockTypedInformer
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{"added", h, args{NewMockTypedInformer(mockCtrl), 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.otherHandler.EXPECT().InfoTypes().Return(types)
			tt.args.otherHandler.EXPECT().getInformationByType(gomock.Any()).Return(otherHandleable, true).MaxTimes(3)
			otherHandleable.EXPECT().clone(tt.args.currentFaceCount).Return(ownHandleable)
			ownHandleable.EXPECT().InfoType().Return(reflect.TypeOf((*string)(nil)).Elem())
			ownHandleable.EXPECT().setInternalID(tt.h.internalIDCounter)
			tt.h.AddInfoFrom(tt.args.otherHandler, tt.args.currentFaceCount)
		})
	}
}

func TestHandler_CloneFaceInfosFrom(t *testing.T) {
	types := []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*float32)(nil)).Elem(), reflect.TypeOf((*float64)(nil)).Elem()}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	otherHandleable := NewMockHandleable(mockCtrl)
	ownHandleable := NewMockHandleable(mockCtrl)
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = ownHandleable
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = ownHandleable
	type args struct {
		faceIndex      uint32
		otherHandler   *MockTypedInformer
		otherFaceIndex uint32
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{"base", h, args{2, NewMockTypedInformer(mockCtrl), 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.otherHandler.EXPECT().InfoTypes().Return(types)
			tt.args.otherHandler.EXPECT().getInformationByType(gomock.Any()).Return(otherHandleable, true).MaxTimes(3)
			ownHandleable.EXPECT().cloneFaceInfosFrom(tt.args.faceIndex, ownHandleable, tt.args.otherFaceIndex).MaxTimes(2)
			tt.h.CloneFaceInfosFrom(tt.args.faceIndex, tt.args.otherHandler, tt.args.otherFaceIndex)
		})
	}
}

func TestHandler_ResetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	handleable := NewMockHandleable(mockCtrl)
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = handleable
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = handleable
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
			handleable.EXPECT().resetFaceInformation(tt.args.faceIndex).MaxTimes(2)
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
	handleable := NewMockHandleable(mockCtrl)
	h := NewHandler()
	h.lookup[reflect.TypeOf((*float32)(nil)).Elem()] = handleable
	h.lookup[reflect.TypeOf((*float64)(nil)).Elem()] = handleable
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
			handleable.EXPECT().permuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3).MaxTimes(2)
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
