package main

import (
	"sync"
)

func main() {

}

// NewSetRWMutex создает SetRWMutex
func NewSetRWMutex() *SetRWMutex {
	return &SetRWMutex{
		nm: make(map[int]struct{}),
	}
}

// SetRWMutex - множество, использующе для безопастного доступа sync.RWMutex
type SetRWMutex struct {
	me sync.RWMutex
	nm map[int]struct{}
}

// Add добавляет число ко множеству
func (s *SetRWMutex) Add(i int) {
	s.me.Lock()
	s.nm[i] = struct{}{}
	s.me.Unlock()
}

// Has проверяет наличие числа в множестве
func (s *SetRWMutex) Has(i int) bool {
	s.me.RLock()
	_, ok := s.nm[i]
	s.me.RUnlock()
	return ok
}

// NewSetMutex создает SetMutex
func NewSetMutex() *SetMutex {
	return &SetMutex{
		nm: make(map[int]struct{}),
	}
}

// SetMutex - множество, использующе для безопастного доступа sync.Mutex
type SetMutex struct {
	me sync.Mutex
	nm map[int]struct{}
}

// Add добавляет число ко множеству
func (s *SetMutex) Add(i int) {
	s.me.Lock()
	s.nm[i] = struct{}{}
	s.me.Unlock()
}

// Has проверяет наличие числа в множестве
func (s *SetMutex) Has(i int) bool {
	s.me.Lock()
	_, ok := s.nm[i]
	s.me.Unlock()
	return ok
}
