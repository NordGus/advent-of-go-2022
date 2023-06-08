package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	inputFileName = "004/input.txt"
)

type Assignment struct {
	Start uint64
	End   uint64
}

type Pair struct {
	Assignments [2]Assignment
	Overlap     bool
}

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	inputs := scanInput(file)
	assignments := buildAssignments(inputs)
	pairs := buildPairs(assignments)
	overlapped := checkOverlap(pairs)

	overlappedPairs := 0

	for pair := range overlapped {
		if pair.Overlap {
			overlappedPairs++
		}
	}

	fmt.Println(overlappedPairs)
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

func buildAssignments(input <-chan string) <-chan [2]Assignment {
	out := make(chan [2]Assignment)

	go func(input <-chan string, out chan<- [2]Assignment) {
		for in := range input {
			pair := strings.Split(in, ",")

			out <- [2]Assignment{parseAssignment(pair[0]), parseAssignment(pair[1])}
		}

		close(out)
	}(input, out)

	return out
}

func parseAssignment(raw string) Assignment {
	pair := strings.Split(raw, "-")

	start, err := strconv.ParseUint(pair[0], 10, 0)
	if err != nil {
		panic(err)
	}

	end, err := strconv.ParseUint(pair[1], 10, 0)
	if err != nil {
		panic(err)
	}

	return Assignment{
		Start: start,
		End:   end,
	}
}

func buildPairs(input <-chan [2]Assignment) <-chan Pair {
	out := make(chan Pair)

	go func(input <-chan [2]Assignment, out chan<- Pair) {
		for in := range input {
			out <- Pair{
				Assignments: in,
			}
		}

		close(out)
	}(input, out)

	return out
}

func checkOverlap(input <-chan Pair) <-chan Pair {
	out := make(chan Pair)

	go func(input <-chan Pair, out chan<- Pair) {
		for in := range input {
			in.Overlap = in.overlaps()

			out <- in
		}

		close(out)
	}(input, out)

	return out
}

func (p *Pair) overlaps() bool {
	if p.Assignments[0].Start <= p.Assignments[1].Start && p.Assignments[0].End >= p.Assignments[1].End {
		return true
	}

	if p.Assignments[0].Start >= p.Assignments[1].Start && p.Assignments[0].End <= p.Assignments[1].End {
		return true
	}

	return false
}
