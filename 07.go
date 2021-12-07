package main

import (
	"io"
	"log"
	"strings"
)

func init() {
	addSolutions(7, problem7)
}

func problem7(ctx *problemContext) {
	b, err := io.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	input := strings.TrimSpace(string(b))
	var crabs []int64
	for _, v := range strings.Split(input, ",") {
		crabs = append(crabs, parseInt(v, 10, 64))
	}
	ctx.reportLoad()

	var min, max int64
	for i, crab := range crabs {
		if i == 0 {
			min, max = crab, crab
		} else {
			if crab < min {
				min = crab
			}
			if crab > max {
				max = crab
			}
		}
	}
	var best int64
	for c := min; c <= max; c++ {
		var sum int64
		for _, crab := range crabs {
			sum += abs(c - crab)
		}
		if c == min || sum < best {
			best = sum
		}
	}
	ctx.reportPart1(best)

	best = 0
	for c := min; c <= max; c++ {
		var sum int64
		for _, crab := range crabs {
			d := abs(c - crab)
			sum += (d + 1) * d / 2
		}
		if c == min || sum < best {
			best = sum
		}
	}
	ctx.reportPart1(best)
}
