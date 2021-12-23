package main

import (
	"container/heap"
	"strings"
)

func init() {
	addSolutions(23, problem23)
}

func problem23(ctx *problemContext) {
	initial := burrowState{roomLen: 2}
	scanner := ctx.scanner()
	var i int
	for scanner.scan() {
		line := scanner.text()
		switch i {
		case 2:
			initial.rooms[0][1] = line[3]
			initial.rooms[1][1] = line[5]
			initial.rooms[2][1] = line[7]
			initial.rooms[3][1] = line[9]
		case 3:
			initial.rooms[0][0] = line[3]
			initial.rooms[1][0] = line[5]
			initial.rooms[2][0] = line[7]
			initial.rooms[3][0] = line[9]
		}
		i++
	}
	ctx.reportLoad()

	ctx.reportPart1(bestBurrowOrg(initial))

	initial.roomLen = 4
	for i := range initial.rooms {
		initial.rooms[i][3] = initial.rooms[i][1]
	}
	initial.rooms[0][1] = 'D'
	initial.rooms[0][2] = 'D'
	initial.rooms[1][1] = 'B'
	initial.rooms[1][2] = 'C'
	initial.rooms[2][1] = 'A'
	initial.rooms[2][2] = 'B'
	initial.rooms[3][1] = 'C'
	initial.rooms[3][2] = 'A'

	ctx.reportPart2(bestBurrowOrg(initial))
}

const hallLen = 11

type burrowState struct {
	rooms   [4][4]byte
	roomLen int
	hall    [hallLen]byte
}

func (st burrowState) String() string {
	var b strings.Builder
	for _, h := range st.hall {
		if h == 0 {
			b.WriteByte('.')
		} else {
			b.WriteByte(h)
		}
	}
	for ri := st.roomLen - 1; ri >= 0; ri-- {
		b.WriteString("\n  ")
		for _, room := range st.rooms {
			if room[ri] == 0 {
				b.WriteByte('.')
			} else {
				b.WriteByte(room[ri])
			}
			b.WriteByte(' ')
		}
	}
	return b.String()
}

type burrowStateCost struct {
	st   burrowState
	cost int64
}

type burrowQueue struct {
	indexes map[burrowState]int
	q       []burrowStateCost
}

func (q *burrowQueue) Len() int { return len(q.q) }
func (q *burrowQueue) Swap(i, j int) {
	q.indexes[q.q[i].st] = j
	q.indexes[q.q[j].st] = i
	q.q[i], q.q[j] = q.q[j], q.q[i]
}

func (q *burrowQueue) Less(i, j int) bool {
	return q.q[i].cost < q.q[j].cost
}

func (q *burrowQueue) Push(x interface{}) {
	sc := x.(burrowStateCost)
	q.indexes[sc.st] = len(q.q)
	q.q = append(q.q, sc)
}

func (q *burrowQueue) Pop() interface{} {
	sc := q.q[len(q.q)-1]
	delete(q.indexes, sc.st)
	q.q = q.q[:len(q.q)-1]
	return sc
}

func (q *burrowQueue) upsert(st burrowState, cost int64) {
	i, ok := q.indexes[st]
	if !ok {
		heap.Push(q, burrowStateCost{st, cost})
		return
	}
	q.q[i].cost = cost
	heap.Fix(q, i)
}

func burrowOK(st burrowState) bool {
	for _, h := range st.hall {
		if h != 0 {
			return false
		}
	}
	for i, room := range st.rooms {
		for k := 0; k < st.roomLen; k++ {
			if !amphMatches(room[k], i) {
				return false
			}
		}
	}
	return true
}

func bestBurrowOrg(initial burrowState) int64 {
	costs := map[burrowState]int64{initial: 0}
	q := &burrowQueue{indexes: make(map[burrowState]int)}
	heap.Push(q, burrowStateCost{initial, 0})
	pushStateCost := func(st burrowState, cost int64) {
		if prev, ok := costs[st]; ok && prev <= cost {
			return
		}
		costs[st] = cost
		q.upsert(st, cost)
	}
	for q.Len() > 0 {
		sc := heap.Pop(q).(burrowStateCost)
		if burrowOK(sc.st) {
			return sc.cost
		}
		for i, room := range sc.st.rooms {
			ihall := 2*i + 2
			roomOK := true
			var ktop int // first unoccupied room index
			for ; ktop < sc.st.roomLen; ktop++ {
				if room[ktop] == 0 {
					break
				}
				if !amphMatches(room[ktop], i) {
					roomOK = false
				}
			}
			if roomOK {
				if ktop == sc.st.roomLen {
					// Room is all full and ready to go.
					// Nothing left to do.
					continue
				}
				// See which amphs we can move in.
				steps := sc.st.roomLen - ktop
				for j := ihall - 1; j >= 0; j-- {
					amph := sc.st.hall[j]
					if amph == 0 {
						continue
					}
					if !amphMatches(amph, i) {
						break
					}
					st := sc.st
					st.rooms[i][ktop] = amph
					st.hall[j] = 0
					cost := sc.cost + int64(steps+ihall-j)*amphCost(amph)
					pushStateCost(st, cost)
					break
				}
				for j := ihall + 1; j < hallLen; j++ {
					amph := sc.st.hall[j]
					if amph == 0 {
						continue
					}
					if !amphMatches(amph, i) {
						break
					}
					st := sc.st
					st.rooms[i][ktop] = amph
					st.hall[j] = 0
					cost := sc.cost + int64(steps+j-ihall)*amphCost(amph)
					pushStateCost(st, cost)
					break
				}
				continue
			}

			// Room still needs to be emptied out.
			steps := sc.st.roomLen - ktop + 1
			for j := ihall - 1; j >= 0; j-- {
				if sc.st.hall[j] != 0 {
					break // blocked
				}
				switch j {
				case 2, 4, 6, 8:
					continue
				}
				st := sc.st
				amph := room[ktop-1]
				st.rooms[i][ktop-1] = 0
				st.hall[j] = amph
				cost := sc.cost + int64(steps+ihall-j)*amphCost(amph)
				pushStateCost(st, cost)
			}
			for j := ihall + 1; j < hallLen; j++ {
				if sc.st.hall[j] != 0 {
					break // blocked
				}
				switch j {
				case 2, 4, 6, 8:
					continue
				}
				st := sc.st
				amph := room[ktop-1]
				st.rooms[i][ktop-1] = 0
				st.hall[j] = amph
				cost := sc.cost + int64(steps+j-ihall)*amphCost(amph)
				pushStateCost(st, cost)
			}
		}
	}
	panic("no solution")
}

func amphCost(a byte) int64 {
	switch a {
	case 'A':
		return 1
	case 'B':
		return 10
	case 'C':
		return 100
	case 'D':
		return 1000
	default:
		panic("bad")
	}
}

func amphMatches(a byte, iroom int) bool {
	if a < 'A' || a > 'D' {
		panic("bad")
	}
	return int(a-'A') == iroom
}
