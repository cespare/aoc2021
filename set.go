package main

// https://github.com/golang/go/discussions/47331

type Set[T comparable] struct {
	m map[T]struct{}
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

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

func (s *Set[T]) Filter(keep func(T) bool) {
	for k := range s.m {
		if !keep(k) {
			delete(s.m, k)
		}
	}
}

func (s *Set[T]) Pop() T {
	for k := range s.m {
		return k
	}
	panic("Pop of empty set")
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
