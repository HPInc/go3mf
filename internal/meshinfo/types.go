//go:generate mockgen -destination types_mock_test.go -package meshinfo -self_package github.com/qmuntal/go3mf/internal/meshinfo github.com/qmuntal/go3mf/internal/meshinfo FaceData,Container,MeshInfo,TypedInformer,FaceQuerier

package meshinfo

import "reflect"

// Color represents a RGB color.
type Color = uint32

// FaceData defines a generic information of a face.
type FaceData interface {
	// Invalidate sets the data with its default values.
	Invalidate()
	// Copy copies all the properties from another object. Do nothing if not the same type.
	Copy(from interface{})
	// Permute some properties of the instance using the desired indexes.
	// Not all the indexes have to be used.
	Permute(index1, index2, index3 uint32)
	// Merge merges both instances.
	Merge(other interface{})
	// HasData returns true if the instances has any kind of data other than its default ones.
	HasData() bool
}

// FaceQuerier can query information for an indexed face.
type FaceQuerier interface {
	// GetFaceData returns the data of the face with the target index.
	GetFaceData(faceIndex uint32) (val FaceData, err error)
}

// Repository defines an interface for interacting with a mesh information repository.
type Repository interface {
	FaceQuerier
	// AddFaceData adds data to the last added face and returns the pointer to the data of the added face.
	// The parameter newFaceCount should indicate the faces information stored in the container, including the new one.
	// If the count is not equal to the one returned by GetCurrentFaceCount an error will be returned.
	AddFaceData(newFaceCount uint32) (val FaceData, err error)
	// GetCurrentFaceCount returns the number of faces information stored in the container.
	GetCurrentFaceCount() uint32
	// InfoType returns the type of the stored data.
	InfoType() reflect.Type
	// Clear removes all the information stored in the container.
	Clear()
}

// Container provides a repository for holding information state of a complete mesh structure.
// It is intended to be used as a low level repository of mesh information,
// such a thin wrapper around an in memory map or a disk serializer.
type Container interface {
	Repository
	// clone creates a copy of the container with all the faces invalidated.
	clone(currentFaceCount uint32) Container
}

// FaceModifier defines methods that can modify an inexed face.
type FaceModifier interface {
	// resetFaceInformation clears the data of an specific face.
	resetFaceInformation(faceIndex uint32)
	// cloneFaceInfosFrom clones the data from another face.
	cloneFaceInfosFrom(faceIndex uint32, otherInfo FaceQuerier, otherFaceIndex uint32)
	// permuteNodeInformation swap the data of the target mesh.
	permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32)
}

// Identificator defines methods to set and get an identification.
type Identificator interface {
	// setInternalID sets an ID for the whole mesh information.
	setInternalID(internalID uint64)
	// getInternalID gets the internal ID of the mesh information.
	getInternalID() uint64
}

// MeshInfo defines the Mesh Information Class.
// This is a base struct for handling all the mesh-related linear information (like face colors, textures, etc...).
type MeshInfo interface {
	Repository
	FaceModifier
	Identificator
	// clone creates a deep copy of this instance.
	clone(currentFaceCount uint32) MeshInfo
}

// TypedInformer inform about specific types of information.
type TypedInformer interface {
	// GetInformationByType retrieves the information of the desired type.
	GetInformationByType(infoType reflect.Type) (*GenericMeshInfo, bool)
	// InfoTypes returns the types of informations stored in the handler.
	InfoTypes() []reflect.Type
	// getInformationByType retrieves the information of the desired type.
	getInformationByType(infoType reflect.Type) (MeshInfo, bool)
}
