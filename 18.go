package main

import (
	"strings"

	"golang.org/x/exp/slices"
)

func init() {
	addSolutions(18, problem18)
}

func problem18(ctx *problemContext) {
	nums := scanSlice(ctx, parseSFNum)
	ctx.reportLoad()

	sum := SliceReduce(nums, nil, addSFNums)
	ctx.reportPart1(sum.mag())

	var best int64
	for i, n0 := range nums {
		for j, n1 := range nums {
			if i != j {
				best = max(best, addSFNums(n0, n1).mag())
			}
		}
	}
	ctx.reportPart2(best)
}

type sfNum []any // list of tokens: rune ('[', ',', ']') or int64

func parseSFNum(s string) sfNum {
	var num sfNum
	for len(s) > 0 {
		i := strings.IndexFunc(s, func(r rune) bool { return r < '0' || r > '9' })
		if i == 0 {
			num = append(num, rune(s[i]))
			s = s[1:]
		} else {
			num = append(num, parseInt(s[:i], 10, 64))
			s = s[i:]
		}
	}
	return num
}

func collectSFNum(vs ...any) sfNum {
	return SliceReduce(vs, sfNum(nil), func(n sfNum, v any) sfNum {
		if n1, ok := v.(sfNum); ok {
			return append(n, n1...)
		}
		return append(n, v)
	})
}

func addSFNums(n0, n1 sfNum) sfNum {
	if n0 == nil {
		return n1
	}
	n := collectSFNum('[', n0, ',', n1, ']')

reduce:
	var ok bool
	if n, ok = n.explode(); ok {
		goto reduce
	}
	if n, ok = n.split(); ok {
		goto reduce
	}
	return n
}

func (n sfNum) explode() (sfNum, bool) {
	var depth int
	idx := slices.IndexFunc(n, func(t any) bool {
		switch t {
		case '[':
			if depth == 4 {
				return true
			}
			depth++
		case ']':
			depth--
		}
		return false
	})
	if idx < 0 {
		return n, false
	}

	lv, rv := n[idx+1].(int64), n[idx+3].(int64)
	left, right := n[:idx], n[idx+5:]
	for i := len(left) - 1; i >= 0; i-- {
		if v0, ok := left[i].(int64); ok {
			left[i] = v0 + lv
			break
		}
	}
	for i := 0; i < len(right); i++ {
		if v0, ok := right[i].(int64); ok {
			right[i] = v0 + rv
			break
		}
	}
	return collectSFNum(left, int64(0), right), true
}

func (n sfNum) split() (sfNum, bool) {
	for i, t := range n {
		if v, ok := t.(int64); ok && v >= 10 {
			lv, rv := v/2, (v+1)/2
			left, right := n[:i], n[i+1:]
			return collectSFNum(left, '[', lv, ',', rv, ']', right), true
		}
	}
	return n, false
}

func (n sfNum) mag() int64 {
	// The multiplier for each individual value is some combination of 2s
	// and 3s. Instead of reconstructing the nesting, figure out the
	// multiplier as we go and add each value directly.
	var m int64
	mul := int64(1)
	for _, t := range n {
		if v, ok := t.(int64); ok {
			m += mul * v
			continue
		}
		switch t {
		case '[':
			mul *= 3
		case ',':
			mul /= 3
			mul *= 2
		case ']':
			mul /= 2
		}
	}
	return m
}
