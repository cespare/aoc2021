package main

import (
	"strings"
)

func init() {
	addSolutions(22, problem22)
}

func problem22(ctx *problemContext) {
	var steps []rebootStep
	scanner := ctx.scanner()
	for scanner.scan() {
		line := scanner.text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		onOff, rest, _ := strings.Cut(scanner.text(), " ")
		min, max := parseCuboid(rest)
		step := rebootStep{
			on:  onOff == "on",
			min: min,
			max: max.add(vec3{1, 1, 1}), // half-open intervals are easier
		}
		steps = append(steps, step)
	}
	ctx.reportLoad()

	t := new(reactorTree)
	inPart1 := func(v vec3) bool {
		return v.x >= -50 && v.x <= 51 &&
			v.y >= -50 && v.y <= 51 &&
			v.z >= -50 && v.z <= 51
	}
	for _, step := range steps {
		if inPart1(step.min) && inPart1(step.max) {
			t.insert(step)
		}
	}
	ctx.reportPart1(t.count())

	t = new(reactorTree)
	for _, step := range steps {
		t.insert(step)
	}
	ctx.reportPart2(t.count())
}

type rebootStep struct {
	on  bool
	min vec3
	max vec3
}

func parseCuboid(s string) (min, max vec3) {
	parts := strings.Split(s, ",")
	xmin, xmax := parseRange(strings.TrimPrefix(parts[0], "x="))
	ymin, ymax := parseRange(strings.TrimPrefix(parts[1], "y="))
	zmin, zmax := parseRange(strings.TrimPrefix(parts[2], "z="))
	min = vec3{
		x: xmin,
		y: ymin,
		z: zmin,
	}
	max = vec3{
		x: xmax,
		y: ymax,
		z: zmax,
	}
	return min, max
}

type reactorTree struct {
	root *reactorNode
}

type reactorNode struct {
	min   int64
	max   int64
	left  *reactorNode
	right *reactorNode
	next  *reactorNode // if x- or y-dimensional
	on    bool         // if z-dimensional
}

func (t *reactorTree) insert(s rebootStep) {
	t.root = t.root.insert(s, 0)
}

func (t *reactorTree) count() int64 {
	return t.root.count(vec3{}, vec3{}, 0)
}

func (n *reactorNode) insert(s rebootStep, dim int) (updated *reactorNode) {
	var min, max int64
	switch dim {
	case 0:
		min, max = s.min.x, s.max.x
	case 1:
		min, max = s.min.y, s.max.y
	case 2:
		min, max = s.min.z, s.max.z
	default:
		panic("bad")
	}
	for min < max {
		ins, upd := n.insert1(min, max)
		min = ins.max
		n = upd
		if dim == 2 {
			ins.on = s.on
		} else {
			ins.next = ins.next.insert(s, dim+1)
		}
	}
	return n
}

func (n *reactorNode) insert1(min, max int64) (inserted, updated *reactorNode) {
	if n == nil {
		ins := &reactorNode{
			min: min,
			max: max,
		}
		return ins, ins
	}
	if min < n.min {
		if n.left == nil {
			ins := &reactorNode{
				min: min,
				max: minim(n.min, max),
			}
			n.left = ins
			return ins, n
		}
		ins, upd := n.left.insert1(min, max)
		ins.max = minim(ins.max, n.min) // don't let child overreach
		n.left = upd
		return ins, n
	}
	if min >= n.max {
		if n.right == nil {
			ins := &reactorNode{
				min: min,
				max: max,
			}
			n.right = ins
			return ins, n
		}
		ins, upd := n.right.insert1(min, max)
		n.right = upd
		return ins, n
	}
	// Split
	if n.min < min {
		n.left = n.clone()
		n.left.max = min
		n.left.right = nil
	}
	if n.max > max {
		n.right = n.clone()
		n.right.min = max
		n.right.left = nil
	}
	n.min = min
	n.max = minim(n.max, max)
	return n, n
}

func (n *reactorNode) clone() *reactorNode {
	if n == nil {
		return nil
	}
	return &reactorNode{
		min:   n.min,
		max:   n.max,
		left:  n.left.clone(),
		right: n.right.clone(),
		next:  n.next.clone(),
		on:    n.on,
	}
}

func (n *reactorNode) count(min, max vec3, dim int) int64 {
	if n == nil {
		return 0
	}
	total := n.left.count(min, max, dim) + n.right.count(min, max, dim)
	switch dim {
	case 0:
		min.x = n.min
		max.x = n.max
		total += n.next.count(min, max, dim+1)
	case 1:
		min.y = n.min
		max.y = n.max
		total += n.next.count(min, max, dim+1)
	case 2:
		min.z = n.min
		max.z = n.max
		if n.on {
			cubeArea := (max.x - min.x) * (max.y - min.y) * (max.z - min.z)
			total += cubeArea
		}
	default:
		panic("bad")
	}
	return total
}
