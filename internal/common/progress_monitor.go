package common

import (
	"math" 
	"github.com/qmuntal/go3mf/pkg/stack"
	"github.com/qmuntal/go3mf/pkg/semaphore"
)

type ProgressMonitor struct {
	progressCallback   ProgressCallback
	userData           interface{}
	lastCallbackResult bool
	levels             *stack.ItemStack
	callbackMutex      *semaphore.Semaphore
}

func NewProgressMonitor() *ProgressMonitor {
	return &ProgressMonitor{
		lastCallbackResult: true,
		callbackMutex: semaphore.NewSemaphore(),
		levels: stack.NewItemStack(),
	}
}

func (p *ProgressMonitor) QueryCancelled() bool{
	return p.Progress(-1, ProgressQueryCanceled)
}

func (p *ProgressMonitor) Progress(progress float64, identifier ProgressIdentifier) bool{
	if p.progressCallback == nil || !p.callbackMutex.CanRun() {
		return true
	}

	var nProgress int
	if progress == -1{
		nProgress = -1
	}else {
		nProgress = int(100.0*(p.level().a + math.Max(math.Min(progress, 1.0), 0.0) * (p.level().b - p.level().a)))
	}
	p.lastCallbackResult = p.progressCallback(nProgress, identifier, p.userData)
	p.callbackMutex.Done()
	return p.lastCallbackResult
}

func (p *ProgressMonitor) PushLevel(relativeStart float64, relativeEnd float64) {
	curLevel := p.level()
	curRange := curLevel.b - curLevel.a
	p.levels.Push(Float64Pair{curLevel.a + curRange*relativeStart, curLevel.a + curRange*relativeEnd})
}

func (p *ProgressMonitor) PopLevel() Float64Pair {
	ret := p.level()
	if (!p.levels.Empty()) {
		p.levels.Pop()
	}
	return ret
}

func (p *ProgressMonitor) ResetLevels() {
	for !p.levels.Empty() {
		p.levels.Pop()
	}
}

func (p *ProgressMonitor) level() Float64Pair {
	if (p.levels.Empty()) {
		p.levels.Push(Float64Pair{0.0, 1.0})
	}
	return (*p.levels.Top()).(Float64Pair)
}

func (p *ProgressMonitor) SetProgressCallback(callback ProgressCallback, userData interface{}){
	p.progressCallback = callback
	p.userData = userData
	p.lastCallbackResult = true
	p.ResetLevels()
}

func (p *ProgressMonitor) ClearProgressCallback(){
	p.SetProgressCallback(nil, nil)
}

func (p *ProgressMonitor) WasAborted() bool {
	return p.lastCallbackResult == false
}

func (p *ProgressMonitor) GetProgressMessage(progressIdentifier ProgressIdentifier) string {
	switch (progressIdentifier) {
		case ProgressQueryCanceled: 
			return ""
		case ProgressDone: 
			return "Done"
		case ProgressCleanup: 
			return "Cleaning up"
		case ProgressReadStream: 
			return "Reading stream"
		case ProgressExtraxtOPCPackage: 
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
		case ProgressWriteNoBjects: 
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