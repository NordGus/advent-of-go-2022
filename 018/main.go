package main

import (
	"bufio"
	"fmt"
	"github.com/NordGus/advent-of-go-2022/018/part2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NordGus/advent-of-go-2022/018/part1"
)

const (
	inputFileName = "018/input.txt"
)

func main() {
	start := time.Now()
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	c1 := part1.NewCloud()
	c2 := part2.NewCloud()

	input := scanInput(file)
	points := parsePoints(input)

	for in := range points {
		c1.AddPoint(in)
		c2.AddLavaPoint(in)
	}

	start1 := time.Now()
	p1 := c1.CountSurfaceAreaOfLavaDroplet()
	fmt.Printf("Part 1: What is the surface area of your scanned lava droplet? %v (took %v)\n", p1, time.Since(start1))

	start2 := time.Now()
	p2 := c2.CountExternalSurfaceAreaOfLavaDroplet()
	fmt.Printf("Part 1: What is the surface area of your scanned lava droplet? %v (took %v)\n", p2, time.Since(start2))

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

func parsePoints(points <-chan string) <-chan [3]int {
	out := make(chan [3]int)

	go func(points <-chan string, out chan<- [3]int) {
		for point := range points {
			pout := [3]int{}

			for i, v := range strings.Split(point, ",") {
				n, _ := strconv.Atoi(v)

				pout[i] = n
			}

			out <- pout
		}

		close(out)
	}(points, out)

	return out
}
