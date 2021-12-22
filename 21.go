package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(21, problem21)
}

func problem21(ctx *problemContext) {
	var p1, p2 int
	scanner := ctx.scanner()
	i := 0
	for scanner.scan() {
		ns := strings.TrimPrefix(
			scanner.text(),
			fmt.Sprintf("Player %d starting position: ", i+1),
		)
		n := parseInt(ns, 10, 64)
		switch i {
		case 0:
			p1 = int(n)
		case 1:
			p2 = int(n)
		default:
			panic("bad")
		}
		i++
	}
	ctx.reportLoad()

	d := newDetDice(p1, p2)
	for player := 0; ; player = 1 - player {
		if d.turn(player) {
			ctx.reportPart1(d.score[1-player] * d.rolls)
			break
		}
	}

	ctx.reportPart2(maxim(diracDice(p1, p2)))
}

type detDice struct {
	p     [2]int
	score [2]int64
	rolls int64
	die   int
}

func newDetDice(p1, p2 int) *detDice {
	return &detDice{
		p:   [2]int{p1, p2},
		die: 1,
	}
}

func (d *detDice) turn(player int) (win bool) {
	n := d.roll() + d.roll() + d.roll()
	d.p[player] = (d.p[player]+n-1)%10 + 1
	d.score[player] += int64(d.p[player])
	return d.score[player] >= 1000
}

func (d *detDice) roll() int {
	r := d.die
	d.die++
	if d.die > 100 {
		d.die = 1
	}
	d.rolls++
	return r
}

type diracState struct {
	p        [2]int
	score    [2]int
	turn     int // 0 or 1
	cur      int
	numRolls int // 0, 1, 2, 3
}

func diracDice(p1, p2 int) (p1Wins, p2Wins int64) {
	results := make(map[diracState][2]int64)
	st := diracState{p: [2]int{p1, p2}}
	r := diracDiceWinner(st, results)
	return r[0], r[1]
}

func diracDiceWinner(st diracState, results map[diracState][2]int64) (r [2]int64) {
	if r, ok := results[st]; ok {
		return r
	}
	defer func() { results[st] = r }()
	i := st.turn
	if st.numRolls == 3 {
		st.p[i] = (st.p[i]+st.cur-1)%10 + 1
		st.score[i] += st.p[i]
		if st.score[i] >= 21 {
			r[i] = 1
			return r
		}
		st.cur = 0
		st.numRolls = 0
		st.turn = 1 - st.turn
		return diracDiceWinner(st, results)
	}
	for roll := 1; roll <= 3; roll++ {
		st.cur += roll
		st.numRolls++
		w := diracDiceWinner(st, results)
		r[0] += w[0]
		r[1] += w[1]
		st.numRolls--
		st.cur -= roll
	}
	return r
}
