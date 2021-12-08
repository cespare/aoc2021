package main

import (
	"io"
	"log"
	"strings"
)

func init() {
	addSolutions(6, problem6)
}

func problem6(ctx *problemContext) {
	b, err := io.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	input := strings.TrimSpace(string(b))
	var ages []int64
	for _, v := range strings.Split(input, ",") {
		ages = append(ages, parseInt(v, 10, 64))
	}
	ctx.reportLoad()

	lf := newLanternfish(ages)
	for i := 0; i < 80; i++ {
		lf.step()
	}
	ctx.reportPart1(lf.count())

	lf = newLanternfish(ages)
	for i := 0; i < 256; i++ {
		lf.step()
	}
	ctx.reportPart2(lf.count())
}

type lanternfish map[int]int64

func newLanternfish(ages []int64) *lanternfish {
	lf := make(lanternfish)
	for _, age := range ages {
		lf[int(age)]++
	}
	return &lf
}

func (lf *lanternfish) step() {
	next := make(lanternfish)
	for age, n := range *lf {
		if age == 0 {
			next[8] += n
			next[6] += n
		} else {
			next[age-1] += n
		}
	}
	*lf = next
}

func (lf *lanternfish) count() int64 {
	var total int64
	for _, n := range *lf {
		total += n
	}
	return total
}
