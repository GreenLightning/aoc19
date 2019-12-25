package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Room struct {
	Name        string
	Connections map[string]*Room
}

type Mode int

const (
	ModeExplore  Mode = 0
	ModeNavigate Mode = 1
	ModeTest     Mode = 2
)

var opposite = map[string]string{"north": "south", "south": "north", "west": "east", "east": "west"}

func main() {
	playFlag := flag.Bool("play", false, "play the game yourself")
	interactiveFlag := flag.Bool("interactive", false, "press enter to advance")

	flag.Parse()

	text := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}

	emulator := makeEmulator(program)
	scanner := bufio.NewScanner(os.Stdin)

	if *playFlag {
		for {
			char, status := emulate(emulator)
			switch status {
			case EmulatorStatusHalted:
				return
			case EmulatorStatusOutput:
				fmt.Print(string(char))
				if char == '\n' {
					time.Sleep(32 * time.Millisecond)
				}
			case EmulatorStatusWaitingForInput:
				if scanner.Scan() {
					emulator.WriteString(scanner.Text())
					emulator.WriteString("\n")
				}
			}
		}
	}

	sendCommand := func(format string, args ...interface{}) {
		cmd := fmt.Sprintf(format, args...)
		if *interactiveFlag {
			fmt.Print(cmd)
		}
		emulator.WriteString(cmd)
	}

	roomNameRegex := regexp.MustCompile(`^== (.+) ==$`)
	listItemRegex := regexp.MustCompile(`^- (.+)$`)
	takenRegex := regexp.MustCompile(`^You take the (.+)\.$`)
	droppedRegex := regexp.MustCompile(`^You drop the (.+)\.$`)

	world := make(map[string]*Room)
	inventory := make(map[string]bool)

	var mode Mode
	var path []*Room
	var checkpoint, floor *Room
	var testDir string

	var availableItems []string
	var itemMask uint64

	var last *Room
	var lastItems []string
	var lastDir string

	var outputBuilder strings.Builder

loop:
	for {
		char, status := emulate(emulator)
		switch status {
		case EmulatorStatusHalted:
			output := outputBuilder.String()
			outputBuilder.Reset()

			var result string

			resultRegex := regexp.MustCompile(`"Oh, hello! You should be able to get in by typing (\d+) on the keypad at the main airlock\."$`)

			for _, line := range strings.Split(output, "\n") {
				if match := resultRegex.FindStringSubmatch(line); match != nil {
					result = match[1]
				}
			}

			if !*interactiveFlag {
				fmt.Println("--- Part One ---")
				fmt.Println(result)
			}

			return

		case EmulatorStatusOutput:
			if *interactiveFlag {
				fmt.Print(string(char))
			}

			outputBuilder.WriteString(string(char))

			if *interactiveFlag && char == '\n' {
				time.Sleep(32 * time.Millisecond)
			}

		case EmulatorStatusWaitingForInput:
			output := outputBuilder.String()
			outputBuilder.Reset()

			var current *Room
			var items []string

			lines := strings.Split(output, "\n")
			for i := 0; i < len(lines); i++ {
				line := lines[i]

				if line == "" || line == "Command?" {
					continue
				}

				if match := roomNameRegex.FindStringSubmatch(line); match != nil {
					name := match[1]

					var description []string
					for ; i+1 < len(lines) && lines[i+1] != ""; i++ {
						description = append(description, lines[i+1])
					}

					current = world[name]
					if current == nil {
						current = &Room{Name: name}
						world[name] = current
					}

					items = nil

					continue
				}

				if line == "Doors here lead:" {
					fresh := (current.Connections == nil)

					if fresh {
						current.Connections = make(map[string]*Room)
					}

					for ; i+1 < len(lines) && lines[i+1] != ""; i++ {
						match := listItemRegex.FindStringSubmatch(lines[i+1])
						direction := match[1]
						if fresh {
							current.Connections[direction] = nil
						}
					}

					continue
				}

				if line == "Items here:" {
					for ; i+1 < len(lines) && lines[i+1] != ""; i++ {
						match := listItemRegex.FindStringSubmatch(lines[i+1])
						item := match[1]
						items = append(items, item)
					}

					continue
				}

				if match := takenRegex.FindStringSubmatch(line); match != nil {
					taken := match[1]
					inventory[taken] = true

					current = last
					for _, item := range lastItems {
						if item != taken {
							items = append(items, item)
						}
					}

					continue
				}

				if match := droppedRegex.FindStringSubmatch(line); match != nil {
					dropped := match[1]
					inventory[dropped] = false

					current = last
					items = append(lastItems, dropped)

					continue
				}

				if strings.HasPrefix(line, `A loud, robotic voice says "Alert!`) {
					if mode == ModeExplore {
						path = path[:len(path)-1]
						checkpoint, floor, testDir = last, current, lastDir
						checkpoint.Connections[testDir] = floor
					}

					last, lastItems, lastDir = nil, nil, ""

					continue
				}

				panic(line)
			}

			if *interactiveFlag {
				if mode == ModeExplore {
					scanner.Scan()
				} else {
					time.Sleep(32 * time.Millisecond)
				}
			}

			if last != nil && lastDir != "" && last.Connections[lastDir] == nil {
				last.Connections[lastDir] = current
				current.Connections[opposite[lastDir]] = last
			}

			last, lastItems, lastDir = current, items, ""

			switch mode {
			case ModeExplore:

				blacklist := []string{
					"photons",
					"escape pod",
					"molten lava",
					"infinite loop",
					"giant electromagnet",
				}

			itemLoop:
				for _, item := range items {
					for _, bad := range blacklist {
						if item == bad {
							continue itemLoop
						}
					}

					sendCommand("take %s\n", item)
					continue loop
				}

				var target string
				for dir, room := range current.Connections {
					if room == nil {
						path = append(path, current)
						target = dir
						break
					}
				}

				if target == "" && len(path) != 0 {
					last := path[len(path)-1]
					for dir, room := range current.Connections {
						if room == last {
							path = path[:len(path)-1]
							target = dir
							break
						}
					}
					if target == "" {
						panic(fmt.Sprintf(`cannot go from "%s" to "%s"`, current.Name, last.Name))
					}
				}

				if target != "" {
					lastDir = target
					sendCommand("%s\n", target)
					continue loop
				}

				path = findPath(current, checkpoint)[1:]
				mode = ModeNavigate
				fallthrough

			case ModeNavigate:
				if len(path) != 0 {
					for dir, room := range current.Connections {
						if room == path[0] {
							path = path[1:]
							sendCommand("%s\n", dir)
							continue loop
						}
					}

					panic(fmt.Sprintf(`cannot go from "%s" to "%s"`, current.Name, path[0].Name))
				}

				availableItems = nil
				for item := range inventory {
					availableItems = append(availableItems, item)
				}
				itemMask = 0
				mode = ModeTest
				fallthrough

			case ModeTest:
				for index := 0; index < len(availableItems); index++ {
					item := availableItems[index]
					targetState := (itemMask&(1<<uint64(index)) != 0)
					if inventory[item] != targetState {
						var action string
						if targetState {
							action = "take"
						} else {
							action = "drop"
						}
						sendCommand("%s %s\n", action, item)
						continue loop
					}
				}

				itemMask++
				sendCommand("%s\n", testDir)
				continue loop
			}
		}
	}
}

func findPath(from, to *Room) []*Room {
	type Item struct {
		Room *Room
		Path []*Room
	}

	var queue []Item
	queue = append(queue, Item{from, []*Room{from}})

	visited := make(map[*Room]bool)
	visited[from] = true

	for len(queue) != 0 {
		item := queue[0]
		queue = queue[1:]

		if item.Room == to {
			return item.Path
		}

		for _, next := range item.Room.Connections {
			if !visited[next] {
				visited[next] = true
				path := make([]*Room, len(item.Path), len(item.Path)+1)
				copy(path, item.Path)
				path = append(path, next)
				queue = append(queue, Item{next, path})
			}
		}
	}

	return nil
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

func (emulator *Emulator) WriteString(s string) (int, error) {
	for _, char := range s {
		emulator.input = append(emulator.input, int64(char))
	}
	return len(s), nil
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
