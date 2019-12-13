package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	Empty  = 0
	Wall   = 1
	Block  = 2
	Paddle = 3
	Ball   = 4
)

// Warning: For my input, this outputs about 150k lines.
var printFlag = flag.Bool("print", false, "print game state before each input is provided")

func main() {
	flag.Parse()

	input := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(input, ",") {
		program = append(program, toInt64(value))
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(countBlocks(program))
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(emulateArcadeCabinet(program))
	}
}

func countBlocks(program []int64) (count int) {
	input := make(chan int64)
	messages := make(chan Message)

	go emulate(program, input, messages)

	grid := make(map[Vector2]int64)

	for {
		message := <-messages
		switch message.Kind {
		case MessageOutput:
			var pos Vector2
			pos.x = int(message.Value)

			message = <-messages
			if message.Kind != MessageOutput {
				panic("unexpected message")
			}
			pos.y = int(message.Value)

			message = <-messages
			if message.Kind != MessageOutput {
				panic("unexpected message")
			}
			grid[pos] = message.Value

		case MessageHalt:
			for _, tile := range grid {
				if tile == Block {
					count++
				}
			}
			return

		default:
			panic("unexpected message")
		}
	}
}

func emulateArcadeCabinet(program []int64) int64 {
	// Insert quarters.
	program[0] = 2

	input := make(chan int64)
	messages := make(chan Message)

	go emulate(program, input, messages)

	grid := make(map[Vector2]int64)
	var score int64

	for {
		message := <-messages
		switch message.Kind {
		case MessageWaitingForInput:
			if *printFlag {
				var min, max Vector2
				for pos := range grid {
					min = min.Min(pos)
					max = max.Max(pos)
				}

				for y := min.y; y <= max.y; y++ {
					for x := min.x; x <= max.x; x++ {
						switch grid[Vector2{x, y}] {
						case Empty:
							fmt.Print(" ")
						case Wall:
							fmt.Print("█")
						case Block:
							fmt.Print("X")
						case Paddle:
							fmt.Print("-")
						case Ball:
							fmt.Print("O")
						}
					}
					fmt.Println()
				}
				fmt.Println("Score: ", score)
			}

			// Find the ball and the paddle, then move the paddle closer to the ball.
			// Once they have the same x position, this will track the ball perfectly,
			// since they both move at the same speed (1 tile / frame).
			var ball, paddle Vector2
			for pos, tile := range grid {
				switch tile {
				case Ball:
					ball = pos
				case Paddle:
					paddle = pos
				}
			}
			input <- int64(sign(ball.x - paddle.x))

		case MessageOutput:
			var pos Vector2
			pos.x = int(message.Value)

			message = <-messages
			if message.Kind != MessageOutput {
				panic("unexpected message")
			}
			pos.y = int(message.Value)

			message = <-messages
			if message.Kind != MessageOutput {
				panic("unexpected message")
			}
			if pos.x == -1 && pos.y == 0 {
				score = message.Value
			} else {
				grid[pos] = message.Value
			}

		case MessageHalt:
			return score

		default:
			panic("unexpected message")
		}
	}
}

// This version of the intcode emulator is synchronous, i.e. it sends a
// MessageWaitingForInput before reading from the input channel. This makes it
// possible to base the input on the previous output, without requiring the
// controlling code to know the exact behavior of the intcode program. Using
// channels in this way allows for a clean separation between the intcode
// emulator and the puzzle specific controlling code.

const (
	MessageWaitingForInput = iota
	MessageOutput
	MessageHalt
)

type Message struct {
	Kind  int
	Value int64
}

func emulate(program []int64, input <-chan int64, messages chan<- Message) {
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
			messages <- Message{Kind: MessageWaitingForInput}
			a := getParameter(1)
			*a = <-input
			ip += 2

		case 4: // OUTPUT
			a := getParameter(1)
			messages <- Message{Kind: MessageOutput, Value: *a}
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
			messages <- Message{Kind: MessageHalt}
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

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
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
