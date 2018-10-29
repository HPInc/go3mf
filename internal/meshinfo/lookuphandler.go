package meshinfo

import (
	"reflect"
	"github.com/qmuntal/go3mf/internal/common"
)

const maxInternalID = 9223372036854775808

// LookupHandler implements Handler.
// It allows to include different kinds of information in one mesh (like Textures AND colors).
type LookupHandler struct {
	lookup            map[reflect.Type]MeshInfo
	internalIDCounter uint64
}

// NewLookupHandler creates a new lookup handler.
func NewLookupHandler() *LookupHandler {
	handler := &LookupHandler{
		lookup:            make(map[reflect.Type]MeshInfo, infoLastType),
		internalIDCounter: 1,
	}
	return handler
}

// InfoTypes returns the types of informations stored in the handler.
func (h *LookupHandler) InfoTypes() []reflect.Type {
	types := make([]reflect.Type, 0, len(h.lookup))
	for infoType := range h.lookup {
		types = append(types, infoType)
	}
	return types
}

// AddInformation adds a new type of information to the handler.
func (h *LookupHandler) AddInformation(info MeshInfo) error {
	infoType := info.InfoType()
	h.lookup[infoType] = info
	info.setInternalID(h.internalIDCounter)
	h.internalIDCounter++
	if h.internalIDCounter > maxInternalID {
		return common.NewError(common.ErrorHandleOverflow)
	}
	return nil
}

// AddFace adds a new face to the handler.
func (h *LookupHandler) AddFace(newFaceCount uint32) error {
	for _, info := range h.lookup {
		data, err := info.AddFaceData(newFaceCount)
		if err != nil {
			return err
		}
		data.Invalidate()
	}
	return nil
}

// GetInformationByType retrieves the information of the desried type.
func (h *LookupHandler) GetInformationByType(infoType reflect.Type) (MeshInfo, bool) {
	info, ok := h.lookup[infoType]
	return info, ok
}

// GetInformationCount returns the number of informations added to the handler.
func (h *LookupHandler) GetInformationCount() uint32 {
	return uint32(len(h.lookup))
}

// AddInfoFromTable adds the information of the target handler.
func (h *LookupHandler) AddInfoFromTable(otherHandler Handler, currentFaceCount uint32) error {
	types := otherHandler.InfoTypes()
	for _, infoType := range types {
		otherInfo, _ := otherHandler.GetInformationByType(infoType)
		if _, ok := h.lookup[infoType]; !ok {
			err := h.AddInformation(otherInfo.Clone(currentFaceCount))
			if err != nil {
				return err
			}
		}
		h.lookup[infoType].mergeInformationFrom(otherInfo)
	}
	return nil
}

// CloneFaceInfosFrom clones the data from another face.
func (h *LookupHandler) CloneFaceInfosFrom(faceIndex uint32, otherHandler Handler, otherFaceIndex uint32) {
	types := otherHandler.InfoTypes()
	for _, infoType := range types {
		otherInfo, _ := otherHandler.GetInformationByType(infoType)
		info, ok := h.lookup[infoType]
		if ok {
			info.cloneFaceInfosFrom(faceIndex, otherInfo, otherFaceIndex)
		}
	}
}

// ResetFaceInformation clears the data of an specific face.
func (h *LookupHandler) ResetFaceInformation(faceIndex uint32) {
	for _, info := range h.lookup {
		info.resetFaceInformation(faceIndex)
	}
}

// RemoveInformation removes the information of the target type.
func (h *LookupHandler) RemoveInformation(infoType reflect.Type) {
	if _, ok := h.lookup[infoType]; ok {
		delete(h.lookup, infoType)
	}
}

// PermuteNodeInformation swap the data of the target mesh.
func (h *LookupHandler) PermuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	for _, info := range h.lookup {
		info.permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3)
	}
}
