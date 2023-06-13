package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/NordGus/advent-of-go-2022/010/structs"
)

const (
	inputFileName = "010/input.txt"
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	input := scanInput(file)
	chips := initChip(input)
	probes := chips[0].RunProgram(firstProbes())
	//render := chips[1].RenderScreen()

	part1 := make(chan int, 1)

	// Part 1 solution
	go func(input <-chan structs.ProbeResult, out chan<- int) {
		sum := 0

		for probe := range probes {
			sum += int(probe.Cycle) * probe.X
		}

		out <- sum
		close(out)
	}(probes, part1)

	fmt.Println("Sum Of Signal Strengths:", <-part1)
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

func initChip(input <-chan string) [2]structs.Chip {
	chips := [2]structs.Chip{structs.NewChip(), structs.NewChip()}

	for in := range input {
		chips[0].ParseAndQueueInstruction(in)
		chips[1].ParseAndQueueInstruction(in)
	}

	return chips
}

func firstProbes() func(uint) bool {
	probes := map[uint]bool{
		20:  true,
		60:  true,
		100: true,
		140: true,
		180: true,
		220: true,
	}

	return func(key uint) bool {
		return probes[key]
	}
}
