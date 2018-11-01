package common

// A Stage is an enumerable for the different progress states
type Stage int

const (
	// StageQueryCanceled when the process has been canceled
	StageQueryCanceled Stage = iota
	// StageDone when the process has finished
	StageDone
	// StageCleanup when the process is cleaning up
	StageCleanup
	// StageReadStream when the process is reading an stream
	StageReadStream
	// StageExtractOPCPackage  whenthe process is extracting an OPC package
	StageExtractOPCPackage
	// StageReadNonRootModels when the process is reading non-root models
	StageReadNonRootModels
	// StageReadRootModel when the process is reading root models
	StageReadRootModel
	// StageReadResources when the process is reading resources
	StageReadResources
	// StageReadMesh when the process is reading a mesh
	StageReadMesh
	// StageReadSlices when the process is reading slices
	StageReadSlices
	// StageReadBuild when the process is reading a build
	StageReadBuild
	// StageCreateOPCPackage when the process is creating an OPC package
	StageCreateOPCPackage
	// StageWriteModelsToStream when the process is writing the models to an stream
	StageWriteModelsToStream
	// StageWriteRootModel when the process is writing the root model
	StageWriteRootModel
	// StageWriteNonRootModels when the process is writing non-root models
	StageWriteNonRootModels
	// StageWriteAttachements when the process is writing the attachements
	StageWriteAttachements
	// StageWriteContentTypes when the process is writing content types
	StageWriteContentTypes
	// StageWriteObjects when the process is writing objects
	StageWriteObjects
	// StageWriteNodes when the process is writing nodes
	StageWriteNodes
	// StageWriteTriangles when the process is writing triangles
	StageWriteTriangles
	// StageWriteSlices when the process is writing slices
	StageWriteSlices
)

// A Float64Pair is a tuple of two float64 values
type Float64Pair struct {
	A float64 // the first element of the tuple
	B float64 // the second element of the tuple
}

// progressCallback defines the signature of the callback which will be called when there is a progress in the process.
// Returns true if the progress should continue and false to abort it.
type progressCallback func(progress int, id Stage, data interface{}) bool

// progressMap defines a mapping between the progress identifiers and the message
var progressMap = map[Stage]string{
	StageQueryCanceled:       "",
	StageDone:                "Done",
	StageCleanup:             "Cleaning up",
	StageReadStream:          "Reading stream",
	StageExtractOPCPackage:   "Extracting OPC package",
	StageReadNonRootModels:   "Reading non-root models",
	StageReadRootModel:       "Reading root model",
	StageReadResources:       "Reading resources",
	StageReadMesh:            "Reading mesh data",
	StageReadSlices:          "Reading slice data",
	StageReadBuild:           "Reading build definition",
	StageCreateOPCPackage:    "Creating OPC package",
	StageWriteModelsToStream: "Writing models to stream",
	StageWriteRootModel:      "Writing root model",
	StageWriteNonRootModels:  "Writing non-root models",
	StageWriteAttachements:   "Writing attachments",
	StageWriteContentTypes:   "Writing content types",
	StageWriteObjects:        "Writing objects",
	StageWriteNodes:          "Writing Nodes",
	StageWriteTriangles:      "Writing triangles",
	StageWriteSlices:         "Writing slices",
}
