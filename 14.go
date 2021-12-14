package main

import (
	"strings"

	"golang.org/x/exp/maps"
)

func init() {
	addSolutions(14, problem14)
}

func problem14(ctx *problemContext) {
	var p polymer
	scanner := ctx.scanner()
	for scanner.scan() {
		p.parse(scanner.text())
	}
	ctx.reportLoad()

	for i := 0; i < 10; i++ {
		p.apply()
	}
	min, max := p.freq()
	ctx.reportPart1(max - min)

	for i := 0; i < 30; i++ {
		p.apply()
	}
	min, max = p.freq()
	ctx.reportPart2(max - min)
}

type polymer struct {
	pairs      map[[2]byte]int64
	start, end byte
	rules      map[[2]byte]byte
}

func (p *polymer) parse(line string) {
	if line == "" {
		return
	}
	if from, to, ok := strings.Cut(line, " -> "); ok {
		if p.rules == nil {
			p.rules = make(map[[2]byte]byte)
		}
		p.rules[[2]byte{from[0], from[1]}] = to[0]
		return
	}
	p.start = line[0]
	p.end = line[len(line)-1]
	p.pairs = make(map[[2]byte]int64)
	for i := 0; i < len(line)-1; i++ {
		p.pairs[[2]byte{line[i], line[i+1]}]++
	}
}

func (p *polymer) apply() {
	next := make(map[[2]byte]int64)
	for k, n := range p.pairs {
		r, ok := p.rules[k]
		if !ok {
			next[k] += n
			continue
		}
		next[[2]byte{k[0], r}] += n
		next[[2]byte{r, k[1]}] += n
	}
	p.pairs = next
}

func (p *polymer) freq() (min, max int64) {
	freqs := make(map[byte]int64)
	for k, n := range p.pairs {
		for _, c := range k {
			freqs[c] += n
		}
	}
	freqs[p.start]++
	freqs[p.end]++
	// Now they're all double-counted.
	vals := maps.Values(freqs)
	min = SliceMin(vals) / 2
	max = SliceMax(vals) / 2
	return min, max
}
