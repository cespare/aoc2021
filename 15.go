package main

import (
	"math"
	"strings"
)

func init() {
	addSolutions(15, problem15)
}

func problem15(ctx *problemContext) {
	var m riskMap
	scanner := ctx.scanner()
	for scanner.scan() {
		m.parseLine(scanner.text())
	}
	ctx.reportLoad()

	m2 := m.multiply()

	ctx.reportPart1(m.bestPath())
	ctx.reportPart2(m2.bestPath())
}

type riskMap struct {
	g [][]int64
}

func (m *riskMap) multiply() *riskMap {
	n := len(m.g) * 5
	g := make([][]int64, n)
	for i := range g {
		g[i] = make([]int64, n)
	}
	for xi := 0; xi < 5; xi++ {
		for yi := 0; yi < 5; yi++ {
			for x := 0; x < len(m.g); x++ {
				for y := 0; y < len(m.g); y++ {
					r := m.g[y][x] + int64(xi+yi)
					r = (r-1)%9 + 1
					xx := xi*len(m.g) + x
					yy := yi*len(m.g) + y
					g[yy][xx] = r
				}
			}
		}
	}
	return &riskMap{g: g}
}

func (m *riskMap) parseLine(s string) {
	s = strings.TrimSpace(s)
	var row []int64
	for i := 0; i < len(s); i++ {
		n := int64(s[i] - '0')
		row = append(row, n)
	}
	m.g = append(m.g, row)
}

func (m *riskMap) at(v vec2) int64 {
	return m.g[v.y][v.x]
}

func (m *riskMap) end() vec2 {
	return vec2{int64(len(m.g[0]) - 1), int64(len(m.g) - 1)}
}

func (m *riskMap) neighbors(v vec2) []vec2 {
	var neigh []vec2
	end := m.end()
	if v.x > 0 {
		neigh = append(neigh, vec2{v.x - 1, v.y})
	}
	if v.y > 0 {
		neigh = append(neigh, vec2{v.x, v.y - 1})
	}
	if v.x < end.x {
		neigh = append(neigh, vec2{v.x + 1, v.y})
	}
	if v.y < end.y {
		neigh = append(neigh, vec2{v.x, v.y + 1})
	}
	return neigh
}

type vec2Risk struct {
	v vec2
	r int64
}

func (m *riskMap) bestPath() int64 {
	end := m.end()
	// Upper is the best "worst-case estimate" we've seen: given a
	// particular risk value of a cell, what would the risk be to extend
	// that to (0,0) if it were all 9s along the way? We can use this to
	// directly eliminate non-promising paths.
	upper := int64(math.MaxInt64)
	costs := map[vec2]int64{end: 0}
	less := func(vr0, vr1 vec2Risk) bool { return vr0.r < vr1.r }
	frontier := NewPriorityQueue(less)
	frontier.Push(vec2Risk{end, 0})
	for frontier.Len() > 0 {
		nr := frontier.Pop()
		// If the best-case estimate for this cell (extending to (0, 0)
		// along a path of all 1s) is no better than the current upper
		// bound, eliminate it.
		if nr.r+nr.v.x+nr.v.y >= upper {
			continue
		}
		est := nr.r + (nr.v.x+nr.v.y)*9
		if est < upper {
			upper = est
		}
		for _, n := range m.neighbors(nr.v) {
			r := nr.r + m.at(nr.v)
			prev, ok := costs[n]
			if ok && r >= prev {
				continue
			}
			costs[n] = r
			frontier.Push(vec2Risk{n, r})
		}
	}
	return costs[vec2{0, 0}]
}
