package meshinfo

import (
	"reflect"
)

const maxInternalID = 9223372036854775808

// Handler allows to include different kinds of information in one mesh (like Textures AND colors).
type Handler struct {
	lookup            map[reflect.Type]MeshInfo
	internalIDCounter uint64
}

// NewHandler creates a new lookup handler.
func NewHandler() *Handler {
	handler := &Handler{
		lookup:            make(map[reflect.Type]MeshInfo, 0),
		internalIDCounter: 1,
	}
	return handler
}

// InfoTypes returns the types of informations stored in the handler.
func (h *Handler) InfoTypes() []reflect.Type {
	types := make([]reflect.Type, 0, len(h.lookup))
	for infoType := range h.lookup {
		types = append(types, infoType)
	}
	return types
}

// AddInformation adds a information to the handler.
func (h *Handler) AddInformation(info *GenericMeshInfo) error {
	return h.addInformation(info)
}

// addInformation adds a new type of information to the handler.
func (h *Handler) addInformation(info MeshInfo) error {
	infoType := info.InfoType()
	h.lookup[infoType] = info
	info.setInternalID(h.internalIDCounter)
	h.internalIDCounter++
	if h.internalIDCounter > maxInternalID {
		return new(HandlerOverflowError)
	}
	return nil
}

// AddFace adds a new face to the handler.
func (h *Handler) AddFace(newFaceCount uint32) error {
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
func (h *Handler) GetInformationByType(infoType reflect.Type) (*GenericMeshInfo, bool) {
	info, ok := h.lookup[infoType]
	return info.(*GenericMeshInfo), ok
}

// getInformationByType retrieves the information of the desried type.
func (h *Handler) getInformationByType(infoType reflect.Type) (MeshInfo, bool) {
	info, ok := h.lookup[infoType]
	return info, ok
}

// GetInformationCount returns the number of informations added to the handler.
func (h *Handler) GetInformationCount() uint32 {
	return uint32(len(h.lookup))
}

// AddInfoFromTable adds the information of the target handler.
func (h *Handler) AddInfoFromTable(informer TypedInformer, currentFaceCount uint32) error {
	types := informer.InfoTypes()
	for _, infoType := range types {
		otherInfo, _ := informer.getInformationByType(infoType)
		if _, ok := h.lookup[infoType]; !ok {
			err := h.addInformation(otherInfo.Clone(currentFaceCount))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CloneFaceInfosFrom clones the data from another face.
func (h *Handler) CloneFaceInfosFrom(faceIndex uint32, informer TypedInformer, otherFaceIndex uint32) {
	types := informer.InfoTypes()
	for _, infoType := range types {
		otherInfo, _ := informer.getInformationByType(infoType)
		info, ok := h.lookup[infoType]
		if ok {
			info.cloneFaceInfosFrom(faceIndex, otherInfo, otherFaceIndex)
		}
	}
}

// ResetFaceInformation clears the data of an specific face.
func (h *Handler) ResetFaceInformation(faceIndex uint32) {
	for _, info := range h.lookup {
		info.resetFaceInformation(faceIndex)
	}
}

// RemoveInformation removes the information of the target type.
func (h *Handler) RemoveInformation(infoType reflect.Type) {
	if _, ok := h.lookup[infoType]; ok {
		delete(h.lookup, infoType)
	}
}

// PermuteNodeInformation swap the data of the target mesh.
func (h *Handler) PermuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	for _, info := range h.lookup {
		info.permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3)
	}
}
