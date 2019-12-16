package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input := readFile("input.txt")

	var values []byte
	for _, char := range input {
		values = append(values, byte(char-'0'))
	}

	var longvalues []byte
	for i := 0; i < 10_000; i++ {
		for _, value := range values {
			longvalues = append(longvalues, value)
		}
	}

	offset := toInt(input[:7])

	{
		fmt.Println("--- Part One ---")

		for phase := 0; phase < 100; phase++ {
			result := make([]byte, len(values))
			for ri := 0; ri < len(result); ri++ {
				acc := 0
				length := ri + 1
				for index := ri; index < len(values); {
					for i := 0; i < length && index < len(values); i++ {
						acc += int(values[index])
						index++
					}
					index += length
					for i := 0; i < length && index < len(values); i++ {
						acc -= int(values[index])
						index++
					}
					index += length
				}
				result[ri] = byte(abs(acc % 10))
			}
			values = result
		}

		for i := 0; i < 8; i++ {
			fmt.Print(values[i])
		}
		fmt.Println()
	}

	{
		fmt.Println("--- Part Two ---")

		// This assumes that offset >= len(longvalues)/2, in which case each
		// output element is computed from the sum of the input elements with
		// an equal or greater index. This can be efficiently computed by
		// iterating backwards and using an inclusive scan.

		longvalues := longvalues[offset:]
		for phase := 0; phase < 100; phase++ {
			acc := 0
			for index := len(longvalues) - 1; index >= 0; index-- {
				acc += int(longvalues[index])
				longvalues[index] = byte(abs(acc % 10))
			}
		}

		for i := 0; i < 8; i++ {
			fmt.Print(longvalues[i])
		}
		fmt.Println()
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
