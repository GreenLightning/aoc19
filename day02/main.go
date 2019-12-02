package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input := readFile("input.txt")

	var program []int
	for _, value := range strings.Split(input, ",") {
		program = append(program, toInt(value))
	}

	{
		fmt.Println("--- Part One ---")
		result, fault := emulate(program, 12, 02)
		if fault {
			panic("unexpected fault")
		}
		fmt.Println(result)
	}

	{
		fmt.Println("--- Part Two ---")
	loop:
		for noun := 0; noun < 100; noun++ {
			for verb := 0; verb < 100; verb++ {
				result, _ := emulate(program, noun, verb)
				if result == 19690720 {
					fmt.Printf("%02d%02d\n", noun, verb)
					break loop
				}
			}
		}
	}
}

func emulate(program []int, noun, verb int) (result int, fault bool) {
	// Copy the program into memory, so that we do not modify the original.
	memory := make([]int, len(program))
	copy(memory, program)

	// Copy inputs into memory.
	memory[1], memory[2] = noun, verb

	ip := 0
	for {
		opcode := memory[ip]
		switch opcode {
		case 1:
			a, b, c := memory[ip+1], memory[ip+2], memory[ip+3]
			memory[c] = memory[a] + memory[b]
			ip += 4
		case 2:
			a, b, c := memory[ip+1], memory[ip+2], memory[ip+3]
			memory[c] = memory[a] * memory[b]
			ip += 4
		case 99:
			return memory[0], false
		default:
			return 0, true
		}
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	check(err)
	return strings.TrimSpace(string(bytes))
}

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	check(err)
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
