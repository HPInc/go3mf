package semaphore

import (
	"reflect"
	"testing"
)

func TestNewSemaphore(t *testing.T) {
	tests := []struct {
		name string
		want *Semaphore
	}{
		{
			name: "new",
			want: &Semaphore{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSemaphore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSemaphore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemaphore_CanRun(t *testing.T) {
	s := NewSemaphore()
	if !s.CanRun() {
		t.Error("wrong semaphore state, should be green")
	}
	if s.CanRun() {
		t.Error("wrong semaphore state, should be red")
	}
}

func TestSemaphore_Done(t *testing.T) {
	s := NewSemaphore()
	if !s.CanRun() {
		t.Error("wrong semaphore state, should be green")
	}
	s.Done()
	if !s.CanRun() {
		t.Error("wrong semaphore state, should be green")
	}
}
