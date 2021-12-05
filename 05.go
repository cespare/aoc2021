package main

import (
	"strings"
)

func init() {
	addSolutions(5, problem5)
}

func problem5(ctx *problemContext) {
	var vents []vent
	scanner := ctx.scanner()
	for scanner.scan() {
		parts := strings.Split(scanner.text(), " -> ")
		c0 := strings.Split(parts[0], ",")
		x0 := parseInt(c0[0], 10, 64)
		y0 := parseInt(c0[1], 10, 64)
		c1 := strings.Split(parts[1], ",")
		x1 := parseInt(c1[0], 10, 64)
		y1 := parseInt(c1[1], 10, 64)
		vents = append(vents, vent{vec2{x0, y0}, vec2{x1, y1}})
	}
	ctx.reportLoad()

	ventLines1 := make(map[vec2]int)
	for _, vent := range vents {
		v0, v1 := vent.v0, vent.v1
		switch {
		case v0.x == v1.x:
			if v0.y > v1.y {
				v0, v1 = v1, v0
			}
			for y := v0.y; y <= v1.y; y++ {
				ventLines1[vec2{v0.x, y}]++
			}
		case v0.y == v1.y:
			if v0.x > v1.x {
				v0, v1 = v1, v0
			}
			for x := v0.x; x <= v1.x; x++ {
				ventLines1[vec2{x, v0.y}]++
			}
		}
	}
	var overlaps int
	for _, n := range ventLines1 {
		if n > 1 {
			overlaps++
		}
	}
	ctx.reportPart1(overlaps)

	ventLines2 := make(map[vec2]int)
	for _, vent := range vents {
		v0, v1 := vent.v0, vent.v1
		switch {
		case v0.x == v1.x:
			if v0.y > v1.y {
				v0, v1 = v1, v0
			}
			for y := v0.y; y <= v1.y; y++ {
				ventLines2[vec2{v0.x, y}]++
			}
		case v0.y == v1.y:
			if v0.x > v1.x {
				v0, v1 = v1, v0
			}
			for x := v0.x; x <= v1.x; x++ {
				ventLines2[vec2{x, v0.y}]++
			}
		default:
			if v0.y > v1.y {
				v0, v1 = v1, v0
			}
			if v0.x < v1.x {
				for x, y := v0.x, v0.y; x <= v1.x && y <= v1.y; x, y = x+1, y+1 {
					ventLines2[vec2{x, y}]++
				}
			} else {
				for x, y := v0.x, v0.y; x >= v1.x && y <= v1.y; x, y = x-1, y+1 {
					ventLines2[vec2{x, y}]++
				}
			}
		}
	}
	overlaps = 0
	for _, n := range ventLines2 {
		if n > 1 {
			overlaps++
		}
	}
	ctx.reportPart2(overlaps)
}

type vent struct {
	v0 vec2
	v1 vec2
}
