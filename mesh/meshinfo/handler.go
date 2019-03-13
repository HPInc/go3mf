package meshinfo

// DataType represents a type of data.
type DataType int

const (
	BaseMaterialType DataType = iota
	TextureCoordsType
	NodeColorType
)

// Handler allows to include specific types of information in one mesh (like Textures AND colors).
type Handler struct {
	lookup map[DataType]Handleable
}

// NewHandler creates a new handler.
func NewHandler() *Handler {
	handler := &Handler{
		lookup: make(map[DataType]Handleable, 0),
	}
	return handler
}

// AddBaseMaterialInfo adds a BaseMaterialInfo to the handler.
func (h *Handler) AddBaseMaterialInfo(currentFaceCount uint32) *FacesData {
	f := newFacesData(newbaseMaterialContainer(currentFaceCount))
	h.addInformation(f)
	return f
}

// AddTextureCoordsInfo adds a TextureCoordsInfo to the handler.
func (h *Handler) AddTextureCoordsInfo(currentFaceCount uint32) *FacesData {
	f := newFacesData(newtextureCoordsContainer(currentFaceCount))
	h.addInformation(f)
	return f
}

// AddNodeColorInfo adds a NodeColorInfo to the handler.
func (h *Handler) AddNodeColorInfo(currentFaceCount uint32) *FacesData {
	f := newFacesData(newnodeColorContainer(currentFaceCount))
	h.addInformation(f)
	return f
}

// BaseMaterialInfo returns the base material information. If it is not created the second parameters will be false.
func (h *Handler) BaseMaterialInfo() (data *FacesData, ok bool) {
	info, ok := h.lookup[BaseMaterialType]
	if ok {
		data, ok = info.(*FacesData)
	}
	return
}

// TextureCoordsInfo returns the texture coords information. If it is not created the second parameters will be false.
func (h *Handler) TextureCoordsInfo() (data *FacesData, ok bool) {
	info, ok := h.lookup[TextureCoordsType]
	if ok {
		data, ok = info.(*FacesData)
	}
	return
}

// NodeColorInfo returns the node color information. If it is not created the second parameters will be false.
func (h *Handler) NodeColorInfo() (data *FacesData, ok bool) {
	info, ok := h.lookup[NodeColorType]
	if ok {
		data, ok = info.(*FacesData)
	}
	return
}

// AddFace adds a new face to the handler.
func (h *Handler) AddFace(newFaceCount uint32) {
	for _, info := range h.lookup {
		info.AddFaceData(newFaceCount).Invalidate()
	}
}

// InformationCount returns the number of informations added to the handler.
func (h *Handler) InformationCount() uint32 {
	return uint32(len(h.lookup))
}

// AddInfoFrom adds the information of the target handler.
func (h *Handler) AddInfoFrom(informer TypedInformer, currentFaceCount uint32) {
	types := informer.InfoTypes()
	for _, infoType := range types {
		otherInfo, _ := informer.InformationByType(infoType)
		if _, ok := h.lookup[infoType]; !ok {
			h.addInformation(otherInfo.clone(currentFaceCount))
		}
	}
}

// ResetFaceInformation clears the data of an specific face.
func (h *Handler) ResetFaceInformation(faceIndex uint32) {
	for _, info := range h.lookup {
		info.resetFaceInformation(faceIndex)
	}
}

// PermuteNodeInformation swap the data of the target mesh.
func (h *Handler) PermuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	for _, info := range h.lookup {
		info.permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3)
	}
}

// RemoveAllInformations clears all the data from the handler.
func (h *Handler) RemoveAllInformations() {
	for infoType := range h.lookup {
		h.removeInformation(infoType)
	}
}

// InfoTypes returns the types of informations stored in the handler.
func (h *Handler) InfoTypes() []DataType {
	types := make([]DataType, 0, len(h.lookup))
	for infoType := range h.lookup {
		types = append(types, infoType)
	}
	return types
}

// addInformation adds a new type of information to the handler.
func (h *Handler) addInformation(info Handleable) {
	h.lookup[info.InfoType()] = info
}

// InformationByType retrieves the information of the desried type.
func (h *Handler) InformationByType(infoType DataType) (Handleable, bool) {
	info, ok := h.lookup[infoType]
	return info, ok
}

// removeInformation removes the information of the target type.
func (h *Handler) removeInformation(infoType DataType) {
	if _, ok := h.lookup[infoType]; ok {
		delete(h.lookup, infoType)
	}
}

// CopyFaceInfosFrom clones the data from another face.
func (h *Handler) CopyFaceInfosFrom(faceIndex uint32, informer TypedInformer, otherFaceIndex uint32) {
	types := informer.InfoTypes()
	for _, infoType := range types {
		otherInfo, _ := informer.InformationByType(infoType)
		info, ok := h.lookup[infoType]
		if ok {
			info.copyFaceInfosFrom(faceIndex, otherInfo, otherFaceIndex)
		}
	}
}
