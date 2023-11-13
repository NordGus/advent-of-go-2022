package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	inputFileName = "019/input.txt"
)

func main() {
	start := time.Now()
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	input := scanInput(file)

	for in := range input {
		fmt.Println(in)
	}

	fmt.Printf("took in total: %v\n", time.Since(start))
}

func scanInput(input *os.File) <-chan string {
	out := make(chan string, 5)

	scanner := bufio.NewScanner(input)

	go func(scanner *bufio.Scanner, out chan<- string) {
		for scanner.Scan() {
			out <- scanner.Text()
		}

		input.Close()
		close(out)
	}(scanner, out)

	return out
}
