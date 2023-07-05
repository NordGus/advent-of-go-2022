package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
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
	packets := parsePackets(input)

	packets = append(packets, []any{[]any{2.}}, []any{[]any{6.}})

	sort.Slice(packets, func(i int, j int) bool {
		return compare(packets[i], packets[j]) < 0
	})

	key := 1

	for i, packet := range packets {
		p := fmt.Sprint(packet)

		if p == "[[2]]" || p == "[[6]]" {
			key *= i + 1
		}
	}

	fmt.Println("The decoder key:", key)
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

func parsePackets(input <-chan string) []any {
	out := make([]any, 0, 1_000)

	for in := range input {
		if in == "" {
			continue
		}

		var entry any

		// let's use some weird untyped json voodoo magic shenanigans
		err := json.Unmarshal([]byte(in), &entry)
		if err != nil {
			panic(err)
		}

		out = append(out, entry)
	}

	return out
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
