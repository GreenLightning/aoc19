package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	modules := readNumbers("input.txt")

	total, totalRecursive := 0, 0
	for _, mass := range modules {
		total += mass/3 - 2

		for {
			fuel := mass/3 - 2
			if fuel <= 0 {
				break
			}
			totalRecursive += fuel
			mass = fuel
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(total)
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(totalRecursive)
	}
}

func readNumbers(filename string) []int {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var numbers []int
	for scanner.Scan() {
		numbers = append(numbers, toInt(scanner.Text()))
	}
	return numbers
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
