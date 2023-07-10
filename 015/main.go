package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NordGus/advent-of-go-2022/015/structs"
)

const (
	inputFileName    = "015/input.txt"
	part1ComparisonY = 2_000_000
	part2LowerLimit  = 0
	part2UpperLimit  = 4_000_000
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	grid := structs.NewGrid()

	input := scanInput(file)
	raw := extractSensorAndBeaconCoordinates(input)
	parsed := parseCoordinates(raw)

	for coordinates := range parsed {
		sensor := coordinates[0:2]
		beacon := coordinates[2:]

		grid.AddSensor(sensor[0], sensor[1], beacon[0], beacon[1])
	}

	start := time.Now()
	part1, err := grid.HowManyPositionsCannotContainABeaconAt(part1ComparisonY)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: In the row where y=%v, how many positions cannot contain a beacon?: %v (took %s)\n", part1ComparisonY, part1, time.Since(start))

	start = time.Now()
	part2, err := grid.TuningFrequencyOfOfDistressBeacon(part2LowerLimit, part2UpperLimit)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Part 2: Find the only possible position for the distress beacon for x and y between %v and %v. What is its tuning frequency?: %v (took %s)\n",
		part2LowerLimit,
		part2UpperLimit,
		part2,
		time.Since(start),
	)
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

func extractSensorAndBeaconCoordinates(input <-chan string) <-chan [4]string {
	out := make(chan [4]string)

	go func(input <-chan string, out chan<- [4]string) {
		for in := range input {
			data := [4]string{}

			stage1 := strings.Split(in, ": ")
			stage2 := strings.Split(strings.ReplaceAll(stage1[0], "Sensor at ", ""), ", ")
			stage3 := strings.Split(strings.ReplaceAll(stage1[1], "closest beacon is at ", ""), ", ")

			data[0] = strings.ReplaceAll(stage2[0], "x=", "")
			data[1] = strings.ReplaceAll(stage2[1], "y=", "")
			data[2] = strings.ReplaceAll(stage3[0], "x=", "")
			data[3] = strings.ReplaceAll(stage3[1], "y=", "")

			out <- data
		}

		close(out)
	}(input, out)

	return out
}

func parseCoordinates(input <-chan [4]string) <-chan [4]int {
	out := make(chan [4]int)

	go func(input <-chan [4]string, out chan<- [4]int) {
		for in := range input {
			data := [4]int{}

			for i := 0; i < len(in); i++ {
				n, err := strconv.ParseInt(in[i], 10, 0)
				if err != nil {
					panic(err)
				}

				data[i] = int(n)
			}

			out <- data
		}

		close(out)
	}(input, out)

	return out
}
