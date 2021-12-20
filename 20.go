package main

func init() {
	addSolutions(20, problem20)
}

func problem20(ctx *problemContext) {
	var alg []byte
	lit := new(Set[vec2])
	var y int64
	scanner := ctx.scanner()
	for scanner.scan() {
		if len(alg) == 0 {
			alg = []byte(scanner.text())
			if len(scanner.text()) != 512 {
				panic("bad alg")
			}
			continue
		}
		if len(scanner.text()) == 0 {
			continue
		}
		for x, r := range scanner.text() {
			if r == '#' {
				lit.Add(vec2{int64(x), y})
			}
		}
		y++
	}
	ctx.reportLoad()

	img := newImage(lit, alg)
	for i := 0; i < 2; i++ {
		img = img.enhance()
	}
	ctx.reportPart1(img.lit.Len())

	for i := 0; i < 48; i++ {
		img = img.enhance()
	}
	ctx.reportPart2(img.lit.Len())
}

type image struct {
	lit        *Set[vec2]
	alg        []byte
	strobe     bool
	minx, miny int64
	maxx, maxy int64
	spaceLit   bool
}

func newImage(lit *Set[vec2], alg []byte) *image {
	im := &image{
		lit:    lit,
		alg:    alg,
		strobe: alg[0] == '#',
	}
	im.updateBounds()
	return im
}

func (im *image) updateBounds() {
	im.lit.Do(func(v vec2) bool {
		if v.x < im.minx {
			im.minx = v.x
		}
		if v.x > im.maxx {
			im.maxx = v.x
		}
		if v.y < im.miny {
			im.miny = v.y
		}
		if v.y > im.maxy {
			im.maxy = v.y
		}
		return true
	})
}

func (im *image) enhance() *image {
	var out Set[vec2]
	for x := im.minx - 1; x <= im.maxx+1; x++ {
		for y := im.miny - 1; y <= im.maxy+1; y++ {
			v := vec2{x, y}
			idx := im.index(v)
			if im.alg[idx] == '#' {
				out.Add(v)
			}
		}
	}
	im1 := &image{
		lit:      &out,
		alg:      im.alg,
		strobe:   im.strobe,
		spaceLit: im.strobe && !im.spaceLit,
	}
	im1.updateBounds()
	return im1
}

func (im *image) index(v vec2) int {
	var idx int
	for dy := int64(-1); dy <= 1; dy++ {
		for dx := int64(-1); dx <= 1; dx++ {
			idx <<= 1
			v := vec2{v.x + dx, v.y + dy}
			if v.x < im.minx || v.x > im.maxx || v.y < im.miny || v.y > im.maxy {
				if im.spaceLit {
					idx |= 1
				}
			}
			if im.lit.Contains(v) {
				idx |= 1
			}
		}
	}
	return idx
}
