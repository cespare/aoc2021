package main

func init() {
	addSolutions(24, problem24)
}

func problem24(ctx *problemContext) {
	ctx.reportLoad()
	ctx.reportPart1(65984919997939)
	ctx.reportPart2(11211619541713)
}

/*
14 blocks of (with 3 constants: S, A, B):

inp w
mul x 0
add x z
mod x 26
div z S
add x A
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y B
mul y x
add z y

In pseudocode:

w = input digit
x = z%26
z = z/S
x = x+A
x = x==w ? 0 : 1
y = x*25 + 1
z = z*y
y = (w+B) * x
z = z+y

Note: only the z value is carried from block to block.

My puzzle values are:

1.  S=1  A=12  B=7
2.  S=1  A=11  B=15
3.  S=1  A=12  B=2
4.  S=26 A=-3  B=15
5.  S=1  A=10  B=14
6.  S=26 A=-9  B=2
7.  S=1, A=10  B=15
8.  S=26 A=-7  B=1
9.  S=26 A=-11 B=15
10. S=26 A=-4  B=15
11. S=1  A=14  B=12
12. S=1  A=11  B=2
13. S=26 A=-8  B=13
14. S=26 A=-10 B=13

Note: S==1 when A>0, S==26 when A<0.

So there are really two different kinds of blocks, one where S is 1
and one where S is 26. There are 7 of each kind.

For S=1, note that A>=10. The x==w check cannot return true for any input digit
in [1, 9].

Simplifying the pseudocode as much as possible, the S=1 block is:

w = input digit
z = z*26 + w + B

For the S=26 block, there are two cases. If the input digit makes the equality
test work (that is, w == z%26+A):

w = input digit
z = z/26

Otherwise, if w != x%26+A:

w = input digit
z = z + w + B

(The z/26 is canceled out by z*26.)

We can think of the puzzle as manipulating a 14-digit base-26 number, one digit
at a time, where the digits are the values of z. (This works because B<=15,
so the value of w+B we add on is less than 26.)

In the S=1 case, we are adding a digit (w+B) to the number.

In the S=26 case, we either drop the last digit (if w == z%26+A) or fail to do
so. But because we have 7 each of S==1 and S==26, we *must* drop a digit at each
opportunity if we want to get to z=0 at the end.

So now there are just two kinds of steps in a possible solution:

S=1:  z = z*26 + w + B
S=26: z = z/26

Let's call the values of z z1, z2, ..., z14.
Similarly, the input digits will be d1, d2, ..., d14.

For the w == z%26+A condition to hold, it means that (for example)
d4 = (last digit of z3) + A3

I'll represent base-26 digits separated by colons rather than writing out a
polynomial.

z1  = d1+7
z2  = d1+7:d2+15
z3  = d1+7:d2+15:d3+2
z4  = d1+7:d2+15          d4 = d3+2-3  =>  d4 = d3-1
z5  = d1+7:d2+15:d5+14
z6  = d1+7:d2+15          d6 = d5+14-9  =>  d6 = d5+5
z7  = d1+7:d2+15:d7+15
z8  = d1+7:d2+15          d8 = d7+15-7  =>  d8 = d7+8
z9  = d1+7                d9 = d2+15-11  =>  d9 = d2+4
z10 = 0                   d10 = d1+7-4  =>  d10 = d1+3
z11 = d11+12
z12 = d11+12:d12+2
z13 = d11+12              d13 = d12+2-8  =>  d13 = d12-6
z14 = 0                   d14 = d11+12-10  =>  d14 = d11+2

Thus, the constraints on the digits are as follows:

d4 = d3 - 1
d6 = d5 + 5
d8 = d7 + 8
d9 = d2 + 4
d10 = d1 + 3
d13 = d12 - 6
d14 = d11 + 2

Reordering/rearranging/duplicating:

d1 = d10 - 3
d2 = d9 - 4
d3 = d4 + 1
d4 = d3 - 1
d5 = d6 - 5
d6 = d5 + 5
d7 = d8 - 8
d8 = d7 + 8
d9 = d2 + 4
d10 = d1 + 3
d11 = d14 - 2
d12 = d13 + 6
d13 = d12 - 6
d14 = d11 + 2

From here we can easily choose the largest possible digits.

d1:  6
d2:  5
d3:  9
d4:  8
d5:  4
d6:  9
d7:  1
d8:  9
d9:  9
d10: 9
d11: 7
d12: 9
d13: 3
d14: 9

So the part 1 solution is 65984919997939.

Similarly, the smallest possible number is:

d1:  1
d2:  1
d3:  2
d4:  1
d5:  1
d6:  6
d7:  1
d8:  9
d9:  5
d10: 4
d11: 1
d12: 7
d13: 1
d14: 3

So the part 2 solution is 11211619541713.
*/
