package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	Up    = Vector2{0, -1}
	Down  = Vector2{0, 1}
	Left  = Vector2{-1, 0}
	Right = Vector2{1, 0}
)

var (
	turnLeft  = map[Vector2]Vector2{Up: Left, Left: Down, Down: Right, Right: Up}
	turnRight = map[Vector2]Vector2{Up: Right, Right: Down, Down: Left, Left: Up}
)

var printFlag = flag.Bool("print", false, "print camera image")

func main() {
	flag.Parse()

	text := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}

	var grid []string
	var width, height int

	// Run the program and extract the camera image into grid.
	{
		input := make(chan int64)
		output := make(chan int64)
		halt := make(chan bool)

		go emulate(program, input, output, halt)

		var builder strings.Builder

	loop:
		for {
			select {
			case char := <-output:
				builder.WriteRune(rune(char))

			case <-halt:
				break loop
			}
		}

		grid = strings.Split(strings.TrimSpace(builder.String()), "\n")
		width, height = len(grid[0]), len(grid)
	}

	if *printFlag {
		for _, line := range grid {
			fmt.Println(line)
		}
	}

	{
		fmt.Println("--- Part One ---")
		sumOfAlignmentParameters := 0
		for y := 1; y+1 < height; y++ {
			for x := 1; x+1 < width; x++ {
				if grid[y][x] == '#' && grid[y-1][x] == '#' && grid[y+1][x] == '#' && grid[y][x-1] == '#' && grid[y][x+1] == '#' {
					sumOfAlignmentParameters += x * y
				}
			}
		}
		fmt.Println(sumOfAlignmentParameters)
	}

	{
		fmt.Println("--- Part Two ---")

		// Wake up the robot.
		program[0] = 2

		// Find the robot.
		var pos, dir Vector2
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				switch grid[y][x] {
				case '^':
					pos = Vector2{x, y}
					dir = Up
				case 'v':
					pos = Vector2{x, y}
					dir = Down
				case '<':
					pos = Vector2{x, y}
					dir = Left
				case '>':
					pos = Vector2{x, y}
					dir = Right
				}
			}
		}

		isScaffold := func(pos Vector2) bool {
			return pos.x >= 0 && pos.y >= 0 && pos.x < width && pos.y < height && grid[pos.y][pos.x] == '#'
		}

		// Gather commands to follow the path.
		// Walk straight for as long as possible, then check if we can turn left or right.
		// If we cannot do either, we have reached the end of the path.
		var path MoveList
		for {
			length := 0
			for isScaffold(pos.Plus(dir)) {
				pos = pos.Plus(dir)
				length++
			}
			if length != 0 {
				path = append(path, strconv.Itoa(length))
			}

			if newDir := turnLeft[dir]; isScaffold(pos.Plus(newDir)) {
				dir = newDir
				path = append(path, "L")
			} else if newDir := turnRight[dir]; isScaffold(pos.Plus(newDir)) {
				dir = newDir
				path = append(path, "R")
			} else {
				break
			}
		}

		result := compressPath(path, []MoveList{path}, nil)
		if len(result) == 0 {
			panic("no solution found")
		}

		input := make(chan int64, 100)
		output := make(chan int64)
		halt := make(chan bool)

		go emulate(program, input, output, halt)

		functions := result[0]
		main := strings.Join(functions[0], ",")
		a := strings.Join(functions[1], ",")
		b := strings.Join(functions[2], ",")
		c := strings.Join(functions[3], ",")

		for _, c := range fmt.Sprintf("%s\n%s\n%s\n%s\nn\n", main, a, b, c) {
			input <- int64(c)
		}

	loop2:
		for {
			select {
			case char := <-output:
				if char >= 128 {
					fmt.Println(char)
				}

			case <-halt:
				break loop2
			}
		}
	}
}

type MoveList []string

func compressPath(path MoveList, fragments []MoveList, functions []MoveList) (result [][]MoveList) {
	if len(functions) == 2 {
		// The last function must be the shortest remaining fragment.
		var lastFunction MoveList
		if len(fragments) != 0 {
			lastFunction = fragments[0]
		}
		for _, fragment := range fragments {
			if len(fragment) < len(lastFunction) {
				lastFunction = fragment
			}
		}

		// Check memory limit.
		if len(strings.Join(lastFunction, ",")) > 20 {
			return nil
		}

		// Each remaining fragment must equal the last function (or multiple copies thereof).
		for _, fragment := range fragments {
			for hasPrefix(fragment, lastFunction) {
				fragment = fragment[len(lastFunction):]
			}
			if len(fragment) != 0 {
				return nil
			}
		}

		newFunctions := make([]MoveList, 0, 3)
		newFunctions = append(newFunctions, functions...)
		newFunctions = append(newFunctions, lastFunction)

		// Replace path with function calls to compute main function.
		var mainFunction MoveList
		for len(path) != 0 {
			for i, function := range newFunctions {
				if hasPrefix(path, function) {
					mainFunction = append(mainFunction, string('A'+i))
					path = path[len(function):]
				}
			}
		}

		// Check memory limit for main function.
		if len(strings.Join(mainFunction, ",")) > 20 {
			return nil
		}

		program := make([]MoveList, 0, 4)
		program = append(program, mainFunction)
		program = append(program, newFunctions...)

		result = append(result, program)
		return
	}

	visited := make(map[string]bool)

	// Collect unique candidates.
	var candidates []MoveList
	for _, fragment := range fragments {
		for length := 1; length <= len(fragment); length++ {
			candidate := fragment[:length]
			text := strings.Join(candidate, ",")
			if len(text) <= 20 && !visited[text] {
				visited[text] = true
				candidates = append(candidates, candidate)
			}
		}
	}

	// Try each candidate.
	for _, candidate := range candidates {
		// Split fragments by candidate.
		var newFragments []MoveList
		for _, fragment := range fragments {
			for {
				i := index(fragment, candidate)
				if i == -1 {
					break
				}
				if i != 0 {
					newFragments = append(newFragments, fragment[:i])
				}
				fragment = fragment[i+len(candidate):]
			}
			if len(fragment) != 0 {
				newFragments = append(newFragments, fragment)
			}
		}

		// Add candidate to functions.
		newFunctions := make([]MoveList, 0, 3)
		newFunctions = append(newFunctions, functions...)
		newFunctions = append(newFunctions, candidate)

		subresult := compressPath(path, newFragments, newFunctions)
		result = append(result, subresult...)
	}

	return
}

func hasPrefix(list, prefix MoveList) bool {
	if len(list) < len(prefix) {
		return false
	}
	for i, move := range prefix {
		if list[i] != move {
			return false
		}
	}
	return true
}

func index(list, sublist MoveList) int {
	for i := 0; i <= len(list)-len(sublist); i++ {
		if hasPrefix(list[i:], sublist) {
			return i
		}
	}
	return -1
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
