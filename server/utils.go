package server

import "sync"

type safeList[T any] struct {
	mu   sync.Mutex
	list []T
}

func NewSafeList[T any]() safeList[T] {
	return safeList[T]{}
}

func (l *safeList[T]) Append(item T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.list = append(l.list, item)
}
func (l *safeList[T]) ToList() []T {
	return l.list
}
