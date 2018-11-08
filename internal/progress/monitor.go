package common

import (
	"math"

	"github.com/qmuntal/go3mf/pkg/semaphore"
	"github.com/qmuntal/go3mf/pkg/stack"
)

// Monitor is the reference implementation for the Progress interface.
// It uses semaphores for managing concurrent notification and stacks for managing the process.
type Monitor struct {
	progressCallback   progressCallback
	userData           interface{}
	lastCallbackResult bool
	levels             *stack.ItemStack
	callbackMutex      *semaphore.Semaphore
}

// NewMonitor creates a default ProgresMonitor.
func NewMonitor() *Monitor {
	return &Monitor{
		lastCallbackResult: true,
		callbackMutex:      semaphore.NewSemaphore(),
		levels:             stack.NewItemStack(),
	}
}

// QueryCancelled cancels the current process with a ProgressQueryCanceled identifier.
func (p *Monitor) QueryCancelled() bool {
	return p.Progress(-1, StageQueryCanceled)
}

// Progress updates the progress of the current process.
// If the callback is nil or there is another progress being notified it does nothing and return true.
func (p *Monitor) Progress(progress float64, identifier Stage) bool {
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
func (p *Monitor) PushLevel(relativeStart float64, relativeEnd float64) {
	curLevel := p.level()
	curRange := curLevel.B - curLevel.A
	p.levels.Push(float64Pair{curLevel.A + curRange*relativeStart, curLevel.A + curRange*relativeEnd})
}

// PopLevel removes a level from the progress
func (p *Monitor) PopLevel() (a, b float64) {
	ret := p.level()
	if !p.levels.Empty() {
		p.levels.Pop()
	}
	return ret.A, ret.B
}

// ResetLevels empty the level stack
func (p *Monitor) ResetLevels() {
	for !p.levels.Empty() {
		p.levels.Pop()
	}
}

func (p *Monitor) level() float64Pair {
	if p.levels.Empty() {
		p.levels.Push(float64Pair{0.0, 1.0})
	}
	return (*p.levels.Top()).(float64Pair)
}

// SetProgressCallback restarts the progress and specifies the callback to be executed on every step of the progress.
// Optionaly usedData can be defined, which will be passed as parameter to the callback.
func (p *Monitor) SetProgressCallback(callback progressCallback, userData interface{}) {
	p.progressCallback = callback
	p.userData = userData
	p.lastCallbackResult = true
	p.ResetLevels()
}

// ClearProgressCallback restarts the process and clears the progress callback.
func (p *Monitor) ClearProgressCallback() {
	p.SetProgressCallback(nil, nil)
}

// WasAborted returns true if the callback asked for aborting the progress, false otherwise.
func (p *Monitor) WasAborted() bool {
	return p.lastCallbackResult == false
}

// ProgressMessage stringify the progress identifiers.
func (p *Monitor) ProgressMessage(progressIdentifier Stage) string {
	if val, ok := progressMap[progressIdentifier]; ok {
		return val
	}
	return "Unknown Progress Identifier"
}
