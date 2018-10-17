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
	ProgressReadRresources
	ProgressReadMesh
	ProgressReadSlices
	ProgressReadBuild
	ProgressCreateOPCPackage
	ProgressWriteModelsToStram
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