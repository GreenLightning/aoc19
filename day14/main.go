package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Part struct {
	Quantity int
	Name     string
}

type Reaction struct {
	Inputs []Part
	Output Part
}

func main() {
	lines := readLines("input.txt")

	reactions := make(map[string]Reaction)

	reactionRegex := regexp.MustCompile(`^(\d+ \w+(?:, \d+ \w+)*) => (\d+ \w+)$`)
	partRegex := regexp.MustCompile(`^(\d+) (\w+)$`)

	parsePart := func(text string) (part Part) {
		match := partRegex.FindStringSubmatch(text)
		part.Quantity = toInt(match[1])
		part.Name = match[2]
		return
	}

	for _, line := range lines {
		match := reactionRegex.FindStringSubmatch(line)
		inputs, output := match[1], match[2]
		var reaction Reaction
		for _, input := range strings.Split(inputs, ", ") {
			reaction.Inputs = append(reaction.Inputs, parsePart(input))
		}
		reaction.Output = parsePart(output)
		reactions[reaction.Output.Name] = reaction
	}

	var oreRequiredForOneFuel int

	{
		fmt.Println("--- Part One ---")
		required := map[string]int{"FUEL": 1}
		reduce(required, reactions)
		oreRequiredForOneFuel = required["ORE"]
		fmt.Println(oreRequiredForOneFuel)
	}

	{
		fmt.Println("--- Part Two ---")
		availableOre := 1_000_000_000_000
		fuel := availableOre / oreRequiredForOneFuel // we can make at least this much fuel
		required := map[string]int{"FUEL": fuel}
		for {
			required["FUEL"]++ // try to make one more fuel
			reduce(required, reactions)
			if required["ORE"] > availableOre {
				break
			}
			fuel++
		}
		fmt.Println(fuel)
	}
}

func reduce(required map[string]int, reactions map[string]Reaction) {
	for {
		changed := false

		for name, amount := range required {
			if amount > 0 {
				if reaction, ok := reactions[name]; ok {
					changed = true
					factor := (amount + reaction.Output.Quantity - 1) / reaction.Output.Quantity
					required[name] -= factor * reaction.Output.Quantity
					for _, input := range reaction.Inputs {
						required[input.Name] += factor * input.Quantity
					}
				}
			}
		}

		if !changed {
			return
		}
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
