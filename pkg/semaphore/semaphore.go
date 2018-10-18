package semaphore

import "sync/atomic"

// Semaphore is used to check if a code block is being executed by another goroutine.
type Semaphore struct {
	semaphore int32
}

// NewSemaphore creates a new default semaphore.
func NewSemaphore() *Semaphore {
	return &Semaphore{}
}

// CanRun should be called when entering to the synchronization block.
// Returns true if no thread is running inside the block and false otherwhise.
// If the result is true the caller is responsible for freeing the block with a Done().
func (l *Semaphore) CanRun() bool {
	return atomic.CompareAndSwapInt32(&l.semaphore, 0, 1)
}

// Done should be called to free a synchronization block after calling CanRun() with a true result.
func (l *Semaphore) Done() {
	atomic.CompareAndSwapInt32(&l.semaphore, 1, 0)
}
