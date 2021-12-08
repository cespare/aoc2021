package main

import (
	"strings"

	"golang.org/x/exp/maps"
	"github.com/cespare/permute"
)

func init() {
	addSolutions(8, problem8)
}

func problem8(ctx *problemContext) {
	scanner := ctx.scanner()
	var entries []*ssdEntry
	for scanner.scan() {
		entries = append(entries, parseSSDEntry(scanner.text()))
	}
	ctx.reportLoad()

	var n int
	for _, entry := range entries {
		for _, v := range entry.out {
			switch len(v) {
			case 2, 3, 4, 7:
				n++
			}
		}
	}
	ctx.reportPart1(n)

	var sum int64
	for _, entry := range entries {
		sum += int64(entry.determineOutput())
	}
	ctx.reportPart2(sum)
}

type ssdEntry struct {
	in  [10]string
	out [4]string
}

func parseSSDEntry(s string) *ssdEntry {
	parts := strings.Split(s, " | ")
	in := strings.Fields(parts[0])
	out := strings.Fields(parts[1])
	for i, s := range in {
		in[i] = sortString(s)
	}
	for i, s := range out {
		out[i] = sortString(s)
	}
	return &ssdEntry{
		in:  *(*[10]string)(in),
		out: *(*[4]string)(out),
	}
}

func (e *ssdEntry) determineOutput() int {
	var candidates Set[int]
	for _, s := range e.in {
		m := ssdCombos[s]
		if candidates.Len() == 0 {
			candidates.Add(maps.Keys(m)...)
		} else {
			candidates.Filter(func(n int) bool {
				_, ok := m[n]
				return ok
			})
		}
	}
	perm := candidates.Pop()
	var result int
	digit := 1000
	for _, s := range e.out {
		i := ssdCombos[s][perm]
		result += i*digit
		digit /= 10
	}
	return result
}

// sorted string -> index of permutation -> digit meaning
var ssdCombos = make(map[string]map[int]int)

func init() {
	template := [10]string{
		"abcefg",
		"cf",
		"acdeg",
		"acdfg",
		"bcdf",
		"abdfg",
		"abdefg",
		"acf",
		"abcdefg",
		"abcdfg",
	}
	ord := []int{0, 1, 2, 3, 4, 5, 6}
	p := permute.Ints(ord)
	i := 0
	for p.Permute() {
		for d, t := range template {
			b := make([]byte, len(t))
			for j, c := range t {
				b[j] = byte(ord[c-'a']) + 'a'
			}
			SliceSort(b)
			s := string(b)
			m, ok := ssdCombos[s]
			if !ok {
				m = make(map[int]int)
				ssdCombos[s] = m
			}
			m[i] = d
		}
		i++
	}
}

func sortString(s string) string {
	b := []byte(s)
	SliceSort(b)
	return string(b)
}
