package meshinfo

// LookupMeshInformationHandler implements MeshInformationHandler.
// It allows to include different kinds of information in one mesh (like Textures AND colors).
type LookupMeshInformationHandler struct {
	informations []MeshInformation
	lookup       map[InformationType]MeshInformation
}
