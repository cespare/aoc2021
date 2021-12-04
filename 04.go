package main

import (
	"strings"
)

func init() {
	addSolutions(4, problem4)
}

func problem4(ctx *problemContext) {
	var nums []int64
	var boards [][][]int64
	scanner := ctx.scanner()
	for scanner.scan() {
		if len(nums) == 0 {
			for _, f := range strings.Split(scanner.text(), ",") {
				nums = append(nums, parseInt(f, 10, 64))
			}
			continue
		}
		if scanner.text() == "" {
			boards = append(boards, nil)
			continue
		}
		var row []int64
		for _, f := range strings.Fields(scanner.text()) {
			row = append(row, parseInt(f, 10, 64))
		}
		boards[len(boards)-1] = append(boards[len(boards)-1], row)
	}
	ctx.reportLoad()

	bingoState := make([][][]bool, len(boards))
	for i := range bingoState {
		bingoState[i] = make([][]bool, 5)
		for j := range bingoState[i] {
			bingoState[i][j] = make([]bool, 5)
		}
	}

	alreadyWon := make([]bool, len(boards))
	foundFirst := false
	var last int64
	for _, n := range nums {
		for i, board := range boards {
			if alreadyWon[i] {
				continue
			}
			for j, row := range board {
				for k, val := range row {
					if val == n {
						bingoState[i][j][k] = true
					}
				}
			}
			if isBingo(bingoState[i]) {
				var score int64
				for j, row := range board {
					for k, val := range row {
						if !bingoState[i][j][k] {
							score += val
						}
					}
				}
				score *= n
				if !foundFirst {
					ctx.reportPart1(score)
					foundFirst = true
				}
				last = score
				alreadyWon[i] = true
			}
		}
	}
	ctx.reportPart2(last)
}

func isBingo(state [][]bool) bool {
rowLoop:
	for _, row := range state {
		for _, v := range row {
			if !v {
				continue rowLoop
			}
		}
		return true
	}
colLoop:
	for k := 0; k < 5; k++ {
		for _, row := range state {
			if !row[k] {
				continue colLoop
			}
		}
		return true
	}

	return false
}
