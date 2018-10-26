package meshinfo

// LookupHandler implements Handler.
// It allows to include different kinds of information in one mesh (like Textures AND colors).
type LookupHandler struct {
	informations []MeshInfo
	lookup       map[InformationType]MeshInfo
}
