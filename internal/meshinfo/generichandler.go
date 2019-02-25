package meshinfo

const maxInternalID = 9223372036854775808

// genericHandler allows to include different kinds of information in one mesh (like Textures AND colors).
type genericHandler struct {
	lookup            map[DataType]Handleable
	internalIDCounter uint64
}

// newGenericHandler creates a new generic handler.
func newgenericHandler() *genericHandler {
	handler := &genericHandler{
		lookup:            make(map[DataType]Handleable, 0),
		internalIDCounter: 1,
	}
	return handler
}

// AddFace adds a new face to the handler.
func (h *genericHandler) AddFace(newFaceCount uint32) {
	for _, info := range h.lookup {
		data := info.AddFaceData(newFaceCount)
		data.Invalidate()
	}
}

// InformationCount returns the number of informations added to the handler.
func (h *genericHandler) InformationCount() uint32 {
	return uint32(len(h.lookup))
}

// AddInfoFrom adds the information of the target handler.
func (h *genericHandler) AddInfoFrom(informer TypedInformer, currentFaceCount uint32) {
	types := informer.InfoTypes()
	for _, infoType := range types {
		otherInfo, _ := informer.InformationByType(infoType)
		if _, ok := h.lookup[infoType]; !ok {
			h.addInformation(otherInfo.clone(currentFaceCount))
		}
	}
}

// ResetFaceInformation clears the data of an specific face.
func (h *genericHandler) ResetFaceInformation(faceIndex uint32) {
	for _, info := range h.lookup {
		info.resetFaceInformation(faceIndex)
	}
}

// PermuteNodeInformation swap the data of the target mesh.
func (h *genericHandler) PermuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	for _, info := range h.lookup {
		info.permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3)
	}
}

// RemoveAllInformations clears all the data from the handler.
func (h *genericHandler) RemoveAllInformations() {
	for infoType := range h.lookup {
		h.removeInformation(infoType)
	}
}

// InfoTypes returns the types of informations stored in the handler.
func (h *genericHandler) InfoTypes() []DataType {
	types := make([]DataType, 0, len(h.lookup))
	for infoType := range h.lookup {
		types = append(types, infoType)
	}
	return types
}

// addInformation adds a new type of information to the handler.
func (h *genericHandler) addInformation(info Handleable) {
	infoType := info.InfoType()
	h.lookup[infoType] = info
	info.setInternalID(h.internalIDCounter)
	h.internalIDCounter++
	if h.internalIDCounter > maxInternalID {
		panic(new(HandlerOverflowError))
	}
}

// InformationByType retrieves the information of the desried type.
func (h *genericHandler) InformationByType(infoType DataType) (Handleable, bool) {
	info, ok := h.lookup[infoType]
	return info, ok
}

// removeInformation removes the information of the target type.
func (h *genericHandler) removeInformation(infoType DataType) {
	if _, ok := h.lookup[infoType]; ok {
		delete(h.lookup, infoType)
	}
}

// CopyFaceInfosFrom clones the data from another face.
func (h *genericHandler) CopyFaceInfosFrom(faceIndex uint32, informer TypedInformer, otherFaceIndex uint32) {
	types := informer.InfoTypes()
	for _, infoType := range types {
		otherInfo, _ := informer.InformationByType(infoType)
		info, ok := h.lookup[infoType]
		if ok {
			info.copyFaceInfosFrom(faceIndex, otherInfo, otherFaceIndex)
		}
	}
}
