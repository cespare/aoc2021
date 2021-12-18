package main

import (
	"constraints"
	"sort"
)

// Stuff not yet in golang.org/x/exp/slices.
// See https://github.com/golang/go/issues/47619#issuecomment-915428658

func SliceSort[E constraints.Ordered](x []E) {
	sort.Slice(x, func(i, j int) bool { return x[i] < x[j] })
}

func SortFunc[E any](x []E, less func(a, b E) bool) {
	sort.Slice(x, func(i, j int) bool { return less(x[i], x[j]) })
}

// Extra slice stuff.

func SliceMin[E constraints.Ordered](x []E) E {
	if len(x) == 0 {
		panic("SliceMin of 0 elements")
	}
	min := x[0]
	for i := 1; i < len(x); i++ {
		if x[i] < min {
			min = x[i]
		}
	}
	return min
}

func SliceMax[E constraints.Ordered](x []E) E {
	if len(x) == 0 {
		panic("SliceMax of 0 elements")
	}
	max := x[0]
	for i := 1; i < len(x); i++ {
		if x[i] > max {
			max = x[i]
		}
	}
	return max
}

func SliceReduce[S ~[]E, E, R any](s S, initial R, fn func(R, E) R) R {
	r := initial
	for _, e := range s {
		r = fn(r, e)
	}
	return r
}
