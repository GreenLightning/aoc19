package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	input := readFile("input.txt")
	regex := regexp.MustCompile(`^(\d{6})-(\d{6})$`)
	match := regex.FindStringSubmatch(input)
	min, max := toInt(match[1]), toInt(match[2])

	total, isolatedTotal := 0, 0
	for current := min; current <= max; current++ {
		password := []byte(fmt.Sprintf("%06d", current))

		hasMatchingDigits, hasIsolatedMatchingDigits := false, false
		for i := 0; i+1 < len(password); i++ {
			if password[i] == password[i+1] {
				hasMatchingDigits = true
				if (i-1 < 0 || password[i-1] != password[i]) && (i+2 >= len(password) || password[i+2] != password[i]) {
					hasIsolatedMatchingDigits = true
					break
				}
			}
		}

		decreasing := false
		for i := 0; i+1 < len(password); i++ {
			if password[i+1] < password[i] {
				decreasing = true
				break
			}
		}

		if hasMatchingDigits && !decreasing {
			total++
		}

		if hasIsolatedMatchingDigits && !decreasing {
			isolatedTotal++
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(total)
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(isolatedTotal)
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	check(err)
	return strings.TrimSpace(string(bytes))
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
