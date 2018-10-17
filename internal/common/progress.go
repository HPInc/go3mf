package common

import (
	"math" 
	"github.com/qmuntal/go3mf/pkg/stack"
	"github.com/qmuntal/go3mf/pkg/semaphore"
)

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
	}
}

func (p *ProgressMonitor) QueryCancelled() bool{
	return p.Progress(-1, ProgressQueryCanceled);
}

func (p *ProgressMonitor) Progress(progress float64, identifier ProgressIdentifier) bool{
	if p.progressCallback == nil || !p.callbackMutex.CanRun() {
		return true
	}

	var nProgress int
	if progress == -1{
		nProgress = -1
	}else {
		nProgress = int(100.0*(p.Level().a + math.Max(math.Min(progress, 1.0), 0.0) * (p.Level().b - p.Level().a)))
	}
	p.lastCallbackResult = p.progressCallback(nProgress, identifier, p.userData)
	p.callbackMutex.Done()
	return p.lastCallbackResult;
}

func (p *ProgressMonitor) PushLevel(relativeStart float64, relativeEnd float64) {
	curLevel := p.Level()
	curRange := curLevel.b - curLevel.a;
	p.levels.Push(Float64Pair{curLevel.a + curRange*relativeStart, curLevel.a + curRange*relativeEnd})
}

func (p *ProgressMonitor) PopLevel() Float64Pair {
	ret := p.Level();
	if (!p.levels.Empty()) {
		p.levels.Pop();
	}
	return ret;
}

func (p *ProgressMonitor) ResetLevels() {
	for !p.levels.Empty() {
		p.levels.Pop();
	}
}

func (p *ProgressMonitor) Level() Float64Pair {
	if (p.levels.Empty()) {
		p.levels.Push(Float64Pair{0.0, 1.0});
	}
	return (*p.levels.Top()).(Float64Pair);
}

func (p *ProgressMonitor) SetProgressCallback(callback ProgressCallback, userData interface{}){
	p.progressCallback = callback;
	p.userData = userData;
	p.lastCallbackResult = true;
	p.ResetLevels();
}

func (p *ProgressMonitor) ClearProgressCallback(){
	p.SetProgressCallback(nil, nil);
}

func (p *ProgressMonitor) WasAborted() bool {
	return p.lastCallbackResult == false;
}

func (p *ProgressMonitor) GetProgressMessage(progressIdentifier ProgressIdentifier, progressString *string) {
	switch (progressIdentifier) {
		case ProgressQueryCanceled: 
			*progressString = "";
		case ProgressDone: 
			*progressString = "Done";
		case ProgressCleanup: 
			*progressString = "Cleaning up";
		case ProgressReadStream: 
			*progressString = "Reading stream";
		case ProgressExtraxtOPCPackage: 
			*progressString = "Extracting OPC package";
		case ProgressReadNonRootModels: 
			*progressString = "Reading non-root models";
		case ProgressReadRootModel:
			 *progressString = "Reading root model";
		case ProgressReadResources: 
			*progressString = "Reading resources";
		case ProgressReadMesh: 
			*progressString = "Reading mesh data";
		case ProgressReadSlices: 
			*progressString = "Reading slice data";
		case ProgressReadBuild: 
			*progressString = "Reading build definition";
		case ProgressCreateOPCPackage: 
			*progressString = "Creating OPC package";
		case ProgressWriteModelsToStream: 
			*progressString = "Writing models to stream";
		case ProgressWriteRootModel: 
			*progressString = "Writing root model";
		case ProgressWriteNonRootModels: 
			*progressString = "Writing non-root models";
		case ProgressWriteAttachements: 
			*progressString = "Writing attachments";
		case ProgressWriteContentTypes: 
			*progressString = "Writing content types";
		case ProgressWriteNoBjects: 
			*progressString = "Writing objects";
		case ProgressWriteNodes: 
			*progressString = "Writing Nodes";
		case ProgressWriteTriangles: 
			*progressString = "Writing triangles";
		case ProgressWriteSlices: 
			*progressString = "Writing slices";
		default: 
			*progressString = "Unknown Progress Identifier";
	}
}