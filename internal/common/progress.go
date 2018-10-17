package common

type ProgressIdentifier int

const (
	ProgressQueryCanceled ProgressIdentifier = iota
	ProgressDone
	ProgressCleanup
	ProgressReadStream
	ProgressExtraxtOPCPackage
	ProgressReadNonRootModels
	ProgressReadRootModel
	ProgressReadResources
	ProgressReadMesh
	ProgressReadSlices
	ProgressReadBuild
	ProgressCreateOPCPackage
	ProgressWriteModelsToStream
	ProgressWriteRootModel
	ProgressWriteNonRootModels
	ProgressWriteAttachements
	ProgressWriteContentTypes
	ProgressWriteNoBjects
	ProgressWriteNodes
	ProgressWriteTriangles
	ProgressWriteSlices
)

type ProgressCallback func(progress int, id ProgressIdentifier, data interface{}) bool

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
