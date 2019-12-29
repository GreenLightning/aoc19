package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Layer [5][5]bool

var printFlag = flag.Bool("print", false, "print final state for part two")

func main() {
	flag.Parse()

	lines := readLines("input.txt")

	{
		fmt.Println("--- Part One ---")

		var state uint32

		for y, line := range lines {
			for x, char := range line {
				if char == '#' {
					state |= (1 << uint32(5*y+x))
				}
			}
		}

		seen := make(map[uint32]bool)

		for !seen[state] {
			seen[state] = true

			var next uint32

			for y := 0; y < 5; y++ {
				for x := 0; x < 5; x++ {
					neighbors := 0
					if x > 0 && state&(1<<uint32(5*y+x-1)) != 0 {
						neighbors++
					}
					if x < 4 && state&(1<<uint32(5*y+x+1)) != 0 {
						neighbors++
					}
					if y > 0 && state&(1<<uint32(5*(y-1)+x)) != 0 {
						neighbors++
					}
					if y < 4 && state&(1<<uint32(5*(y+1)+x)) != 0 {
						neighbors++
					}
					var bit uint32 = 1 << uint32(5*y+x)
					if ((state&bit != 0) && neighbors == 1) || ((state&bit == 0) && neighbors >= 1 && neighbors <= 2) {
						next |= bit
					}
				}
			}

			state = next
		}

		fmt.Println(state)
	}

	{
		fmt.Println("--- Part Two ---")

		var layer Layer
		for y, line := range lines {
			for x, char := range line {
				if char == '#' {
					layer[y][x] = true
				}
			}
		}

		state := make(map[int]Layer)
		state[0] = layer
		min, max := 0, 0

		for minute := 0; minute < 200; minute++ {
			next := make(map[int]Layer)

			for index := min - 1; index <= max+1; index++ {
				var nextLayer Layer

				for y := 0; y < 5; y++ {
					for x := 0; x < 5; x++ {
						if x == 2 && y == 2 {
							continue
						}

						neighbors := 0

						if x > 0 && state[index][y][x-1] {
							neighbors++
						}
						if x < 4 && state[index][y][x+1] {
							neighbors++
						}
						if y > 0 && state[index][y-1][x] {
							neighbors++
						}
						if y < 4 && state[index][y+1][x] {
							neighbors++
						}

						if x == 0 && state[index-1][2][1] {
							neighbors++
						}
						if x == 4 && state[index-1][2][3] {
							neighbors++
						}
						if y == 0 && state[index-1][1][2] {
							neighbors++
						}
						if y == 4 && state[index-1][3][2] {
							neighbors++
						}

						if x == 1 && y == 2 {
							for i := 0; i < 5; i++ {
								if state[index+1][i][0] {
									neighbors++
								}
							}
						}

						if x == 3 && y == 2 {
							for i := 0; i < 5; i++ {
								if state[index+1][i][4] {
									neighbors++
								}
							}
						}

						if y == 1 && x == 2 {
							for i := 0; i < 5; i++ {
								if state[index+1][0][i] {
									neighbors++
								}
							}
						}

						if y == 3 && x == 2 {
							for i := 0; i < 5; i++ {
								if state[index+1][4][i] {
									neighbors++
								}
							}
						}

						nextLayer[y][x] = (state[index][y][x] && neighbors == 1) || (!state[index][y][x] && neighbors >= 1 && neighbors <= 2)
					}
				}

				next[index] = nextLayer
			}

			state = next
			min, max = min-1, max+1
		}

		bugs := 0
		for _, layer := range state {
			bugs += count(layer)
		}

		fmt.Println(bugs)

		if *printFlag {
			for count(state[min]) == 0 && min < max {
				min++
			}
			for count(state[max]) == 0 && min < max {
				max--
			}
			for index := min; index <= max; index++ {
				layer := state[index]

				fmt.Printf("Depth %d:\n", index)
				for y := 0; y < 5; y++ {
					for x := 0; x < 5; x++ {
						if y == 2 && x == 2 {
							fmt.Print("?")
						} else if layer[y][x] {
							fmt.Print("#")
						} else {
							fmt.Print(".")
						}
					}
					fmt.Println()
				}
			}
		}
	}
}

func count(layer Layer) (bugs int) {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if layer[y][x] {
				bugs++
			}
		}
	}
	return
}

func readLines(filename string) []string {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
