package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

var (
	Up    = Vector2{0, -1}
	Down  = Vector2{0, 1}
	Left  = Vector2{-1, 0}
	Right = Vector2{1, 0}
)

var directions = []Vector2{Up, Down, Left, Right}

// An entity can be an entrance or a key.
type Entity struct {
	Char        byte
	Position    Vector2
	Connections []Connection
}

// A connection represents a path from an entity to a key. Note that there are
// no connections back to an entrance, because we never want to go there
// explicitly (but we might walk over an entrance implicitly by following
// another connection).
type Connection struct {
	Key      *Entity
	Required uint32
	Distance int
}

func main() {
	lines := readLines("input.txt")

	// Convert input into byte arrays, so that we can modify it later.
	grid := make([][]byte, len(lines))
	for y, line := range lines {
		grid[y] = []byte(line)
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(run(grid))
	}

	// Modify grid for part two.
	{
		// Find the original entrance.
		var cx, cy int
		for y, line := range grid {
			for x, char := range line {
				if char == '@' {
					cx, cy = x, y
				}
			}
		}

		// Check surroundings and exit early if the input is not valid (e.g.
		// the examples for part one).
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dy == 0 && dx == 0 {
					continue
				}
				if grid[cy+dy][cx+dx] != '.' {
					fmt.Println("Input not valid for part two.")
					os.Exit(1)
				}
			}
		}

		// The new entrances need different characters, because I want to use
		// the characters as unique map keys.
		grid[cy-1][cx-1], grid[cy-1][cx], grid[cy-1][cx+1] = '@', '#', '$'
		grid[cy+0][cx-1], grid[cy+0][cx], grid[cy+0][cx+1] = '#', '#', '#'
		grid[cy+1][cx-1], grid[cy+1][cx], grid[cy+1][cx+1] = '%', '#', '&'
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(run(grid))
	}
}

func run(grid [][]byte) int {
	numKeys := 0
	entities := make(map[byte]*Entity)

	// Find entities.
	for y, line := range grid {
		for x, char := range line {
			if isKey(char) {
				numKeys++
			}
			if isEntrance(char) || isKey(char) {
				entities[char] = &Entity{
					Char:     char,
					Position: Vector2{x, y},
				}
			}
		}
	}

	// Connect the entities, by exploring the grid around each one.
	for _, entity := range entities {

		type Item struct {
			Position Vector2
			Required uint32
			Distance int
		}

		var open []Item
		open = append(open, Item{entity.Position, 0, 0})

		visited := make(map[Vector2]bool)
		visited[entity.Position] = true

		for len(open) != 0 {
			current := open[0]
			open = open[1:]

			for _, dir := range directions {
				next := current.Position.Plus(dir)
				if visited[next] {
					continue
				}

				char := grid[next.y][next.x]
				if char == '.' || isEntrance(char) {
					// We can walk over these tiles regularly.
					visited[next] = true
					open = append(open, Item{next, current.Required, current.Distance + 1})
				} else if isDoor(char) {
					// If it is a door, we must have the key for the door.
					visited[next] = true
					open = append(open, Item{next, current.Required | bitFromDoor(char), current.Distance + 1})
				} else if isKey(char) {
					// If it is a key, we must have the key to walk over it.
					// This prevents walking over a key without picking it up.
					visited[next] = true
					open = append(open, Item{next, current.Required | bitFromKey(char), current.Distance + 1})

					// We also record the connection to the key,
					// which obviously does not require the key itself.
					entity.Connections = append(entity.Connections, Connection{
						Key:      entities[char],
						Required: current.Required,
						Distance: current.Distance + 1,
					})
				}
			}
		}
	}

	// Find starting positions.
	var positions []*Entity
	for _, entity := range entities {
		if isEntrance(entity.Char) {
			positions = append(positions, entity)
		}
	}

	// Find shortest path.
	var allKeys uint32 = (1 << numKeys) - 1
	return find(positions, 0, allKeys, 0, math.MaxInt32)
}

func find(positions []*Entity, unlocked uint32, allKeys uint32, currentDistance, bestDistance int) int {
	// If we have collected all keys, we have reached the end of the path.
	if isUnlocked(unlocked, allKeys) {
		return currentDistance
	}

	// Try all reachable keys.
	for index, position := range positions {
		for _, conn := range position.Connections {
			if !isUnlocked(unlocked, conn.Required) {
				// We do not yet have all the keys required to walk this way.
				continue
			}

			if isUnlocked(unlocked, bitFromKey(conn.Key.Char)) {
				// We have already collected this key.
				// Do not go there again.
				continue
			}

			newDistance := currentDistance + conn.Distance
			if newDistance > bestDistance {
				continue
			}

			newPositions := make([]*Entity, len(positions))
			copy(newPositions, positions)
			newPositions[index] = conn.Key
			newUnlocked := unlocked | bitFromKey(conn.Key.Char)
			distance := find(newPositions, newUnlocked, allKeys, newDistance, bestDistance)
			bestDistance = min(bestDistance, distance)
		}
	}

	return bestDistance
}

func isEntrance(char byte) bool {
	return char == '@' || char == '$' || char == '%' || char == '&'
}

func isKey(char byte) bool {
	return char >= 'a' && char <= 'z'
}

func isDoor(char byte) bool {
	return char >= 'A' && char <= 'Z'
}

func bitFromKey(char byte) uint32 {
	return 1 << (char - 'a')
}

func bitFromDoor(char byte) uint32 {
	return 1 << (char - 'A')
}

func isUnlocked(unlocked uint32, bits uint32) bool {
	return unlocked&bits == bits
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

func min(x, y int) int {
	if y < x {
		return y
	}
	return x
}
