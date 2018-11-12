//go:generate mockgen -destination types_mock_test.go -package meshinfo -self_package github.com/qmuntal/go3mf/internal/meshinfo github.com/qmuntal/go3mf/internal/meshinfo FaceData,Container,TypedInformer,FaceQuerier,Handleable

package meshinfo

import "reflect"

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
	GetFaceData(faceIndex uint32) FaceData
}

type dataCreator interface {
	// AddFaceData adds data to the last added face and returns the pointer to the data of the added face.
	// The parameter newFaceCount should indicate the faces information stored in the container, including the new one.
	AddFaceData(newFaceCount uint32) FaceData
}

// Container provides a repository for holding information state of a complete mesh structure.
// It is intended to be used as a low level repository of mesh information,
// such a thin wrapper around an in memory map or a disk serializer.
type Container interface {
	dataCreator
	FaceQuerier
	// FaceCount returns the number of faces information stored in the container.
	FaceCount() uint32
	// InfoType returns the type of the stored data.
	InfoType() reflect.Type
	// Clear removes all the information stored in the container.
	Clear()
	// clone creates a copy of the container with all the faces invalidated.
	clone(currentFaceCount uint32) Container
}

// faceModifier defines methods that can modify an inexed face.
type faceModifier interface {
	// resetFaceInformation clears the data of an specific face.
	resetFaceInformation(faceIndex uint32)
	// cloneFaceInfosFrom clones the data from another face.
	cloneFaceInfosFrom(faceIndex uint32, otherInfo FaceQuerier, otherFaceIndex uint32)
	// permuteNodeInformation swap the data of the target mesh.
	permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32)
}

// Handleable defines an interface than can be handled by Handler.
type Handleable interface {
	FaceQuerier
	faceModifier
	dataCreator
	// InfoType returns the type of the stored data.
	InfoType() reflect.Type
	// setInternalID sets an ID for the whole mesh information.
	setInternalID(internalID uint64)
	// clone creates a deep copy of this instance.
	clone(currentFaceCount uint32) Handleable
}

// TypedInformer inform about specific types of information.
type TypedInformer interface {
	// InformationByType retrieves the information of the desired type.
	InformationByType(infoType reflect.Type) (*FacesData, bool)
	// InfoTypes returns the types of informations stored in the handler.
	InfoTypes() []reflect.Type
	// informationByType retrieves the information of the desired type.
	informationByType(infoType reflect.Type) (Handleable, bool)
}
