package meshinfo

import (
	"github.com/qmuntal/go3mf/internal/common"
)

const maxInternalID = 9223372036854775808

// LookupHandler implements Handler.
// It allows to include different kinds of information in one mesh (like Textures AND colors).
type LookupHandler struct {
	lookup            map[InformationType]MeshInfo
	internalIDCounter uint64
}

// NewLookupHandler creates a new lookup handler.
func NewLookupHandler() *LookupHandler {
	handler := &LookupHandler{
		lookup:            make(map[InformationType]MeshInfo, infoLastType),
		internalIDCounter: 1,
	}
	for infoType := InfoAbstract; infoType < infoLastType; infoType++ {
		handler.lookup[infoType] = nil
	}
	return handler
}

// AddInformation adds a new type of information to the handler.
func (h *LookupHandler) AddInformation(info MeshInfo) error {
	infoType := info.GetType()
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
		info.Invalidate(data)
	}
	return nil
}

// GetInformationByType retrieves the information of the desried type.
func (h *LookupHandler) GetInformationByType(infoType InformationType) MeshInfo {
	return h.lookup[infoType]
}

// GetInformationCount returns the number of informations added to the handler.
func (h *LookupHandler) GetInformationCount() uint32 {
	count := 0
	for _, info := range h.lookup {
		if info != nil {
			count++
		}
	}
	return uint32(count)
}

// AddInfoFromTable adds the information of the target handler.
func (h *LookupHandler) AddInfoFromTable(otherHandler Handler, currentFaceCount uint32) error {
	for infoType := InfoAbstract; infoType < infoLastType; infoType++ {
		otherInfo := otherHandler.GetInformationByType(infoType)
		if otherInfo != nil {
			if h.lookup[infoType] == nil {
				err := h.AddInformation(otherInfo.Clone(currentFaceCount))
				if err != nil {
					return err
				}
				h.lookup[infoType].mergeInformationFrom(otherInfo)
			}
		}
	}
	return nil
}

// CloneFaceInfosFrom clones the data from another face.
func (h *LookupHandler) CloneFaceInfosFrom(faceIndex uint32, otherHandler Handler, otherFaceIndex uint32) {
	for infoType := InfoAbstract; infoType < infoLastType; infoType++ {
		otherInfo := otherHandler.GetInformationByType(infoType)
		info := h.lookup[infoType]
		if (otherInfo != nil) && (info != nil) {
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
func (h *LookupHandler) RemoveInformation(infoType InformationType) {
	h.lookup[infoType] = nil
}

// PermuteNodeInformation swap the data of the target mesh.
func (h *LookupHandler) PermuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	for _, info := range h.lookup {
		info.permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3)
	}
}
