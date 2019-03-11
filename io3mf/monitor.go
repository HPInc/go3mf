package io3mf

import (
	"math"

	"github.com/qmuntal/go3mf/pkg/semaphore"
	"github.com/qmuntal/go3mf/pkg/stack"
)

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

// A float64Pair is a tuple of two float64 values
type float64Pair struct {
	A float64 // the first element of the tuple
	B float64 // the second element of the tuple
}

// ProgressCallback defines the signature of the callback which will be called when there is a progress in the process.
// Returns true if the progress should continue and false to abort it.
type ProgressCallback func(progress int, id Stage, data interface{}) bool

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

// monitor is the reference implementation for the Progress interface.
// It uses semaphores for managing concurrent notification and stacks for managing the process.
type monitor struct {
	progressCallback    ProgressCallback
	userData            interface{}
	lastCallbackAborted bool
	levels              stack.ItemStack
	callbackMutex       semaphore.Semaphore
}

// QueryCancelled cancels the current process with a ProgressQueryCanceled identifier.
func (p *monitor) QueryCancelled() bool {
	return p.progress(-1, StageQueryCanceled)
}

// Progress updates the progress of the current process.
// If the callback is nil or there is another progress being notified it does nothing and return true.
func (p *monitor) progress(progress float64, identifier Stage) bool {
	if p.progressCallback == nil || !p.callbackMutex.CanRun() {
		return true
	}

	var nProgress int
	if progress == -1 {
		nProgress = -1
	} else {
		nProgress = int(100.0 * (p.level().A + math.Max(math.Min(progress, 1.0), 0.0)*(p.level().B-p.level().A)))
	}
	p.lastCallbackAborted = !p.progressCallback(nProgress, identifier, p.userData)
	p.callbackMutex.Done()
	return p.lastCallbackAborted == false
}

func (p *monitor) pushLevel(relativeStart float64, relativeEnd float64) {
	curLevel := p.level()
	curRange := curLevel.B - curLevel.A
	p.levels.Push(float64Pair{curLevel.A + curRange*relativeStart, curLevel.A + curRange*relativeEnd})
}

// popLevel removes a level from the progress
func (p *monitor) popLevel() (a, b float64) {
	ret := p.level()
	if !p.levels.Empty() {
		p.levels.Pop()
	}
	return ret.A, ret.B
}

// ResetLevels empty the level stack
func (p *monitor) ResetLevels() {
	for !p.levels.Empty() {
		p.levels.Pop()
	}
}

func (p *monitor) level() float64Pair {
	if p.levels.Empty() {
		p.levels.Push(float64Pair{0.0, 1.0})
	}
	return (*p.levels.Top()).(float64Pair)
}

// SetProgressCallback restarts the progress and specifies the callback to be executed on every step of the progress.
// Optionaly usedData can be defined, which will be passed as parameter to the callback.
func (p *monitor) SetProgressCallback(callback ProgressCallback, userData interface{}) {
	p.progressCallback = callback
	p.userData = userData
	p.lastCallbackAborted = false
	p.ResetLevels()
}

// ClearProgressCallback restarts the process and clears the progress callback.
func (p *monitor) ClearProgressCallback() {
	p.SetProgressCallback(nil, nil)
}

// WasAborted returns true if the callback asked for aborting the progress, false otherwise.
func (p *monitor) WasAborted() bool {
	return p.lastCallbackAborted
}

// ProgressMessage stringify the progress identifiers.
func (p *monitor) ProgressMessage(progressIdentifier Stage) string {
	if val, ok := progressMap[progressIdentifier]; ok {
		return val
	}
	return "Unknown Progress Identifier"
}
