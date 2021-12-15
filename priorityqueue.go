package main

import (
	"container/heap"
)

type PriorityQueue[E any] struct {
	s sliceHeap[E]
}

func NewPriorityQueue[E any](less func(E, E) bool) *PriorityQueue[E] {
	return &PriorityQueue[E]{sliceHeap[E]{less: less}}
}

func (pq *PriorityQueue[E]) Push(elem E) {
	heap.Push(&pq.s, elem)
}

func (pq *PriorityQueue[E]) Pop() E {
	return heap.Pop(&pq.s).(E)
}

func (pq *PriorityQueue[E]) Peek() E {
	return pq.s.s[0]
}

func (pq *PriorityQueue[E]) Len() int {
	return len(pq.s.s)
}

func (pq *PriorityQueue[E]) Slice() []E {
	return pq.s.s
}

func (pq *PriorityQueue[E]) SetIndex(f func(E, int)) {
	pq.s.setIndex = f
}
func (pq *PriorityQueue[E]) Fix(i int) {
	heap.Fix(&pq.s, i)
}

func (pq *PriorityQueue[E]) Remove(i int) E {
	return heap.Remove(&pq.s, i).(E)
}

type sliceHeap[E any] struct {
	s        []E
	less     func(E, E) bool
	setIndex func(E, int)
}

func (s *sliceHeap[E]) Len() int { return len(s.s) }

func (s *sliceHeap[E]) Swap(i, j int) {
	s.s[i], s.s[j] = s.s[j], s.s[i]
	if s.setIndex != nil {
		s.setIndex(s.s[i], i)
		s.setIndex(s.s[j], j)
	}
}

func (s *sliceHeap[E]) Less(i, j int) bool {
	return s.less(s.s[i], s.s[j])
}

func (s *sliceHeap[E]) Push(x interface{}) {
	s.s = append(s.s, x.(E))
	if s.setIndex != nil {
		s.setIndex(s.s[len(s.s)-1], len(s.s)-1)
	}
}

func (s *sliceHeap[E]) Pop() interface{} {
	e := s.s[len(s.s)-1]
	if s.setIndex != nil {
		s.setIndex(e, -1)
	}
	s.s = s.s[:len(s.s)-1]
	return e
}
