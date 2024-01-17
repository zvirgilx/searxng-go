package util

import (
	"sync"
)

type Set struct {
	store sync.Map
}

var emptyStruct struct{}

func NewSet() *Set {
	return &Set{store: sync.Map{}}
}

func SetAdd[T int | string](s *Set, value T) {
	_, ok := s.store.Load(value)
	if ok {
		return
	}
	s.store.Store(value, emptyStruct)
}

func SetMerge[T int | string](s1 *Set, s2 *Set) {
	if s2 == nil {
		return
	}
	s2.store.Range(func(key, value any) bool {
		SetAdd(s1, key.(T))
		return true
	})
}

func SetRemove[T int | string](s *Set, value T) {
	s.store.Delete(value)
}

func SetToArray[T int | string](s *Set) []T {
	set := make([]T, 0)
	if s == nil {
		return set
	}
	s.store.Range(func(key, value interface{}) bool {
		set = append(set, key.(T))
		return true
	})
	return set
}
