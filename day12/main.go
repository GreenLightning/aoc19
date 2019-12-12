package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Moon struct {
	pos, vel Vector3
}

func main() {
	lines := readLines("input.txt")

	var input []Moon

	regex := regexp.MustCompile(`^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>$`)
	for _, line := range lines {
		match := regex.FindStringSubmatch(line)
		x, y, z := toInt(match[1]), toInt(match[2]), toInt(match[3])
		moon := Moon{
			pos: Vector3{x, y, z},
			vel: Vector3{0, 0, 0},
		}
		input = append(input, moon)
	}

	{
		fmt.Println("--- Part One ---")

		moons := make([]Moon, len(input))
		copy(moons, input)

		for step := 0; step < 1000; step++ {
			simulate(moons)
		}

		totalEnergy := 0
		for _, moon := range moons {
			potentialEnergy := moon.pos.ManhattenLength()
			kineticEnergy := moon.vel.ManhattenLength()
			totalEnergy += potentialEnergy * kineticEnergy
		}
		fmt.Println(totalEnergy)
	}

	{
		fmt.Println("--- Part Two ---")

		moons := make([]Moon, len(input))
		copy(moons, input)

		xSteps, ySteps, zSteps := 0, 0, 0
		for steps := 1; xSteps == 0 || ySteps == 0 || zSteps == 0; steps++ {
			simulate(moons)

			if xSteps == 0 {
				found := true
				for i, moon := range moons {
					if moon.pos.x != input[i].pos.x || moon.vel.x != input[i].vel.x {
						found = false
						break
					}
				}
				if found {
					xSteps = steps
				}
			}

			if ySteps == 0 {
				found := true
				for i, moon := range moons {
					if moon.pos.y != input[i].pos.y || moon.vel.y != input[i].vel.y {
						found = false
						break
					}
				}
				if found {
					ySteps = steps
				}
			}

			if zSteps == 0 {
				found := true
				for i, moon := range moons {
					if moon.pos.z != input[i].pos.z || moon.vel.z != input[i].vel.z {
						found = false
						break
					}
				}
				if found {
					zSteps = steps
				}
			}
		}

		result := lcm(xSteps, ySteps)
		result = lcm(result, zSteps)
		fmt.Println(result)
	}
}

func simulate(moons []Moon) {
	for ai, a := range moons {
		for bi, b := range moons {
			if bi == ai {
				continue
			}

			a.vel = a.vel.Plus(b.pos.Minus(a.pos).Sign())
		}
		moons[ai] = a
	}

	for index, moon := range moons {
		moons[index].pos = moon.pos.Plus(moon.vel)
	}
}

type Vector3 struct {
	x, y, z int
}

func (v Vector3) Plus(other Vector3) Vector3 {
	return Vector3{
		x: v.x + other.x,
		y: v.y + other.y,
		z: v.z + other.z,
	}
}

func (v Vector3) Minus(other Vector3) Vector3 {
	return Vector3{
		x: v.x - other.x,
		y: v.y - other.y,
		z: v.z - other.z,
	}
}

func (v Vector3) Sign() Vector3 {
	return Vector3{
		x: sign(v.x),
		y: sign(v.y),
		z: sign(v.z),
	}
}

func (v Vector3) ManhattenLength() int {
	return abs(v.x) + abs(v.y) + abs(v.z)
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

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}
