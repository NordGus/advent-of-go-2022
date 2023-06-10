package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/NordGus/advent-of-go-2022/006/structs"
)

const (
	inputFileName     = "006/input.txt"
	startOfPacketSize = 4
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	signal := scanInput(file)
	startOfPacketIndex := detectStartOfPacket(signal)

	for input := range startOfPacketIndex {
		fmt.Println("First start-of-packet:", input)
	}
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

func detectStartOfPacket(input <-chan string) <-chan uint {
	out := make(chan uint)

	go func(input <-chan string, out chan<- uint) {
		decoder := structs.NewDecoder(startOfPacketSize)

		for signal := range input {
			for _, r := range signal {
				decoder.Push(r)

				startOfPacket, index := decoder.IsStartOfPackage()

				if startOfPacket {
					out <- index
					break
				}
			}

			decoder.Clear()
		}

		close(out)
	}(input, out)

	return out
}
