# Advent of Code 2019

This repository contains my solutions to the [Advent of Code](https://adventofcode.com/) puzzles of 2019.

**SPOILER WARNING**: Do NOT look at the solutions until you have solved the puzzles for yourself.

## Previous Years

- [2018](https://github.com/GreenLightning/aoc18)
- [2017](https://github.com/GreenLightning/aoc17)
- [2016](https://github.com/GreenLightning/aoc16)
- [2015](https://github.com/GreenLightning/aoc15)

## Methodology

I am a competitive person, so I try to solve each puzzle as fast as possible
and hopefully land a spot on the
[leaderboard](https://adventofcode.com/leaderboard) (but there are [other
ways](https://adventofcode.com/about) to play as well).

As you might know, each puzzle has two parts. Once I have solved both, I make
them work at the same time (I might first modify the solution of the first
part for the second part), format the code, give the variables more reasonable
names and occasionally write a few comments or do some light refactoring as I
see fit. That is the code that you can see here. Note that it is not the
unformatted mess that produced the solution that I submitted, but it is not
the most well-thought-out code either. Occasionally I will come back to a
problem and work on my solution a little more (for example to optimize its
runtime).

To go fast, I have written myself an input downloader,
[aocdl](https://github.com/GreenLightning/advent-of-code-downloader). This is
a command line utility that you can start before the puzzle is released and it
will automatically download your personal input file once it is available.
That way you do not lose any time dealing with the puzzle input.

I also have created some template code, which I copy into the solutions. I
copy the template code instead of importing a library to keep each solution
self-contained and independent. The special template code for priority queues
and different vector types needs to be adapted to each puzzle anyway.
