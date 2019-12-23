package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	text := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}

	computers := make([]*Emulator, 50)

	for i := range computers {
		computers[i] = makeEmulator(program, int64(i))
	}

	var natInitialized bool
	var natX, natY int64

	delivered := make(map[int64]bool)

	for {
		for _, computer := range computers {
			for waiting := 0; waiting < 2; {
				address, status := emulate(computer)
				if status == EmulatorStatusOutput {
					x, status := emulate(computer)
					if status != EmulatorStatusOutput {
						panic("expected output")
					}

					y, status := emulate(computer)
					if status != EmulatorStatusOutput {
						panic("expected output")
					}

					if address == 255 {
						if !natInitialized {
							natInitialized = true
							fmt.Println("--- Part One ---")
							fmt.Println(y)
						}
						natX, natY = x, y
					} else {
						target := computers[address]
						target.input = append(target.input, x, y)
					}
					waiting = 0
				} else if status == EmulatorStatusWaitingForInput {
					computer.input = append(computer.input, -1)
					waiting++
				} else {
					panic("halted")
				}
			}
		}

		waiting := 0
		for _, computer := range computers {
			if len(computer.input) == 1 {
				waiting++
			}
		}

		if waiting == len(computers) {
			if delivered[natY] {
				fmt.Println("--- Part Two ---")
				fmt.Println(natY)
				return
			}
			delivered[natY] = true
			computers[0].input = append(computers[0].input, natX, natY)
		}
	}
}

// This version of the intcode emulator does not use goroutines.

type EmulatorStatus int

const (
	EmulatorStatusHalted          EmulatorStatus = 0
	EmulatorStatusOutput          EmulatorStatus = 1
	EmulatorStatusWaitingForInput EmulatorStatus = 2
)

type Emulator struct {
	memory           []int64
	input            []int64
	ip, relativeBase int64
}

func makeEmulator(program []int64, input ...int64) *Emulator {
	// Copy the program into memory, so that we do not modify the original.
	memory := make([]int64, len(program))
	copy(memory, program)

	return &Emulator{
		memory: memory,
		input:  input,
	}
}

func emulate(emulator *Emulator, input ...int64) (int64, EmulatorStatus) {
	emulator.input = append(emulator.input, input...)

	getMemoryPointer := func(index int64) *int64 {
		// Grow memory, if index is out of range.
		for int64(len(emulator.memory)) <= index {
			emulator.memory = append(emulator.memory, 0)
		}
		return &emulator.memory[index]
	}

	for {
		instruction := emulator.memory[emulator.ip]
		opcode := instruction % 100

		getParameter := func(offset int64) *int64 {
			parameter := emulator.memory[emulator.ip+offset]
			mode := instruction / pow(10, offset+1) % 10
			switch mode {
			case 0: // position mode
				return getMemoryPointer(parameter)
			case 1: // immediate mode
				return &parameter
			case 2: // relative mode
				return getMemoryPointer(emulator.relativeBase + parameter)
			default:
				panic(fmt.Sprintf("fault: invalid parameter mode: ip=%d instruction=%d offset=%d mode=%d", emulator.ip, instruction, offset, mode))
			}
		}

		switch opcode {

		case 1: // ADD
			a, b, c := getParameter(1), getParameter(2), getParameter(3)
			*c = *a + *b
			emulator.ip += 4

		case 2: // MULTIPLY
			a, b, c := getParameter(1), getParameter(2), getParameter(3)
			*c = *a * *b
			emulator.ip += 4

		case 3: // INPUT
			if len(emulator.input) == 0 {
				return 0, EmulatorStatusWaitingForInput
			}
			a := getParameter(1)
			*a = emulator.input[0]
			emulator.input = emulator.input[1:]
			emulator.ip += 2

		case 4: // OUTPUT
			a := getParameter(1)
			emulator.ip += 2
			return *a, EmulatorStatusOutput

		case 5: // JUMP IF TRUE
			a, b := getParameter(1), getParameter(2)
			if *a != 0 {
				emulator.ip = *b
			} else {
				emulator.ip += 3
			}

		case 6: // JUMP IF FALSE
			a, b := getParameter(1), getParameter(2)
			if *a == 0 {
				emulator.ip = *b
			} else {
				emulator.ip += 3
			}

		case 7: // LESS THAN
			a, b, c := getParameter(1), getParameter(2), getParameter(3)
			if *a < *b {
				*c = 1
			} else {
				*c = 0
			}
			emulator.ip += 4

		case 8: // EQUAL
			a, b, c := getParameter(1), getParameter(2), getParameter(3)
			if *a == *b {
				*c = 1
			} else {
				*c = 0
			}
			emulator.ip += 4

		case 9: // RELATIVE BASE OFFSET
			a := getParameter(1)
			emulator.relativeBase += *a
			emulator.ip += 2

		case 99: // HALT
			return 0, EmulatorStatusHalted

		default:
			panic(fmt.Sprintf("fault: invalid opcode: ip=%d instruction=%d opcode=%d", emulator.ip, instruction, opcode))
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
