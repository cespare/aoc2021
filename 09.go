package main

func init() {
	addSolutions(9, problem9)
}

func problem9(ctx *problemContext) {
	scanner := ctx.scanner()
	var m heightMap
	for scanner.scan() {
		b := []uint8(scanner.text())
		for i, c := range b {
			b[i] = c - '0'
		}
		m = append(m, b)
	}

	ctx.reportLoad()

	var lows []vec2
	var sum int64
	for x := 0; x < m.width(); x++ {
		for y := 0; y < m.height(); y++ {
			v := vec2{int64(x), int64(y)}
			if m.isLow(v) {
				lows = append(lows, v)
				sum += int64(m.at(v) + 1)
			}
		}
	}
	ctx.reportPart1(sum)

	basins := make([]int64, len(lows))
	for i, low := range lows {
		basins[i] = m.basinSize(low)
	}
	SortFunc(basins, func(a, b int64) bool { return a > b })
	ctx.reportPart2(basins[0] * basins[1] * basins[2])
}

type heightMap [][]uint8

func (m heightMap) at(v vec2) uint8 {
	return m[v.y][v.x]
}

func (m heightMap) width() int  { return len(m[0]) }
func (m heightMap) height() int { return len(m) }

func (m heightMap) neighbors(v vec2) []vec2 {
	var n []vec2
	if v.x > 0 {
		n = append(n, vec2{v.x - 1, v.y})
	}
	if v.x < int64(m.width()-1) {
		n = append(n, vec2{v.x + 1, v.y})
	}
	if v.y > 0 {
		n = append(n, vec2{v.x, v.y - 1})
	}
	if v.y < int64(m.height()-1) {
		n = append(n, vec2{v.x, v.y + 1})
	}
	return n
}

func (m heightMap) isLow(v vec2) bool {
	t := m.at(v)
	for _, nv := range m.neighbors(v) {
		if m.at(nv) <= t {
			return false
		}
	}
	return true
}

func (m heightMap) basinSize(low vec2) int64 {
	frontier := SetOf(low)
	basin := SetOf(low)
	for {
		var next Set[vec2]
		frontier.Do(func(v vec2) bool {
			for _, nv := range m.neighbors(v) {
				if basin.Contains(nv) {
					continue
				}
				nt := m.at(nv)
				if nt == 9 || nt < m.at(v) {
					continue
				}
				next.Add(nv)
				basin.Add(nv)
			}
			return true
		})
		if next.Len() == 0 {
			return int64(basin.Len())
		}
		frontier = &next
	}
}
