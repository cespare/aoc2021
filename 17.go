package main

import (
	"io"
	"log"
	"math"
	"strings"
)

func init() {
	addSolutions(17, problem17)
}

func problem17(ctx *problemContext) {
	b, err := io.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.TrimSpace(string(b))
	s = strings.TrimPrefix(s, "target area: ")
	xs, ys, _ := strings.Cut(s, ", ")
	xmin, xmax := parseRange(strings.TrimPrefix(xs, "x="))
	ymin, ymax := parseRange(strings.TrimPrefix(ys, "y="))
	ctx.reportLoad()

	// Figure out some bounds.
	// If we don't have enough vx, we won't make it to vmin before vx drops
	// to 0. We can work this out by solving a quadratic.
	vx0 := int64((math.Sqrt(8*float64(xmin)+1) - 1) / 2)
	// If vx is large enough, we'll skip past vmax on the very first step.
	vx1 := xmax + 1
	// If vy is too small (negative with large magnitude), same thing wrt
	// ymin.
	vy0 := ymin - 1
	// If vy is too large, we'll have the same problem when we return to y=0
	// at the same velocity (but negative).
	vy1 := -vy0

	var best int64
	var numInit int64
	for vy := vy0; vy <= vy1; vy++ {
		for vx := vx0; vx <= vx1; vx++ {
			v := vec2{vx, vy}
			if maxHeight, ok := launch(v, xmin, xmax, ymin, ymax); ok {
				numInit++
				if maxHeight > best {
					best = maxHeight
				}

			}
		}
	}
	ctx.reportPart1(best)
	ctx.reportPart1(numInit)
}

func parseRange(s string) (min, max int64) {
	m, n, _ := strings.Cut(s, "..")
	return parseInt(m, 10, 64), parseInt(n, 10, 64)
}

func launch(v vec2, xmin, xmax, ymin, ymax int64) (maxHeight int64, ok bool) {
	var p vec2
	for {
		// Check
		if p.y > maxHeight {
			maxHeight = p.y
		}
		if p.x > xmax {
			return 0, false
		} else if p.x >= xmin {
			if p.y < ymin {
				return 0, false
			}
			if p.y <= ymax {
				return maxHeight, true
			}
		} else { // p.x < xmin
			if p.y < ymin {
				return 0, false
			}
			if v.x == 0 {
				return 0, false
			}
		}

		// Update
		p.x += v.x
		p.y += v.y
		if v.x > 0 {
			v.x--
		}
		if v.x < 0 {
			v.x++
		}
		v.y--
	}
}
