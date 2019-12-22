package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

const (
	KindDealStack     = 0
	KindDealIncrement = 1
	KindCut           = 2
)

type Shuffle struct {
	Kind  int
	Value int64
}

func main() {
	lines := readLines("input.txt")

	dealStackRegex := regexp.MustCompile(`^deal into new stack$`)
	dealIncrementRegex := regexp.MustCompile(`^deal with increment (\d+)$`)
	cutRegex := regexp.MustCompile(`^cut (-?\d+)$`)

	var input []Shuffle
	for _, line := range lines {
		if match := dealStackRegex.FindStringSubmatch(line); match != nil {
			input = append(input, Shuffle{KindDealStack, 0})

		} else if match := dealIncrementRegex.FindStringSubmatch(line); match != nil {
			increment := toInt64(match[1])
			input = append(input, Shuffle{KindDealIncrement, increment})

		} else if match := cutRegex.FindStringSubmatch(line); match != nil {
			cut := toInt64(match[1])
			input = append(input, Shuffle{KindCut, cut})

		} else {
			panic(line)
		}
	}

	{
		fmt.Println("--- Part One ---")

		const count = 10007

		cards := make([]int, count)
		for index := range cards {
			cards[index] = index
		}

		tmp := make([]int, count)

		shuffles := compact(input, count)

		for _, shuffle := range shuffles {
			switch shuffle.Kind {
			case KindDealStack:
				for index, card := range cards {
					tmp[len(cards)-1-index] = card
				}

			case KindDealIncrement:
				var index int64
				increment := shuffle.Value
				for _, card := range cards {
					tmp[index] = card
					index = (index + increment) % count
				}

			case KindCut:
				cut := shuffle.Value
				copy(tmp, cards[cut:])
				copy(tmp[count-cut:], cards)
			}
			cards, tmp = tmp, cards
		}

		for index, card := range cards {
			if card == 2019 {
				fmt.Println(index)
				break
			}
		}
	}

	{
		fmt.Println("--- Part Two ---")

		const count = 119315717514047
		const iterations = 101741582076661

		var shuffles []Shuffle
		factor := compact(input, count)

		// Exponentiation by squaring.
		for iterationsLeft := count - iterations - 1; iterationsLeft != 0; iterationsLeft /= 2 {
			if iterationsLeft%2 == 1 {
				shuffles = append(shuffles, factor...)
				shuffles = compact(shuffles, count)
			}
			factor = append(factor, factor...)
			factor = compact(factor, count)
		}

		pos := big.NewInt(2020)
		for _, shuffle := range shuffles {
			if shuffle.Kind == KindDealIncrement {
				increment := shuffle.Value
				pos.Mul(pos, big.NewInt(increment))
				pos.Mod(pos, big.NewInt(count))

			} else if shuffle.Kind == KindDealStack {
				pos.Sub(big.NewInt(count-1), pos)

			} else if shuffle.Kind == KindCut {
				cut := shuffle.Value

				if pos.Int64() < cut {
					pos.Add(pos, big.NewInt(count-cut))
				} else {
					pos.Sub(pos, big.NewInt(cut))
				}
			}
		}

		fmt.Println(pos.Int64())
	}

}

func compact(input []Shuffle, count int64) []Shuffle {
	// Compact "deal into stack" shuffles.
	//
	// Two consecutive "deal into stack" shuffles cancel each other. So we
	// iterate over the input list, tracking whether we need to currently need
	// reverse the stack, which changes every time we see a "deal into stack"
	// shuffle. Then, if we need to reverse at the end, we add a single "deal
	// into stack" shuffle to the output.
	//
	// If we currently need to reverse the stack, we have to modify the other
	// shuffles. This boils down to the following two rules, where the list of
	// instructions below the line has the same effect as the list of
	// instructions above the line:
	//
	// deal into new stack
	// cut x
	// -------------------
	// cut count-x
	// deal into new stack
	//
	// deal into new stack
	// deal with increment x
	// ---
	// deal with increment x
	// cut count+1-x
	// deal into new stack
	//
	{
		compacted := make([]Shuffle, 0, len(input))
		reverse := false
		for _, shuffle := range input {
			if shuffle.Kind == KindDealStack {
				reverse = !reverse
				continue
			}
			if !reverse {
				compacted = append(compacted, shuffle)
				continue
			}
			switch shuffle.Kind {
			case KindDealIncrement:
				compacted = append(compacted, shuffle)
				compacted = append(compacted, Shuffle{KindCut, count + 1 - shuffle.Value})

			case KindCut:
				cut := (shuffle.Value + count) % count // normalize negative values
				cut = count - cut                      // reverse cut
				compacted = append(compacted, Shuffle{KindCut, cut})
			}
		}
		if reverse {
			compacted = append(compacted, Shuffle{KindDealStack, 0})
		}
		input = compacted
	}

	// Compact "cut" shuffles.
	//
	// Here we require that the "deal into stack" shuffles have been compacted
	// already, so we can insert the "cut" shuffle before the "deal into
	// stack" shuffle or at the end. Then, we only have to handle "deal with
	// increment" shuffles.
	//
	// cut x
	// cut y
	// ---
	// cut (x+y) % count
	//
	// cut x
	// deal with increment y
	// ---
	// deal with increment y
	// cut (x*y) % count
	//
	{
		compacted := make([]Shuffle, 0, len(input))
		cut := big.NewInt(0)
		for _, shuffle := range input {
			switch shuffle.Kind {
			case KindDealStack:
				if value := cut.Int64(); value != 0 {
					compacted = append(compacted, Shuffle{KindCut, value})
					cut.SetInt64(0)
				}
				compacted = append(compacted, shuffle)

			case KindDealIncrement:
				compacted = append(compacted, shuffle)
				cut.Mul(cut, big.NewInt(shuffle.Value))
				cut.Mod(cut, big.NewInt(count))

			case KindCut:
				cut.Add(cut, big.NewInt(shuffle.Value))
				cut.Mod(cut, big.NewInt(count))
			}
		}
		if value := cut.Int64(); value != 0 {
			compacted = append(compacted, Shuffle{KindCut, value})
			cut.SetInt64(0)
		}
		input = compacted
	}

	// Compact "deal with increment" shuffles.
	//
	// Finally, we just have to combine "deal with increment" shuffles.
	//
	// deal with increment x
	// deal with increment y
	// ---
	// deal with increment (x*y) % count
	//
	{
		compacted := make([]Shuffle, 0, len(input))
		increment := big.NewInt(1)
		for _, shuffle := range input {
			switch shuffle.Kind {
			case KindDealIncrement:
				increment.Mul(increment, big.NewInt(shuffle.Value))
				increment.Mod(increment, big.NewInt(count))

			default:
				if value := increment.Int64(); value != 1 {
					compacted = append(compacted, Shuffle{KindDealIncrement, value})
					increment.SetInt64(1)
				}
				compacted = append(compacted, shuffle)
			}
		}
		if value := increment.Int64(); value != 1 {
			compacted = append(compacted, Shuffle{KindDealIncrement, value})
			increment.SetInt64(1)
		}
		input = compacted
	}

	return input
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

func toInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	check(err)
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
