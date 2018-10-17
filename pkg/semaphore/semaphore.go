package semaphore

import "sync/atomic"

type Semaphore struct {
	semaphore int32
}

func NewSemaphore() *Semaphore {
	return &Semaphore{}
}

func (l *Semaphore) CanRun() bool {
	return atomic.CompareAndSwapInt32(&l.semaphore, 0, 1)
}
func (l *Semaphore) Done() {
	atomic.CompareAndSwapInt32(&l.semaphore, 1, 0)
}
