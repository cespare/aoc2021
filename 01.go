package main

func init() {
	addSolutions(1, problem1)
}

func problem1(ctx *problemContext) {
	var ns []int64
	scanner := ctx.scanner()
	for scanner.scan() {
		ns = append(ns, scanner.int64())
	}
	ctx.reportLoad()

	prev := int64(-1)
	var increases int
	for _, n := range ns {
		if prev >= 0 && n > prev {
			increases++
		}
		prev = n
	}
	ctx.reportPart1(increases)

	prev = int64(-1)
	increases = 0
	for i := 0; i < len(ns)-2; i++ {
		sum := ns[i] + ns[i+1] + ns[i+2]
		if prev >= 0 && sum > prev {
			increases++
		}
		prev = sum
	}
	ctx.reportPart2(increases)
}
