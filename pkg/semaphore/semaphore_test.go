package semaphore

import (
	"testing"
)

func TestSemaphore_CanRun(t *testing.T) {
	s := Semaphore{}
	if !s.CanRun() {
		t.Error("wrong semaphore state, should be green")
	}
	if s.CanRun() {
		t.Error("wrong semaphore state, should be red")
	}
}

func TestSemaphore_Done(t *testing.T) {
	s := Semaphore{}
	if !s.CanRun() {
		t.Error("wrong semaphore state, should be green")
	}
	s.Done()
	if !s.CanRun() {
		t.Error("wrong semaphore state, should be green")
	}
}
