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
	genericHandler
}

// NewHandler creates a new handler.
func NewHandler() *Handler {
	return &Handler{
		genericHandler: *newgenericHandler(),
	}
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
func (h *Handler) BaseMaterialInfo() (*FacesData, bool) {
	info, ok := h.lookup[BaseMaterialType]
	return info.(*FacesData), ok
}

// TextureCoordsInfo returns the texture coords information. If it is not created the second parameters will be false.
func (h *Handler) TextureCoordsInfo() (*FacesData, bool) {
	info, ok := h.lookup[TextureCoordsType]
	return info.(*FacesData), ok
}

// NodeColorInfo returns the node color information. If it is not created the second parameters will be false.
func (h *Handler) NodeColorInfo() (*FacesData, bool) {
	info, ok := h.lookup[NodeColorType]
	return info.(*FacesData), ok
}
