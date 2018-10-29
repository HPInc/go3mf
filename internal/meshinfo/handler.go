package meshinfo

import (
	"reflect"
)

const maxInternalID = 9223372036854775808

// Handler implements Handler.
// It allows to include different kinds of information in one mesh (like Textures AND colors).
type Handler struct {
	lookup            map[reflect.Type]MeshInfo
	internalIDCounter uint64
}

// NewHandler creates a new handler.
func NewHandler() *Handler {
	handler := &Handler{
		lookup:            make(map[reflect.Type]MeshInfo, 0),
		internalIDCounter: 1,
	}
	return handler
}

// AddInformation adds a new type of information to the handler.
func (h *Handler) AddInformation(info MeshInfo) error {
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
func (h *Handler) GetInformationByType(infoType reflect.Type) (MeshInfo, bool) {
	info, ok := h.lookup[infoType]
	return info, ok
}

// GetInformationCount returns the number of informations added to the handler.
func (h *Handler) GetInformationCount() uint32 {
	return uint32(len(h.lookup))
}

// AddInfoFromTable adds the information of the target handler.
func (h *Handler) AddInfoFromTable(otherHandler *Handler, currentFaceCount uint32) error {
	for infoType, otherInfo := range otherHandler.lookup {
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
func (h *Handler) CloneFaceInfosFrom(faceIndex uint32, otherHandler *Handler, otherFaceIndex uint32) {
	for infoType, otherInfo := range otherHandler.lookup {
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
