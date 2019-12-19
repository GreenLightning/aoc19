package main

import (
	"bufio"
	"container/heap"
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

// An entity can be an entrance, a key or a door.
type Entity struct {
	Char     byte
	Index    int
	Position Vector2
	Edges    []Edge
}

type Edge struct {
	Neighbor *Entity
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
			if isEntrance(char) || isKey(char) || isDoor(char) {
				entities[char] = &Entity{
					Char:     char,
					Index:    len(entities),
					Position: Vector2{x, y},
				}
			}
		}
	}

	// Connect the entities, by exploring the grid around each one.
	for _, entity := range entities {
		type Item struct {
			Position Vector2
			Distance int
		}

		var open []Item
		open = append(open, Item{entity.Position, 0})

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
				if char == '.' {
					visited[next] = true
					open = append(open, Item{next, current.Distance + 1})
				} else if isEntrance(char) || isKey(char) || isDoor(char) {
					visited[next] = true
					entity.Edges = append(entity.Edges, Edge{
						Neighbor: entities[char],
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
	return find(entities, positions, 0, allKeys, 0, math.MaxInt32)
}

func find(entities map[byte]*Entity, positions []*Entity, unlocked uint32, allKeys uint32, currentDistance, bestDistance int) int {
	// If we have collected all keys, we have reached the end of the path.
	if isUnlocked(unlocked, allKeys) {
		return currentDistance
	}

	// Find the keys we can reach.
	var keys []Key
	for index, position := range positions {
		keys = append(keys, explore(entities, index, position, unlocked)...)
	}

	// If we cannot reach any more keys, the path is invalid.
	if len(keys) == 0 {
		return math.MaxInt32
	}

	// Try all keys.
	for _, key := range keys {
		newDistance := currentDistance + key.Distance
		if newDistance <= bestDistance {
			newPositions := make([]*Entity, len(positions))
			copy(newPositions, positions)
			newPositions[key.Index] = key.Entity
			newUnlocked := unlocked | bitFromKey(key.Entity.Char)
			distance := find(entities, newPositions, newUnlocked, allKeys, newDistance, bestDistance)
			bestDistance = min(bestDistance, distance)
		}
	}

	return bestDistance
}

type Key struct {
	Entity   *Entity
	Index    int
	Distance int
}

func explore(entities map[byte]*Entity, index int, position *Entity, unlocked uint32) (keys []Key) {
	var open PriorityQueue
	open.Push(&PriorityItem{position, 0})

	visited := make([]bool, len(entities))
	visited[position.Index] = true

	for !open.Empty() {
		current := open.Pop()

		for _, edge := range current.Entity.Edges {
			if visited[edge.Neighbor.Index] {
				continue
			}

			if isEntrance(edge.Neighbor.Char) || isKey(edge.Neighbor.Char) || (isDoor(edge.Neighbor.Char) && isUnlocked(unlocked, bitFromDoor(edge.Neighbor.Char))) {
				visited[edge.Neighbor.Index] = true
				if isKey(edge.Neighbor.Char) && !isUnlocked(unlocked, bitFromKey(edge.Neighbor.Char)) {
					keys = append(keys, Key{edge.Neighbor, index, current.Distance + edge.Distance})
				} else {
					open.Push(&PriorityItem{edge.Neighbor, current.Distance + edge.Distance})
				}
			}
		}
	}

	return
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

type PriorityItem struct {
	Entity   *Entity
	Distance int
}

type PriorityStorage []*PriorityItem

func (s PriorityStorage) Len() int {
	return len(s)
}

func (s PriorityStorage) Less(i, j int) bool {
	return s[i].Distance < s[j].Distance
}

func (s PriorityStorage) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *PriorityStorage) Push(x interface{}) {
	item := x.(*PriorityItem)
	*s = append(*s, item)
}

func (s *PriorityStorage) Pop() interface{} {
	len := len(*s)
	item := (*s)[len-1]
	*s = (*s)[:len-1]
	return item
}

type PriorityQueue struct {
	storage PriorityStorage
}

func (q *PriorityQueue) Len() int {
	return len(q.storage)
}

func (q *PriorityQueue) Empty() bool {
	return len(q.storage) == 0
}

func (q *PriorityQueue) Push(item *PriorityItem) {
	heap.Push(&q.storage, item)
}

func (q *PriorityQueue) Pop() *PriorityItem {
	return heap.Pop(&q.storage).(*PriorityItem)
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
