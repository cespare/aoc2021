package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(13, problem13)
}

func problem13(ctx *problemContext) {
	scanner := ctx.scanner()
	var fd foldDots
	for scanner.scan() {
		fd.addLine(scanner.text())
	}
	ctx.reportLoad()

	fd.fold(fd.folds[0])
	ctx.reportPart1(fd.dots.Len())

	for i := 1; i < len(fd.folds); i++ {
		fd.fold(fd.folds[i])
	}
	fd.print()
	ctx.reportPart2("^^^^^^^^^^^^^")
}

type foldDots struct {
	dots  Set[vec2]
	folds []vec2 // either x or y are 0
}

func (fd *foldDots) addLine(s string) {
	if s == "" {
		return
	}
	if s, ok := trimPrefix(s, "fold along "); ok {
		var fold vec2
		if s, ok := trimPrefix(s, "x="); ok {
			fold = vec2{parseInt(s, 10, 64), 0}
		} else {
			s = strings.TrimPrefix(s, "y=")
			fold = vec2{0, parseInt(s, 10, 64)}
		}
		fd.folds = append(fd.folds, fold)
		return
	}
	xs, ys, ok := strings.Cut(s, ",")
	if !ok {
		panic("bad input")
	}
	fd.dots.Add(vec2{parseInt(xs, 10, 64), parseInt(ys, 10, 64)})
}

func (fd *foldDots) fold(f vec2) {
	var newDots []vec2
	fd.dots.Do(func(v vec2) bool {
		var vv vec2
		if f.y == 0 {
			if v.x < f.x {
				return true
			}
			vv = vec2{f.x - (v.x - f.x), v.y}
		} else {
			if v.y < f.y {
				return true
			}
			vv = vec2{v.x, f.y - (v.y - f.y)}
		}
		fd.dots.Remove(v)
		newDots = append(newDots, vv)
		return true
	})
	for _, dot := range newDots {
		fd.dots.Add(dot)
	}
}

func (fd *foldDots) print() {
	var extent vec2
	fd.dots.Do(func (v vec2) bool {
		if v.x > extent.x {
			extent.x = v.x
		}
		if v.y > extent.y {
			extent.y = v.y
		}
		return true
	})
	for y := int64(0); y <= extent.y; y++ {
		for x := int64(0); x < extent.x; x++ {
			v := vec2{x, y}
			if fd.dots.Contains(v) {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
