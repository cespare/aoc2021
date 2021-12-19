package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(19, problem19)
}

func problem19(ctx *problemContext) {
	var beacons [][]vec3
	scn := ctx.scanner()
	var curBeacon []vec3
	needBeacon := true
	for scn.scan() {
		if needBeacon {
			want := fmt.Sprintf("--- scanner %d ---", len(beacons))
			if scn.text() != want {
				panic(scn.text())
			}
			needBeacon = false
			continue
		}
		if scn.text() == "" {
			beacons = append(beacons, curBeacon)
			curBeacon = nil
			needBeacon = true
			continue
		}
		curBeacon = append(curBeacon, parseBeacon(scn.text()))
	}
	if len(curBeacon) > 0 {
		beacons = append(beacons, curBeacon)
	}
	ctx.reportLoad()

	allBeacons := SetOf(beacons[0]...)
	firstScanner := &beaconScanner{
		beacons: beacons[0],
	}
	solved := []*beaconScanner{firstScanner}
	unsolved := beacons[1:]

	for len(unsolved) > 0 {
	unsolvedLoop:
		for i := 0; i < len(unsolved); i++ {
			bs := unsolved[i]
			for _, o := range scannerOrientations {
				var reoriented []vec3
				for _, v := range bs {
					re := reorientBeacon(v, o)
					reoriented = append(reoriented, re)
				}
				for _, s0 := range solved {
					off, ok := s0.matches(reoriented)
					if !ok {
						continue
					}
					s := &beaconScanner{
						orn:     o,
						offset:  off,
						beacons: moveBeacons(reoriented, off),
					}
					allBeacons.Add(s.beacons...)
					unsolved[i] = unsolved[len(unsolved)-1]
					unsolved = unsolved[:len(unsolved)-1]
					i--
					solved = append(solved, s)
					continue unsolvedLoop
				}
			}
		}
	}
	ctx.reportPart1(allBeacons.Len())

	var best int64
	for _, s0 := range solved {
		for _, s1 := range solved {
			best = max(best, s0.offset.hamming(s1.offset))
		}
	}
	ctx.reportPart2(best)
}

func moveBeacons(bs []vec3, off vec3) []vec3 {
	return SliceMap(bs, func(v vec3) vec3 { return v.add(off) })
}

type beaconScanner struct {
	orn     scannerOrientation // relative to scanner[0]
	offset  vec3               // after reorienting
	beacons []vec3             // after reorientation and offset
}

func parseBeacon(line string) vec3 {
	parts := strings.Split(line, ",")
	return vec3{
		x: parseInt(parts[0], 10, 64),
		y: parseInt(parts[1], 10, 64),
		z: parseInt(parts[2], 10, 64),
	}
}

func (s *beaconScanner) matches(beacons []vec3) (off vec3, ok bool) {
	offsetFreqs := make(map[vec3]int)
	for _, v0 := range s.beacons {
		for _, v1 := range beacons {
			offsetFreqs[v0.sub(v1)]++
		}
	}
	found := false
	for o, freq := range offsetFreqs {
		if freq < 12 {
			continue
		}
		if found {
			panic("multiple matches (TODO: check further)")
		}
		found = true
		off = o
	}
	return off, found
}

type scannerOrientation struct {
	front int // in {1, 2, 3, -1, -2, -3}; indicate {+x, +y, +z, -x, -y, -z}
	rot   int // in {0, 1, 2, 3}; indicate 0, 90, 180, 270 deg around forward axis
}

var scannerOrientations []scannerOrientation

func init() {
	for _, front := range []int{1, 2, 3, -1, -2, -3} {
		for rot := 0; rot <= 3; rot++ {
			o := scannerOrientation{front: front, rot: rot}
			scannerOrientations = append(scannerOrientations, o)
		}
	}
}

func reorientBeacon(v vec3, o scannerOrientation) vec3 {
	switch o.front {
	case 1:
	case 2:
		v.x, v.y = v.y, -v.x
	case 3:
		v.x, v.z = v.z, -v.x
	case -1:
		v.x, v.y = -v.x, -v.y
	case -2:
		v.x, v.y = -v.y, v.x
	case -3:
		v.x, v.z = -v.z, v.x
	default:
		panic("bad")
	}
	switch o.rot {
	case 0:
	case 1:
		v.y, v.z = -v.z, v.y
	case 2:
		v.y, v.z = -v.y, -v.z
	case 3:
		v.y, v.z = v.z, -v.y
	default:
		panic("bad")
	}
	return v
}
