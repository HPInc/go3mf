package stack

// https://flaviocopes.com/golang-data-structure-stack/

import (
	"sync"
)

// Item the type of the stack
type Item interface{}

// ItemStack the stack of Items
type ItemStack struct {
	items []Item
	lock  sync.RWMutex
}

// NewItemStack creates a new ItemStack
func NewItemStack() *ItemStack {
	return &ItemStack{
		items: []Item{},
	}
}

// Push adds an Item to the top of the stack
func (s *ItemStack) Push(t Item) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Pop removes an Item from the top of the stack
func (s *ItemStack) Pop() *Item {
	s.lock.Lock()
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	s.lock.Unlock()
	return &item
}

// Empty removes all the elements from the stack
func (s *ItemStack) Empty() bool {
	return len(s.items) == 0
}

// Top returns the last added element from the stack
func (s *ItemStack) Top() *Item {
	return &s.items[len(s.items)-1]
}
