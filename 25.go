package main

func init() {
	addSolutions(25, problem25)
}

func problem25(ctx *problemContext) {
	var c seaCukes
	scanner := ctx.scanner()
	for scanner.scan() {
		c.g = append(c.g, []byte(scanner.text()))
	}
	c.init()
	ctx.reportLoad()

	for step := 1; ; step++ {
		if !c.step() {
			ctx.reportPart1(step)
			break
		}
	}
}

type seaCukes struct {
	g  [][]byte
	g1 [][]byte
}

func (c *seaCukes) init() {
	w, h := c.dims()
	c.g1 = make([][]byte, h)
	for i := range c.g1 {
		c.g1[i] = make([]byte, w)
	}
}

func (c *seaCukes) at(v vec2) byte {
	return c.g[v.y][v.x]
}

func (c *seaCukes) set1(v vec2, b byte) {
	c.g1[v.y][v.x] = b
}

func (c *seaCukes) step() bool {
	moved := c.moveDir('>', vec2{1, 0})
	moved += c.moveDir('v', vec2{0, 1})
	return moved > 0
}

func (c *seaCukes) moveDir(sym byte, d vec2) int64 {
	for y, row := range c.g {
		copy(c.g1[y], row)
	}
	w, h := c.dims()
	var moved int64
	for x := int64(0); x < w; x++ {
		for y := int64(0); y < h; y++ {
			v := vec2{x, y}
			if c.at(v) != sym {
				continue
			}
			targ := v.add(d)
			if targ.x >= w {
				targ.x -= w
			}
			if targ.y >= h {
				targ.y -= h
			}
			if c.at(targ) == '.' {
				c.set1(v, '.')
				c.set1(targ, sym)
				moved++
			}
		}
	}
	c.g, c.g1 = c.g1, c.g
	return moved
}

func (c *seaCukes) dims() (w, h int64) {
	return int64(len(c.g[0])), int64(len(c.g))
}
