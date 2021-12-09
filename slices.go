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
