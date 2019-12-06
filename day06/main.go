package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	lines := readLines("input.txt")

	// A map from the object in orbit to the object it is orbiting around.
	// This stores the parent-of relationship in the orbit tree.
	orbits := make(map[string]string)

	for _, line := range lines {
		parts := strings.Split(line, ")")
		orbits[parts[1]] = parts[0]
	}

	{
		fmt.Println("--- Part One ---")

		// > What is the total number of direct and indirect orbits in your
		// > map data?

		// For each object in orbit, we walk up the orbit tree until we reach
		// its root (the COM) and count how many other objects the object
		// orbits. The first iteration of the inner loop will be a direct
		// orbit, while every other iteration will be an indirect orbit.

		total := 0
		for object := range orbits {
			for {
				parent, ok := orbits[object]
				if !ok {
					break
				}
				object = parent
				total++
			}
		}

		fmt.Println(total)
	}

	{
		fmt.Println("--- Part Two ---")

		// > What is the minimum number of orbital transfers required to move
		// > from the object YOU are orbiting to the object SAN is orbiting?

		// Or, put another way, what is the shortest path between the object
		// YOU is orbiting to the object SAN is orbiting? Since this is a
		// tree, the shortest path consists of two parts. The first part walks
		// up the tree from the object YOU is orbiting to the first common
		// parent, while the second part walks down to the object SAN is
		// orbiting.

		// We do not know what the first common parent is, so we first walk up
		// from the object YOU is orbiting to the root of the tree. At each
		// step we will store the distance to the current object in the "path"
		// map below.

		path := make(map[string]int)

		object, distance := orbits["YOU"], 0
		for {
			path[object] = distance
			parent, ok := orbits[object]
			if !ok {
				break
			}
			object = parent
			distance++
		}

		// Also, we cannot easily walk down the tree, so we take the second
		// part in reverse. We walk up from the object SAN is orbiting, until
		// we hit one of the objects in the "path" map. That object will be
		// the first common parent and all that is left to do is add up the
		// current distance from the object SAN is orbiting and the distance
		// from the first part, which was stored in the "path" map.

		object, distance = orbits["SAN"], 0
		for {
			pathDistance, ok := path[object]
			if ok {
				distance += pathDistance
				break
			}
			parent, ok := orbits[object]
			if !ok {
				panic("YOU and SAN are not connected")
			}
			object = parent
			distance++
		}

		fmt.Println(distance)
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
