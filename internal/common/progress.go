package common

// A ProgressIdentifier is an enumerable for the different progress states
type ProgressIdentifier int

const (
	// ProgressQueryCanceled when the process has been canceled
	ProgressQueryCanceled ProgressIdentifier = iota
	// ProgressDone when the process has finished
	ProgressDone
	// ProgressCleanup when the process is cleaning up
	ProgressCleanup
	// ProgressReadStream when the process is reading an stream
	ProgressReadStream
	// ProgressExtractOPCPackage  whenthe process is extracting an OPC package
	ProgressExtractOPCPackage
	// ProgressReadNonRootModels when the process is reading non-root models
	ProgressReadNonRootModels
	// ProgressReadRootModel when the process is reading root models
	ProgressReadRootModel
	// ProgressReadResources when the process is reading resources
	ProgressReadResources
	// ProgressReadMesh when the process is reading a mesh
	ProgressReadMesh
	// ProgressReadSlices when the process is reading slices
	ProgressReadSlices
	// ProgressReadBuild when the process is reading a build
	ProgressReadBuild
	// ProgressCreateOPCPackage when the process is creating an OPC package
	ProgressCreateOPCPackage
	// ProgressWriteModelsToStream when the process is writing the models to an stream
	ProgressWriteModelsToStream
	// ProgressWriteRootModel when the process is writing the root model
	ProgressWriteRootModel
	// ProgressWriteNonRootModels when the process is writing non-root models
	ProgressWriteNonRootModels
	// ProgressWriteAttachements when the process is writting the attachements
	ProgressWriteAttachements
	// ProgressWriteContentTypes when the process is writing content types
	ProgressWriteContentTypes
	// ProgressWriteObjects when the process is writing objects
	ProgressWriteObjects
	// ProgressWriteNodes when the process is writing nodes
	ProgressWriteNodes
	// ProgressWriteTriangles when the process is writing triangles
	ProgressWriteTriangles
	// ProgressWriteSlices when the process is writing slices
	ProgressWriteSlices
)

// A Float64Pair is a tuple of two float64 values
type Float64Pair struct {
	A float64 // the first element of the tuple
	B float64 // the second element of the tuple
}

// ProgressCallback defines the signature of the callback which will be called when there is a progress in the process.
// Returns true if the progress should continue and false to abort it.
type ProgressCallback func(progress int, id ProgressIdentifier, data interface{}) bool

// Progress defines an interface that keeps track of the progress of a process.
// The implementation should manage concurrency and optionally enable a client callback.
type Progress interface {
	QueryCancelled() bool
	Progress(progress float64, identifier ProgressIdentifier) bool
	PushLevel(relativeStart float64, relativeEnd float64)
	PopLevel() Float64Pair
	ResetLevels()
	SetProgressCallback(callback ProgressCallback, userData interface{})
	ClearProgressCallback()
	WasAborted() bool
	GetProgressMessage(progressIdentifier ProgressIdentifier) string
}
