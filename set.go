package main

import "fmt"

// https://github.com/golang/go/discussions/47331

type Set[T comparable] struct {
	m map[T]struct{}
}

func SetOf[T comparable](v ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(v))}
	s.Add(v...)
	return s
}

func (s *Set[T]) String() string {
	return fmt.Sprint(s.m)
}

func (s *Set[T]) Len() int {
	return len(s.m)
}

func (s *Set[T]) Add(v ...T) {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	for _, vv := range v {
		s.m[vv] = struct{}{}
	}
}

func (s *Set[T]) AddSet(s2 *Set[T]) {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	for k := range s2.m {
		s.m[k] = struct{}{}
	}
}

func (s *Set[T]) Remove(v T) {
	delete(s.m, v)
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

func (s *Set[T]) Do(f func(T) bool) {
	for v := range s.m {
		if !f(v) {
			return
		}
	}
}

func (s *Set[T]) Filter(keep func(T) bool) {
	for k := range s.m {
		if !keep(k) {
			delete(s.m, k)
		}
	}
}

func (s *Set[T]) Pop() (T, bool) {
	for k := range s.m {
		delete(s.m, k)
		return k, true
	}
	var t T
	return t, false
}

func (s *Set[T]) Slice() []T {
	var ts []T
	for k := range s.m {
		ts = append(ts, k)
	}
	return ts
}

func SetUnion[T comparable](s0, s1 Set[T]) Set[T] {
	var s2 Set[T]
	for k := range s0.m {
		s2.Add(k)
	}
	for k := range s1.m {
		s2.Add(k)
	}
	return s2
}

func SetIntersection[T comparable](s0, s1 Set[T]) Set[T] {
	var s2 Set[T]
	for k := range s0.m {
		if s1.Contains(k) {
			s2.Add(k)
		}
	}
	return s2
}
