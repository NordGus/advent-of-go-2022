package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NordGus/advent-of-go-2022/017/structs"
)

const (
	inputFileName   = "017/input.txt"
	part1RocksLimit = 2022
)

func main() {
	start := time.Now()
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	chamber := structs.NewChamber()

	input := scanInput(file)

	for in := range input {
		chamber.SetJets(in)
	}

	start1 := time.Now()
	part1 := chamber.HowManyUnitsTallWillTheTowerOfRocksBeAfterNRocksHaveStoppedFalling(part1RocksLimit)
	fmt.Printf("Part 1: How many units tall will the tower of rocks be after 2022 rocks have stopped falling? %v (took %v)\n", part1, time.Since(start1))

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
