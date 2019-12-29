package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

var printFlag = flag.Bool("print", false, "print map as the 200th asteroid is destroyed")

func main() {
	flag.Parse()

	lines := readLines("input.txt")

	asteroids := make(map[Vector2]bool)

	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				asteroids[Vector2{x, y}] = true
			}
		}
	}

	var bestVisible int
	var bestLocation Vector2
	for location := range asteroids {
		visible := len(findVisibleAsteroids(location, asteroids))
		if visible > bestVisible {
			bestVisible = visible
			bestLocation = location
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(bestVisible)
	}

	var vaporizationOrder []Vector2
	for len(asteroids) > 1 {
		// Find the asteroids that will be vaporized during this rotation of the laser.
		list := findVisibleAsteroids(bestLocation, asteroids)

		// Sort them by angle to find the exact order they will be vaproized in.
		calculateAngle := func(asteroid Vector2) float64 {
			dist := asteroid.Minus(bestLocation)
			return 2.0*math.Pi - (math.Atan2(float64(dist.x), float64(dist.y)) + math.Pi)
		}
		sort.Slice(list, func(i, j int) bool {
			return calculateAngle(list[i]) < calculateAngle(list[j])
		})

		// Add them to the global list and remove them from the asteroid field.
		vaporizationOrder = append(vaporizationOrder, list...)
		for _, asteroid := range list {
			delete(asteroids, asteroid)
		}
	}

	target := vaporizationOrder[199]

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(target.x*100 + target.y)
	}

	if *printFlag {
		for _, asteroid := range vaporizationOrder[200:] {
			asteroids[asteroid] = true
		}
		for y, line := range lines {
			for x := range line {
				if x == bestLocation.x && y == bestLocation.y {
					fmt.Print("X")
				} else if x == target.x && y == target.y {
					fmt.Print("O")
				} else if asteroids[Vector2{x, y}] {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
	}
}

func findVisibleAsteroids(location Vector2, asteroids map[Vector2]bool) []Vector2 {
	// Maps from a direction to the asteroid visible in this direction.
	// A direction is the shortest vector with integer coordinates,
	// that the position is a multiple of.
	visible := make(map[Vector2]Vector2)

	for asteroid := range asteroids {
		if asteroid == location {
			continue
		}

		// Calculate distance from location to asteroid and "normalize" it to a
		// direction by dividing by the GCD of the coordinates.
		dist := asteroid.Minus(location)
		dir := dist.DividedBy(gcd(abs(dist.x), abs(dist.y)))

		// If there already is another asteroid in this direction and it is
		// closer than the current asteroid, then we ignore the current
		// asteroid. Otherwise we will replace the other one below.
		if occluder, ok := visible[dir]; ok {
			occluderDist := occluder.Minus(location)
			if (dir.x != 0 && occluderDist.x/dir.x < dist.x/dir.x) || (dir.y != 0 && occluderDist.y/dir.y < dist.y/dir.y) {
				continue
			}
		}

		visible[dir] = asteroid
	}

	var result []Vector2
	for _, asteroid := range visible {
		result = append(result, asteroid)
	}

	return result
}

type Vector2 struct {
	x, y int
}

func (v Vector2) Minus(other Vector2) Vector2 {
	return Vector2{
		x: v.x - other.x,
		y: v.y - other.y,
	}
}

func (v Vector2) DividedBy(factor int) Vector2 {
	return Vector2{
		x: v.x / factor,
		y: v.y / factor,
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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
