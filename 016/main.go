package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	inputFileName = "016/input.txt"
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	input := scanInput(file)

	for in := range input {
		fmt.Println(in)
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
