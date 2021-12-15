package main

type Queue[E any] struct {
	q []E
}

func QueueOf[E any](e ...E) *Queue[E] {
	return &Queue[E]{q: e}
}

func (q *Queue[E]) Len() int { return len(q.q) }
func (q *Queue[E]) Push(e E) { q.q = append(q.q, e) }

func (q *Queue[E]) Pop() (E, bool) {
	if len(q.q) == 0 {
		var e E
		return e, false
	}
	e := q.q[0]
	q.q = q.q[1:]
	return e, true
}
