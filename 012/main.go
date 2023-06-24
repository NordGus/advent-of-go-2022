package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/NordGus/advent-of-go-2022/012/structs"
)

const (
	inputFileName = "012/input.txt"
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	input := scanInput(file)
	grid := parseGrid(input)
	heightmap := parseHeightmap(grid)
	graph := parseGraph(heightmap)

	g := <-graph
	fewestStepsBetweenStartToFinish := g.GetFewestStepsFromStartToFinish()

	fmt.Println("Fewest Steps Required:", fewestStepsBetweenStartToFinish)
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

func parseGrid(input <-chan string) <-chan [][]string {
	out := make(chan [][]string, 1)

	go func(input <-chan string, out chan<- [][]string) {
		grid := make([][]string, 0)

		for in := range input {
			line := make([]string, len(in))

			for i := 0; i < len(in); i++ {
				line[i] = string(in[i])
			}

			grid = append(grid, line)
		}

		out <- grid

		close(out)
	}(input, out)

	return out
}

func parseHeightmap(input <-chan [][]string) <-chan structs.Heightmap {
	out := make(chan structs.Heightmap, 1)

	go func(input <-chan [][]string, out chan<- structs.Heightmap) {
		grid := <-input
		heightmap := structs.NewHeightmap(len(grid), len(grid[0]))

		for i, row := range grid {
			for j := 0; j < len(row); j++ {
				nodeType := structs.NormalNode

				if row[j] == "S" {
					nodeType = structs.StartNode
					row[j] = "a"
				}

				if row[j] == "E" {
					nodeType = structs.EndNode
					row[j] = "z"
				}

				heightmap.AddNode(mapHeight(row[j]), i, j, nodeType)
			}
		}

		out <- heightmap

		close(out)
	}(input, out)

	return out
}

func parseGraph(input <-chan structs.Heightmap) <-chan structs.Heightmap {
	out := make(chan structs.Heightmap, 1)

	go func(input <-chan structs.Heightmap, out chan<- structs.Heightmap) {
		heightmap := <-input

		heightmap.BuildGraph()

		out <- heightmap

		close(out)
	}(input, out)

	return out
}

func mapHeight(value string) int {
	values := map[string]int{
		"a": 0,
		"b": 1,
		"c": 2,
		"d": 3,
		"e": 4,
		"f": 5,
		"g": 6,
		"h": 7,
		"i": 8,
		"j": 9,
		"k": 10,
		"l": 11,
		"m": 12,
		"n": 13,
		"o": 14,
		"p": 15,
		"q": 16,
		"r": 17,
		"s": 18,
		"t": 19,
		"u": 20,
		"v": 21,
		"w": 22,
		"x": 23,
		"y": 24,
		"z": 25,
	}

	return values[value]
}
