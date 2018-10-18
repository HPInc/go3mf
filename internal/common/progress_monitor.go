package common

import (
	"github.com/qmuntal/go3mf/pkg/semaphore"
	"github.com/qmuntal/go3mf/pkg/stack"
	"math"
)

// ProgressMonitor is the reference implementation for the Progress interface.
// It uses semaphores for managing concurrent notification and stacks for managing the process.
type ProgressMonitor struct {
	progressCallback   ProgressCallback
	userData           interface{}
	lastCallbackResult bool
	levels             *stack.ItemStack
	callbackMutex      *semaphore.Semaphore
}

// NewProgressMonitor creates a default ProgresMonitor.
func NewProgressMonitor() *ProgressMonitor {
	return &ProgressMonitor{
		lastCallbackResult: true,
		callbackMutex:      semaphore.NewSemaphore(),
		levels:             stack.NewItemStack(),
	}
}

// QueryCancelled cancels the current process with a ProgressQueryCanceled identifier.
func (p *ProgressMonitor) QueryCancelled() bool {
	return p.Progress(-1, ProgressQueryCanceled)
}

// Progress updates the progress of the current process.
// If the callback is nil or there is another progress being notified it does nothing and return true.
func (p *ProgressMonitor) Progress(progress float64, identifier ProgressIdentifier) bool {
	if p.progressCallback == nil || !p.callbackMutex.CanRun() {
		return true
	}

	var nProgress int
	if progress == -1 {
		nProgress = -1
	} else {
		nProgress = int(100.0 * (p.level().A + math.Max(math.Min(progress, 1.0), 0.0)*(p.level().B-p.level().A)))
	}
	p.lastCallbackResult = p.progressCallback(nProgress, identifier, p.userData)
	p.callbackMutex.Done()
	return p.lastCallbackResult
}

// PushLevel adds a new level to the progress
func (p *ProgressMonitor) PushLevel(relativeStart float64, relativeEnd float64) {
	curLevel := p.level()
	curRange := curLevel.B - curLevel.A
	p.levels.Push(Float64Pair{curLevel.A + curRange*relativeStart, curLevel.A + curRange*relativeEnd})
}

// PopLevel removes a level from the progress
func (p *ProgressMonitor) PopLevel() Float64Pair {
	ret := p.level()
	if !p.levels.Empty() {
		p.levels.Pop()
	}
	return ret
}

// ResetLevels empty the level stack
func (p *ProgressMonitor) ResetLevels() {
	for !p.levels.Empty() {
		p.levels.Pop()
	}
}

func (p *ProgressMonitor) level() Float64Pair {
	if p.levels.Empty() {
		p.levels.Push(Float64Pair{0.0, 1.0})
	}
	return (*p.levels.Top()).(Float64Pair)
}

// SetProgressCallback restarts the progress and specifies the callback to be executed on every step of the progress.
// Optionaly usedData can be defined, which will be passed as parameter to the callback.
func (p *ProgressMonitor) SetProgressCallback(callback ProgressCallback, userData interface{}) {
	p.progressCallback = callback
	p.userData = userData
	p.lastCallbackResult = true
	p.ResetLevels()
}

// ClearProgressCallback restarts the process and clears the progress callback.
func (p *ProgressMonitor) ClearProgressCallback() {
	p.SetProgressCallback(nil, nil)
}

// WasAborted returns true if the callback asked for aborting the progress, false otherwise.
func (p *ProgressMonitor) WasAborted() bool {
	return p.lastCallbackResult == false
}

// GetProgressMessage stringify the progress identifiers.
func (p *ProgressMonitor) GetProgressMessage(progressIdentifier ProgressIdentifier) string {
	switch progressIdentifier {
	case ProgressQueryCanceled:
		return ""
	case ProgressDone:
		return "Done"
	case ProgressCleanup:
		return "Cleaning up"
	case ProgressReadStream:
		return "Reading stream"
	case ProgressExtractOPCPackage:
		return "Extracting OPC package"
	case ProgressReadNonRootModels:
		return "Reading non-root models"
	case ProgressReadRootModel:
		return "Reading root model"
	case ProgressReadResources:
		return "Reading resources"
	case ProgressReadMesh:
		return "Reading mesh data"
	case ProgressReadSlices:
		return "Reading slice data"
	case ProgressReadBuild:
		return "Reading build definition"
	case ProgressCreateOPCPackage:
		return "Creating OPC package"
	case ProgressWriteModelsToStream:
		return "Writing models to stream"
	case ProgressWriteRootModel:
		return "Writing root model"
	case ProgressWriteNonRootModels:
		return "Writing non-root models"
	case ProgressWriteAttachements:
		return "Writing attachments"
	case ProgressWriteContentTypes:
		return "Writing content types"
	case ProgressWriteObjects:
		return "Writing objects"
	case ProgressWriteNodes:
		return "Writing Nodes"
	case ProgressWriteTriangles:
		return "Writing triangles"
	case ProgressWriteSlices:
		return "Writing slices"
	default:
		return "Unknown Progress Identifier"
	}
}
