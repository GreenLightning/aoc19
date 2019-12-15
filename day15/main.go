package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	Wall = 0
	Path = 1
)

const (
	North = 1
	South = 2
	West  = 3
	East  = 4
)

var (
	Up    = Vector2{0, -1}
	Down  = Vector2{0, 1}
	Left  = Vector2{-1, 0}
	Right = Vector2{1, 0}
)

var (
	commands   = []int64{North, South, West, East}
	directions = []Vector2{Up, Down, Left, Right}
	reverse    = map[int64]int64{North: South, South: North, West: East, East: West}
	direction  = map[int64]Vector2{North: Up, South: Down, West: Left, East: Right}
)

type QueueItem struct {
	Position Vector2
	Distance int
	Next     *QueueItem
}

var printFlag = flag.Bool("print", false, "print map of the discovered area")

func main() {
	flag.Parse()

	text := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}

	input := make(chan int64)
	output := make(chan int64)
	halt := make(chan bool)

	go emulate(program, input, output, halt)

	var pos Vector2

	grid := make(map[Vector2]int)
	grid[pos] = Path

	var oxygenPos Vector2
	var oxygenDistance int

	// Get a complete map of the area and record the position and distance of the oxygen system.
	{
		var queue []QueueItem
		queue = append(queue, QueueItem{Position: pos, Distance: 0})

		for len(queue) != 0 {
			item := queue[0]
			queue = queue[1:]

			pos = navigate(pos, item.Position, grid, input, output)

			for _, cmd := range commands {
				next, nextDistance := pos.Plus(direction[cmd]), item.Distance+1
				if _, ok := grid[next]; !ok {
					// Try command if we do not know what lies in this direction.
					input <- cmd
					switch <-output {
					case 0:
						grid[next] = Wall
					case 2:
						if oxygenDistance == 0 {
							oxygenDistance = nextDistance
							oxygenPos = next
						}
						fallthrough
					case 1:
						grid[next] = Path
						queue = append(queue, QueueItem{Position: next, Distance: nextDistance})
						// Command succeeded, go back to try other commands.
						input <- reverse[cmd]
						<-output
					}
				}
			}
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(oxygenDistance)
	}

	var maxDistance int

	// Fill complete map and record maximum distance.
	{
		var queue []QueueItem
		queue = append(queue, QueueItem{Position: oxygenPos, Distance: 0})

		visited := make(map[Vector2]bool)
		visited[oxygenPos] = true

		for len(queue) != 0 {
			item := queue[0]
			queue = queue[1:]

			maxDistance = max(maxDistance, item.Distance)

			for _, dir := range directions {
				next := item.Position.Plus(dir)
				if !visited[next] && grid[next] == Path {
					visited[next] = true
					queue = append(queue, QueueItem{Position: next, Distance: item.Distance + 1})
				}
			}
		}
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(maxDistance)
	}

	if *printFlag {
		fmt.Println("----------------")

		var min, max Vector2
		for pos := range grid {
			min = min.Min(pos)
			max = max.Max(pos)
		}

		for y := min.y; y <= max.y; y++ {
			for x := min.x; x <= max.x; x++ {
				pos := Vector2{x, y}
				value, ok := grid[pos]
				if x == 0 && y == 0 {
					fmt.Print("S")
				} else if pos == oxygenPos {
					fmt.Print("O")
				} else if !ok {
					fmt.Print("?")
				} else if value == Path {
					fmt.Print(" ")
				} else {
					fmt.Print("█")
				}
			}
			fmt.Println()
		}
	}
}

// Move from pos to target following only known paths in grid.
// Returns target (i.e. the new position after moving).
func navigate(pos, target Vector2, grid map[Vector2]int, input chan int64, output chan int64) Vector2 {
	var link *QueueItem

	// Find shortest route from target to pos (note reversed order).
	{
		var queue []QueueItem
		queue = append(queue, QueueItem{Position: target, Distance: 0})

		visited := make(map[Vector2]bool)
		visited[target] = true

		for len(queue) != 0 {
			item := queue[0]
			queue = queue[1:]

			if item.Position == pos {
				link = &item
				break
			}

			for _, dir := range directions {
				next := item.Position.Plus(dir)
				if !visited[next] && grid[next] == Path {
					visited[next] = true
					queue = append(queue, QueueItem{Position: next, Distance: item.Distance + 1, Next: &item})
				}
			}
		}
	}

	// Follow path backwards from pos to target and apply correct commands.
	for link.Next != nil {
		for cmd, dir := range direction {
			if link.Position.Plus(dir) == link.Next.Position {
				input <- cmd
				<-output
				break
			}
		}

		link = link.Next
	}

	return target
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
