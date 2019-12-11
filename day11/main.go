package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(input, ",") {
		program = append(program, toInt64(value))
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(len(emulateEmergencyHullPaintingRobot(program, 0)))
	}

	{
		fmt.Println("--- Part Two ---")

		grid := emulateEmergencyHullPaintingRobot(program, 1)

		var min, max Vector2
		for pos := range grid {
			min = min.Min(pos)
			max = max.Max(pos)
		}

		for y := min.y; y <= max.y; y++ {
			for x := min.x; x <= max.x; x++ {
				if grid[Vector2{x, y}] == 1 {
					fmt.Print("█")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}
	}
}

func emulateEmergencyHullPaintingRobot(program []int64, initialPanel int64) map[Vector2]int64 {
	up := Vector2{0, -1}
	right := Vector2{1, 0}
	down := Vector2{0, 1}
	left := Vector2{-1, 0}

	input := make(chan int64, 1)
	output := make(chan int64)
	halt := make(chan bool)

	go emulate(program, input, output, halt)

	grid := make(map[Vector2]int64)
	pos, dir := Vector2{0, 0}, up

	grid[pos] = initialPanel

	for {
		input <- grid[pos]

		select {
		case value := <-output:
			grid[pos] = value

			if turn := <-output; turn == 1 {
				// turn right
				switch dir {
				case up:
					dir = right
				case right:
					dir = down
				case down:
					dir = left
				case left:
					dir = up
				}
			} else {
				// turn left
				switch dir {
				case up:
					dir = left
				case left:
					dir = down
				case down:
					dir = right
				case right:
					dir = up
				}
			}

			pos = pos.Plus(dir)

		case <-halt:
			return grid
		}
	}
}

func emulate(program []int64, input <-chan int64, output chan<- int64, halt chan<- bool) {
	// Copy the program into memory, so that we do not modify the original.
	memory := make([]int64, len(program))
	copy(memory, program)

	getMemoryPointer := func(index int64) *int64 {
		// Grow memory, if index is out of range.
		for int64(len(memory)) <= index {
			memory = append(memory, 0)
		}
		return &memory[index]
	}

	var ip, relativeBase int64
	for {
		instruction := memory[ip]
		opcode := instruction % 100

		getParameter := func(offset int64) *int64 {
			parameter := memory[ip+offset]
			mode := instruction / pow(10, offset+1) % 10
			switch mode {
			case 0: // position mode
				return getMemoryPointer(parameter)
			case 1: // immediate mode
				return &parameter
			case 2: // relative mode
				return getMemoryPointer(relativeBase + parameter)
			default:
				panic(fmt.Sprintf("fault: invalid parameter mode: ip=%d instruction=%d offset=%d mode=%d", ip, instruction, offset, mode))
			}
		}

		switch opcode {

		case 1: // ADD
			a, b, c := getParameter(1), getParameter(2), getParameter(3)
			*c = *a + *b
			ip += 4

		case 2: // MULTIPLY
			a, b, c := getParameter(1), getParameter(2), getParameter(3)
			*c = *a * *b
			ip += 4

		case 3: // INPUT
			a := getParameter(1)
			*a = <-input
			ip += 2

		case 4: // OUTPUT
			a := getParameter(1)
			output <- *a
			ip += 2

		case 5: // JUMP IF TRUE
			a, b := getParameter(1), getParameter(2)
			if *a != 0 {
				ip = *b
			} else {
				ip += 3
			}

		case 6: // JUMP IF FALSE
			a, b := getParameter(1), getParameter(2)
			if *a == 0 {
				ip = *b
			} else {
				ip += 3
			}

		case 7: // LESS THAN
			a, b, c := getParameter(1), getParameter(2), getParameter(3)
			if *a < *b {
				*c = 1
			} else {
				*c = 0
			}
			ip += 4

		case 8: // EQUAL
			a, b, c := getParameter(1), getParameter(2), getParameter(3)
			if *a == *b {
				*c = 1
			} else {
				*c = 0
			}
			ip += 4

		case 9: // RELATIVE BASE OFFSET
			a := getParameter(1)
			relativeBase += *a
			ip += 2

		case 99: // HALT
			halt <- true
			return

		default:
			panic(fmt.Sprintf("fault: invalid opcode: ip=%d instruction=%d opcode=%d", ip, instruction, opcode))
		}
	}
}

// Integer power: compute a**b using binary powering algorithm
// See Donald Knuth, The Art of Computer Programming, Volume 2, Section 4.6.3
// Source: https://groups.google.com/d/msg/golang-nuts/PnLnr4bc9Wo/z9ZGv2DYxXoJ
func pow(a, b int64) int64 {
	var p int64 = 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func toInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	check(err)
	return result
}

type Vector2 struct {
	x, y int
}

func (v Vector2) Plus(other Vector2) Vector2 {
	return Vector2{
		x: v.x + other.x,
		y: v.y + other.y,
	}
}

func (v Vector2) Min(other Vector2) Vector2 {
	return Vector2{
		x: min(v.x, other.x),
		y: min(v.y, other.y),
	}
}

func (v Vector2) Max(other Vector2) Vector2 {
	return Vector2{
		x: max(v.x, other.x),
		y: max(v.y, other.y),
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	check(err)
	return strings.TrimSpace(string(bytes))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func min(x, y int) int {
	if y < x {
		return y
	}
	return x
}

func max(x, y int) int {
	if y > x {
		return y
	}
	return x
}
