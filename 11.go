package main

func init() {
	addSolutions(11, problem11)
}

func problem11(ctx *problemContext) {
	var o octopuses
	scanner := ctx.scanner()
	y := 0
	for scanner.scan() {
		s := scanner.text()
		for x := 0; x < 10; x++ {
			o.grid[y][x] = s[x] - '0'
		}
		y++
	}
	ctx.reportLoad()

	for i := 0; ; i++ {
		n := o.step()
		if i == 99 {
			ctx.reportPart1(o.flashes)
		}
		if n == 100 {
			ctx.reportPart2(i+1)
			break
		}
	}
}

type octopuses struct {
	grid    [10][10]uint8
	flashes int64
}

func (o *octopuses) at(v vec2) uint8 {
	return o.grid[v.y][v.x]
}

func (o *octopuses) inc(v vec2) (flashed bool) {
	x := o.at(v)
	if x < 9 {
		x++
	} else {
		x = 255
		o.flashes++
		flashed = true
	}
	o.grid[v.y][v.x] = x
	return flashed
}

func (o *octopuses) step() int {
	var flashed Set[vec2]
	var frontier Set[vec2]
	for y := int64(0); y < 10; y++ {
		for x := int64(0); x < 10; x++ {
			v := vec2{x, y}
			if o.inc(v) {
				flashed.Add(v)
				frontier.Add(v)
			}
		}
	}

	for {
		var next Set[vec2]
		for frontier.Len() > 0 {
			v := frontier.Pop()
			for dy := int64(-1); dy <= 1; dy++ {
				for dx := int64(-1); dx <= 1; dx++ {
					if dy == 0 && dx == 0 {
						continue
					}
					y := v.y+dy
					x := v.x+dx
					if y < 0 || y >= 10 {
						continue
					}
					if x < 0 || x >= 10 {
						continue
					}
					vv := vec2{x, y}
					if flashed.Contains(vv) {
						continue
					}
					if o.inc(vv) {
						flashed.Add(vv)
						next.Add(vv)
					}
				}
			}
		}
		if next.Len() == 0 {
			break
		}
		frontier = next
	}
	for y := int64(0); y < 10; y++ {
		for x := int64(0); x < 10; x++ {
			if o.grid[y][x] == 255 {
				o.grid[y][x] = 0
			}
		}
	}
	return flashed.Len()
}
