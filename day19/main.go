package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var printFlag = flag.Bool("print", false, "print beam")

var program []int64

func main() {
	flag.Parse()

	text := readFile("input.txt")

	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}

	if *printFlag {
		for y := 0; y < 80; y++ {
			for x := 0; x < 80; x++ {
				if probe(x, y) {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
	}

	{
		fmt.Println("--- Part One ---")

		count := 0

		for y := 0; y < 50; y++ {
			for x := 0; x < 50; x++ {
				if probe(x, y) {
					count++
				}
			}
		}

		fmt.Println(count)
	}

	{
		fmt.Println("--- Part Two ---")

		startX, startY := 0, 0
		for {
			if !probe(startX, startY) {
				startX++
			}

			x, y := startX, startY

			for {
				if probe(x, y) && probe(x+99, y) && probe(x, y+99) {
					fmt.Println(x*10000 + y)
					return
				}

				x++

				if !probe(x+99, y) {
					startY++
					break
				}
			}
		}
	}
}

func probe(x, y int) bool {
	input := make(chan int64)
	output := make(chan int64)
	halt := make(chan bool, 1)

	go emulate(program, input, output, halt)

	input <- int64(x)
	input <- int64(y)

	return <-output == 1
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
