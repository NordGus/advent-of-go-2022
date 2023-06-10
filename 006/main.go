package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/NordGus/advent-of-go-2022/006/structs"
)

const (
	inputFileName      = "006/input.txt"
	startOfPacketSize  = 4
	startOfMessageSize = 14
	streamCount        = 5
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	signal := scanInput(file)
	startOfPacketIndex, startOfMessageIndex := detect(signal)

	for startOfPacket := range startOfPacketIndex {
		fmt.Println("First start-of-packet:", startOfPacket)
	}

	for startOfMessage := range startOfMessageIndex {
		fmt.Println("First start-of-message:", startOfMessage)
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

func detect(input <-chan string) (<-chan uint, <-chan uint) {
	packetSignal := make(chan string, streamCount)
	messageSignal := make(chan string, streamCount)

	startOfPacket := detectStartOfPacket(packetSignal)
	startOfMessage := detectStartOfMessage(messageSignal)

	go func(input <-chan string) {
		for signal := range input {
			packetSignal <- signal
			messageSignal <- signal
		}
		close(packetSignal)
		close(messageSignal)
	}(input)

	return startOfPacket, startOfMessage
}

func detectStartOfPacket(input <-chan string) <-chan uint {
	out := make(chan uint)

	go func(input <-chan string, out chan<- uint) {
		decoder := structs.NewDecoder(startOfPacketSize)

		for signal := range input {
			for _, r := range signal {
				decoder.Push(r)

				start, index := decoder.IsStart()

				if start {
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

func detectStartOfMessage(input <-chan string) <-chan uint {
	out := make(chan uint)

	go func(input <-chan string, out chan<- uint) {
		decoder := structs.NewDecoder(startOfMessageSize)

		for signal := range input {
			for _, r := range signal {
				decoder.Push(r)

				start, index := decoder.IsStart()

				if start {
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
