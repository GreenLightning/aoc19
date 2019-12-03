package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := readLines("input.txt")

	grid := make(map[Vector2]int)

	parseSegment := func(segment string) (dir Vector2, length int) {
		switch segment[0] {
		case 'R':
			dir.x++
		case 'L':
			dir.x--
		case 'U':
			dir.y++
		case 'D':
			dir.y--
		}
		length = toInt(segment[1:])
		return
	}

	pos, steps := Vector2{}, 0
	for _, segment := range strings.Split(lines[0], ",") {
		dir, length := parseSegment(segment)
		for i := 0; i < length; i++ {
			pos = pos.Plus(dir)
			steps++
			grid[pos] = steps
		}
	}

	closestManhatten, closestSteps := math.MaxInt32, math.MaxInt32

	pos, steps = Vector2{}, 0
	for _, segment := range strings.Split(lines[1], ",") {
		dir, length := parseSegment(segment)
		for i := 0; i < length; i++ {
			pos = pos.Plus(dir)
			steps++
			if grid[pos] != 0 {
				distance := pos.ManhattenLength()
				closestManhatten = min(closestManhatten, distance)
				totalSteps := grid[pos] + steps
				closestSteps = min(closestSteps, totalSteps)
			}
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(closestManhatten)
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(closestSteps)
	}
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

func (v Vector2) Minus(other Vector2) Vector2 {
	return Vector2{
		x: v.x - other.x,
		y: v.y - other.y,
	}
}

func (v Vector2) ManhattenLength() int {
	return abs(v.x) + abs(v.y)
}

func (v Vector2) ManhattenDistance(o Vector2) int {
	return v.Minus(o).ManhattenLength()
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(x, y int) int {
	if y < x {
		return y
	}
	return x
}
