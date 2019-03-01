package stack

import (
	"testing"
)

func TestPush(t *testing.T) {
	s := ItemStack{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if size := len(s.items); size != 3 {
		t.Errorf("wrong count, expected 3 and got %d", size)
	}
}

func TestPop(t *testing.T) {
	s := ItemStack{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Pop()
	if size := len(s.items); size != 2 {
		t.Errorf("wrong count, expected 2 and got %d", size)
	}

	s.Pop()
	s.Pop()
	if size := len(s.items); size != 0 {
		t.Errorf("wrong count, expected 0 and got %d", size)
	}
}

func TestEmpty(t *testing.T) {
	s := ItemStack{}
	if !s.Empty() {
		t.Errorf("expected to be empty, got %d", len(s.items))
	}
	s.Push(1)
	if s.Empty() {
		t.Error("expected not to be empty")
	}
}

func TestTop(t *testing.T) {
	s := ItemStack{}
	s.Push(1)
	r := *s.Top()
	if r != 1 {
		t.Errorf("expected top to be 1, got %d", r)
	}
	s.Push(2)
	r = *s.Top()
	if r != 2 {
		t.Errorf("expected top to be 2, got %d", r)
	}
}
