package main

func init() {
	addSolutions(3, problem3)
}

func problem3(ctx *problemContext) {
	var ns []int64
	var numBits int
	scanner := ctx.scanner()
	for scanner.scan() {
		if numBits == 0 {
			numBits = len(scanner.text())
		}
		ns = append(ns, parseInt(scanner.text(), 2, 64))
	}
	ctx.reportLoad()

	counts := make([]int, numBits)
	for _, n := range ns {
		for i := 0; i < numBits; i++ {
			if n&(1<<i) > 0 {
				counts[i]++
			} else {
				counts[i]--
			}
		}
	}
	var gamma int64
	for i, c := range counts {
		if c > 0 {
			gamma += int64(1) << i
		}
	}
	epsilon := ^gamma & (int64(1)<<numBits - 1)
	ctx.reportPart1(gamma * epsilon)

	ox := filterDiag(ns, numBits, true)
	co2 := filterDiag(ns, numBits, false)
	ctx.reportPart2(ox * co2)
}

func filterDiag(ns []int64, numBits int, most bool) int64 {
	for i := numBits - 1; ; i-- {
		if len(ns) == 1 {
			return ns[0]
		}
		ns = filterDiag1(ns, i, numBits, most)
	}
	panic("not found")
}

func filterDiag1(ns []int64, i int, numBits int, most bool) []int64 {
	var count int
	for _, n := range ns {
		if n&(1<<i) > 0 {
			count++
		} else {
			count--
		}
	}
	var filtered []int64
	for _, n := range ns {
		if (n&(1<<i) > 0 == (count >= 0)) == most {
			filtered = append(filtered, n)
		}
	}
	return filtered
}
