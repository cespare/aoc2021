package main

import (
	"fmt"
	"strconv"
	"strings"
)

func init() {
	addSolutions(18, problem18)
}

func problem18(ctx *problemContext) {
	var nums []sfNum
	scanner := ctx.scanner()
	for scanner.scan() {
		num, rem := parseSFNum(scanner.text())
		if rem != "" {
			panic("extra")
		}
		nums = append(nums, num)
	}
	ctx.reportLoad()

	var sum sfNum
	for _, num := range nums {
		if sum == nil {
			sum = num
		} else {
			sum = addSFNums(sum, num)
		}
	}
	// fmt.Println(sum)

	ctx.reportPart1(sum.mag())

	var max int64
	for i, num0 := range nums {
		for j, num1 := range nums {
			if i == j {
				continue
			}
			mag := addSFNums(num0, num1).mag()
			if mag > max {
				max = mag
			}
		}
	}
	ctx.reportPart2(max)
}

type sfNum interface {
	setAsLeftChild(sfNum)
	setAsRightChild(sfNum)
	setSibling(sfNum)

	String() string
	explode(int) bool
	split() bool
	bubbleUp(v int64, left bool)
	bubbleDown(v int64, left bool)
	mag() int64
	clone() sfNum
}

type sfNode struct {
	parent  sfNum
	sibling sfNum
	isLeft  bool // left or right of parent
}

func (n *sfNode) setAsLeftChild(p sfNum) {
	n.parent = p
	n.isLeft = true
}

func (n *sfNode) setAsRightChild(p sfNum) {
	n.parent = p
	n.isLeft = false
}

func (n *sfNode) setSibling(sibling sfNum) {
	n.sibling = sibling
}

type sfLit struct {
	sfNode
	v int64
}

type sfPair struct {
	sfNode
	left  sfNum
	right sfNum
}

func parseSFNum(s string) (num sfNum, rem string) {
	if s[0] != '[' {
		panic("bad")
	}
	s = s[1:]
	n := new(sfPair)
	if s[0] == '[' {
		n.left, s = parseSFNum(s)
	} else {
		var v int64
		v, s = parseIntForward(s)
		n.left = &sfLit{v: v}
	}
	if s[0] != ',' {
		panic("bad")
	}
	s = s[1:]
	if s[0] == '[' {
		n.right, s = parseSFNum(s)
	} else {
		var v int64
		v, s = parseIntForward(s)
		n.right = &sfLit{v: v}
	}
	n.fixChildPointers()
	if s[0] != ']' {
		panic("bad")
	}
	s = s[1:]
	return n, s
}

func parseIntForward(s string) (int64, string) {
	i := strings.IndexFunc(s, func(r rune) bool { return r < '0' || r > '9' })
	return parseInt(s[:i], 10, 64), s[i:]
}

func (n *sfLit) String() string {
	return strconv.FormatInt(n.v, 10)
}

func (n *sfPair) String() string {
	return fmt.Sprintf("[%s,%s]", n.left, n.right)
}

func addSFNums(n0, n1 sfNum) sfNum {
	n0 = n0.clone()
	n1 = n1.clone()
	n := &sfPair{
		left:  n0,
		right: n1,
	}
	n.fixChildPointers()
	// fmt.Println("after addition:", n)

	// Reduce

	for {
		if n.explode(0) {
			// fmt.Println("after explode:", n)
			continue
		}
		if n.split() {
			// fmt.Println("after split:", n)
			continue
		}
		return n
	}
}

func (n *sfLit) explode(int) bool { return false }

func (n *sfPair) explode(depth int) bool {
	if depth > 4 {
		panic("unexpected depth")
	}
	if depth == 4 {
		n.bubbleUp(n.left.(*sfLit).v, true)
		n.bubbleUp(n.right.(*sfLit).v, false)
		lit := &sfLit{
			sfNode: n.sfNode,
			v:      0,
		}
		parent := lit.parent.(*sfPair)
		if lit.isLeft {
			parent.left = lit
		} else {
			parent.right = lit
		}
		lit.sibling.setSibling(lit)
		return true
	}
	if n.left.explode(depth + 1) {
		return true
	}
	if n.right.explode(depth + 1) {
		return true
	}
	return false
}

func (n *sfLit) bubbleUp(int64, bool) { panic("unexpected") }

func (n *sfLit) bubbleDown(v int64, _ bool) { n.v += v }

func (n *sfPair) bubbleUp(v int64, left bool) {
	if n.parent == nil {
		return
	}
	if left != n.isLeft {
		n.sibling.bubbleDown(v, !left)
		return
	}
	n.parent.bubbleUp(v, left)
}

func (n *sfPair) bubbleDown(v int64, left bool) {
	if left {
		n.left.bubbleDown(v, left)
	} else {
		n.right.bubbleDown(v, left)
	}
}

func (n *sfLit) split() bool {
	if n.v < 10 {
		return false
	}
	pair := &sfPair{
		sfNode: n.sfNode,
		left:   &sfLit{v: n.v / 2},
		right:  &sfLit{v: (n.v + 1) / 2},
	}
	pair.fixChildPointers()
	parent := pair.parent.(*sfPair)
	if pair.isLeft {
		parent.left = pair
	} else {
		parent.right = pair
	}
	return true
}

func (n *sfPair) split() bool {
	return n.left.split() || n.right.split()
}

func (n *sfLit) mag() int64 {
	return n.v
}

func (n *sfPair) mag() int64 {
	return 3*n.left.mag() + 2*n.right.mag()
}

func (n *sfLit) clone() sfNum {
	n1 := *n
	n1.parent = nil
	return &n1
}

func (n *sfPair) clone() sfNum {
	n1 := *n
	n1.parent = nil
	n1.left = n1.left.clone()
	n1.right = n1.right.clone()
	n1.fixChildPointers()
	return &n1
}

func (n *sfPair) fixChildPointers() {
	n.left.setAsLeftChild(n)
	n.right.setAsRightChild(n)
	n.left.setSibling(n.right)
	n.right.setSibling(n.left)
}
