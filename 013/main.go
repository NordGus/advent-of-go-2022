package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	inputFileName = "013/input.txt"
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	input := scanInput(file)
	pairs := parsePairs(input)
	results := comparePackets(pairs)

	indexSum := 0

	for res := range results {
		if res.inOrder {
			indexSum += res.index
		}
	}

	fmt.Println("Sum of Indexed in Right Order:", indexSum)
}

func scanInput(input *os.File) <-chan string {
	out := make(chan string)

	scanner := bufio.NewScanner(input)

	go func(scanner *bufio.Scanner, out chan<- string) {
		for scanner.Scan() {
			out <- scanner.Text()
		}

		close(out)
	}(scanner, out)

	return out
}

func parsePairs(input <-chan string) <-chan [2]any {
	out := make(chan [2]any)

	go func(input <-chan string, out chan<- [2]any) {
		pair := [2]any{}
		idx := 0
		for in := range input {
			if in == "" {
				out <- pair
				pair = [2]any{}
				idx = 0
				continue
			}

			// let's use some weird untyped json voodoo magic shenanigans
			err := json.Unmarshal([]byte(in), &pair[idx])
			if err != nil {
				panic(err)
			}

			idx++
		}

		out <- pair

		close(out)
	}(input, out)

	return out
}

type result struct {
	index   int
	inOrder bool
}

func comparePackets(input <-chan [2]any) <-chan result {
	out := make(chan result)

	go func(input <-chan [2]any, out chan<- result) {
		idx := 1
		for pair := range input {
			// check if each packet in the pair is in order
			out <- result{index: idx, inOrder: inOrder(pair[0], pair[1])}
			idx++
		}

		close(out)
	}(input, out)

	return out
}

func inOrder(left any, right any) bool {
	return compare(left, right) <= 0
}

func compare(left any, right any) int {
	l, lok := left.([]any)
	r, rok := right.([]any)

	switch {
	case !lok && !rok: // if there's no list in the input, return the difference
		return int(left.(float64) - right.(float64))
	case !lok: // if only the left side is not a list make it a list
		l = []any{left}
	case !rok: // if only the right side is not a list make it a list
		r = []any{right}
	}

	for i := 0; i < len(l) && i < len(r); i++ {
		// recursion voodoo magic
		res := compare(l[i], r[i])

		if res != 0 {
			return res
		}
	}

	// if there's no elements to loop through return the difference in sizes
	return len(l) - len(r)
}
