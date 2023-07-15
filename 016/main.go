package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NordGus/advent-of-go-2022/016/structs"
)

const (
	inputFileName = "016/input.txt"
)

type parsedData struct {
	name      string
	rate      int64
	neighbors []string
}

const (
	volcanoTimer int64 = 30
)

func main() {
	start := time.Now()
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	volcano := structs.NewVolcano()

	input := scanInput(file)
	data := parseInput(input)

	for in := range data {
		volcano.ParseValve(in.name, in.rate, in.neighbors)
	}

	volcano1 := volcano.Simplify()

	start1 := time.Now()
	part1 := volcano1.ReleaseTheMostPressureWithin(volcanoTimer)
	fmt.Printf("Part 1: What is the most pressure you can release? %v (took %v)\n", part1, time.Since(start1))

	fmt.Printf("took in total: %v\n", time.Since(start))
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

func parseInput(input <-chan string) <-chan parsedData {
	out := make(chan parsedData)

	go func(input <-chan string, out chan<- parsedData) {
		for in := range input {
			name := strings.Split(strings.ReplaceAll(in, "Valve ", ""), " has flow ")[0]
			unparsedRate := strings.Split(strings.Split(in, " has flow rate=")[1], "; ")[0]
			neighbors := strings.Split(strings.Split(strings.ReplaceAll(in, "leads to valve", "lead to valves"), " lead to valves ")[1], ", ")

			data := parsedData{
				name:      name,
				neighbors: neighbors,
			}

			rate, err := strconv.ParseInt(unparsedRate, 10, 0)
			if err != nil {
				panic(err)
			}

			data.rate = rate

			out <- data
		}

		close(out)
	}(input, out)

	return out
}
