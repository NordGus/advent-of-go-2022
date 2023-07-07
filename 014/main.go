package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NordGus/advent-of-go-2022/014/structs"
)

const (
	inputFileName = "014/input.txt"
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	cave := structs.NewCave()
	input := scanInput(file)
	lines := parseRockFormations(input)

	for in := range lines {
		cave.AddLine(in...)
	}

	start1 := time.Now()
	fmt.Printf("How many units of sand come to rest before sand starts flowing into the abyss below? (Part 1): %v (took %s)\n", cave.HowManyUnitsOfSandBeforeOverflowing(), time.Since(start1))
	start2 := time.Now()
	fmt.Printf("How many units of sand come to rest? (Part 2): %v (took %s)\n", cave.HowManyUnitsOfSandBeforeBlockage(), time.Since(start2))
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

func parseRockFormations(input <-chan string) <-chan [][2]int {
	out := make(chan [][2]int)

	go func(input <-chan string, out chan<- [][2]int) {
		for in := range input {
			raw := strings.Split(in, " -> ")
			p := make([][2]int, len(raw))

			for i := 0; i < len(raw); i++ {
				rc := strings.Split(raw[i], ",")

				x, err := strconv.ParseInt(rc[0], 10, 0)
				if err != nil {
					panic(err)
				}

				y, err := strconv.ParseInt(rc[1], 10, 0)
				if err != nil {
					panic(err)
				}

				p[i] = [2]int{int(x), int(y)}
			}

			out <- p
		}

		close(out)
	}(input, out)

	return out
}
