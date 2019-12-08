package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

func main() {
	input := []byte(readFile("input.txt"))

	// Convert from character values to numbers.
	for i := range input {
		input[i] -= '0'
	}

	width, height := 25, 6
	size := width * height
	layers := len(input) / size

	{
		fmt.Println("--- Part One ---")

		bestHistogram := [3]int{math.MaxInt32, math.MaxInt32, math.MaxInt32}
		for layer := 0; layer < layers; layer++ {
			var histogram [3]int
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					value := input[layer*size+y*width+x]
					histogram[value]++
				}
			}
			if histogram[0] < bestHistogram[0] {
				bestHistogram = histogram
			}
		}

		fmt.Println(bestHistogram[1] * bestHistogram[2])
	}

	{
		fmt.Println("--- Part Two ---")

		image := make([]byte, size)
		for i := range image {
			image[i] = 2 // transparent
		}

		for layer := 0; layer < layers; layer++ {
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					value := input[layer*size+y*width+x]
					if image[y*width+x] == 2 {
						image[y*width+x] = value
					}
				}
			}
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				switch image[y*width+x] {
				case 0:
					fmt.Print(" ")
				case 1:
					fmt.Print("█")
				default:
					fmt.Print("?")
				}
			}
			fmt.Println()
		}
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	check(err)
	return strings.TrimSpace(string(bytes))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
