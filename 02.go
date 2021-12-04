package main

func init() {
	addSolutions(2, problem2)
}

func problem2(ctx *problemContext) {
	var pos1 vec2
	var pos2 vec3
	scanner := ctx.scanner()
	for scanner.scan() {
		line := scanner.text()
		if v, ok := trimPrefix(line, "forward "); ok {
			pos1.x += parseInt(v, 10, 64)
			pos2.x += parseInt(v, 10, 64)
			pos2.y += pos2.z * parseInt(v, 10, 64)
			continue
		}
		if v, ok := trimPrefix(line, "up "); ok {
			pos1.y -= parseInt(v, 10, 64)
			pos2.z -= parseInt(v, 10, 64)
			continue
		}
		if v, ok := trimPrefix(line, "down "); ok {
			pos1.y += parseInt(v, 10, 64)
			pos2.z += parseInt(v, 10, 64)
			continue
		}
		panic(line)
	}
	ctx.reportLoad()

	ctx.reportPart1(pos1.x * pos1.y)
	ctx.reportPart2(pos2.x * pos2.y)
}
