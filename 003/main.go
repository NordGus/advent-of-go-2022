package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	inputFileName = "003/input.txt"

	rucksackCompartmentCount = 2
)

type Rucksack struct {
	Compartments []string
	CommonItems  map[string]uint
	Priority     uint
}

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	inputs := scanInput(file)
	rucksacks := buildRucksack(inputs)
	searchedRucksacks := checkForCommonItems(rucksacks)
	prioritizedRucksacks := setRucksackPriority(searchedRucksacks)

	var prioritySum uint

	for rucksack := range prioritizedRucksacks {
		prioritySum += rucksack.Priority
	}

	fmt.Println(prioritySum)
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

func buildRucksack(input <-chan string) <-chan Rucksack {
	out := make(chan Rucksack)

	go func(input <-chan string, output chan<- Rucksack) {
		for in := range input {
			inRunes := []rune(in)
			inLen := len(inRunes)
			stepSize := inLen / rucksackCompartmentCount

			if inLen%rucksackCompartmentCount != 0 {
				fmt.Println("invalid input", in)
				continue
			}

			compartments := make([]string, rucksackCompartmentCount)
			commonItems := make(map[string]uint)

			for i := 0; i < rucksackCompartmentCount; i++ {
				start := i * stepSize
				end := (i + 1) * stepSize
				compartments[i] = string(inRunes[start:end])
			}

			output <- Rucksack{
				Compartments: compartments,
				CommonItems:  commonItems,
			}
		}

		close(output)
	}(input, out)

	return out
}

func checkForCommonItems(input <-chan Rucksack) <-chan Rucksack {
	out := make(chan Rucksack)

	go func(input <-chan Rucksack, output chan<- Rucksack) {
		for in := range input {
			for i := 0; i < rucksackCompartmentCount; i++ {
				for _, item := range in.Compartments[i] {
					it := string(item)
					count, ok := in.CommonItems[it]

					if !ok && i == 0 {
						in.CommonItems[it] = 1
					}

					if ok && count <= uint(i) {
						in.CommonItems[it] += 1
					}
				}
			}

			output <- in
		}

		close(output)
	}(input, out)

	return out
}

func setRucksackPriority(input <-chan Rucksack) <-chan Rucksack {
	out := make(chan Rucksack)

	go func(input <-chan Rucksack, output chan<- Rucksack) {
		for in := range input {
			for item, count := range in.CommonItems {
				if count == rucksackCompartmentCount {
					in.Priority = prioritize(item)
				}
			}

			output <- in
		}

		close(output)
	}(input, out)

	return out
}

func prioritize(item string) uint {
	priorities := map[string]uint{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
		"f": 6,
		"g": 7,
		"h": 8,
		"i": 9,
		"j": 10,
		"k": 11,
		"l": 12,
		"m": 13,
		"n": 14,
		"o": 15,
		"p": 16,
		"q": 17,
		"r": 18,
		"s": 19,
		"t": 20,
		"u": 21,
		"v": 22,
		"w": 23,
		"x": 24,
		"y": 25,
		"z": 26,
		"A": 27,
		"B": 28,
		"C": 29,
		"D": 30,
		"E": 31,
		"F": 32,
		"G": 33,
		"H": 34,
		"I": 35,
		"J": 36,
		"K": 37,
		"L": 38,
		"M": 39,
		"N": 40,
		"O": 41,
		"P": 42,
		"Q": 43,
		"R": 44,
		"S": 45,
		"T": 46,
		"U": 47,
		"V": 48,
		"W": 49,
		"X": 50,
		"Y": 51,
		"Z": 52,
	}

	return priorities[item]
}
