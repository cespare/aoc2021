package main

func init() {
	addSolutions(10, problem10)
}

func problem10(ctx *problemContext) {
	var navLines []string
	scanner := ctx.scanner()
	for scanner.scan() {
		navLines = append(navLines, scanner.text())
	}

	ctx.reportLoad()

	var stacks [][]byte
	var score int64
	for _, line := range navLines {
		fail, stack := checkNavLine(line)
		if fail == 0 {
			stacks = append(stacks, stack)
			continue
		}
		switch fail {
		case 0:
		case ')':
			score += 3
		case ']':
			score += 57
		case '}':
			score += 1197
		case '>':
			score += 25137
		}
	}

	ctx.reportPart1(score)

	var scores []int64
	for _, stack := range stacks {
		var score int64
		for i := len(stack) - 1; i >= 0; i-- {
			score *= 5
			switch stack[i] {
			case '(':
				score += 1
			case '[':
				score += 2
			case '{':
				score += 3
			case '<':
				score += 4
			}
		}
		scores = append(scores, score)
	}
	SliceSort(scores)
	ctx.reportPart2(scores[len(scores)/2])
}

func checkNavLine(s string) (fail byte, stack []byte) {
	for i := 0; i < len(s); i++ {
		c := s[i]
		want, ok := navPair[c]
		if !ok {
			stack = append(stack, c)
			continue
		}
		if len(stack) == 0 || stack[len(stack)-1] != want {
			return c, nil
		}
		stack = stack[:len(stack)-1]
	}
	return 0, stack
}

var navPair = map[byte]byte{
	')': '(',
	']': '[',
	'}': '{',
	'>': '<',
}
