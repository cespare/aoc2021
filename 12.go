package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(12, problem12)
}

func problem12(ctx *problemContext) {
	scanner := ctx.scanner()
	var m caveMap
	for scanner.scan() {
		m.addEdge(scanner.text())
	}
	ctx.reportLoad()

	ctx.reportPart1(m.paths("start", make(map[string]int), false))
	ctx.reportPart2(m.paths("start", make(map[string]int), true))
}

type caveMap struct {
	adj map[string]*Set[string]
}

func (m *caveMap) print() {
	for k, v := range m.adj {
		fmt.Printf("%s -> %s\n", k, v)
	}
}

func (m *caveMap) addEdge(s string) {
	v0, v1, ok := strings.Cut(s, "-")
	if !ok {
		panic("bad input")
	}
	if m.adj == nil {
		m.adj = make(map[string]*Set[string])
	}
	s0, ok := m.adj[v0]
	if !ok {
		s0 = new(Set[string])
		m.adj[v0] = s0
	}
	s0.Add(v1)

	s1, ok := m.adj[v1]
	if !ok {
		s1 = new(Set[string])
		m.adj[v1] = s1
	}
	s1.Add(v0)
}

func (m *caveMap) paths(cur string, visited map[string]int, allow2 bool) int {
	if cur == "end" {
		return 1
	}
	if n := visited[cur]; n > 0 && cur[0] >= 'a' && cur[0] <= 'z' {
		if allow2 {
			allow2 = false
		} else {
			return 0
		}
	}
	visited[cur]++
	var total int
	m.adj[cur].Do(func(n string) bool {
		if n != "start" {
			total += m.paths(n, visited, allow2)
		}
		return true
	})
	visited[cur]--
	return total
}
