package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/NordGus/advent-of-go-2022/020/encrypted"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	inputFileName = "020/input.txt"
)

func main() {
	var (
		start = time.Now()
		ctx   = context.Background()
	)

	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	encryptedFile := encrypted.New()

	input := scanInput(ctx, file)
	numbers := parseNumber(ctx, input)

	for number := range numbers {
		encryptedFile.AddItem(number)
	}

	start1 := time.Now()
	p1 := encryptedFile.MixFilePart1(1000, 2000, 3000)
	fmt.Printf("Part 1: What is the sum of the three numbers that form the grove coordinates? %v (took %v)\n", p1, time.Since(start1))

	fmt.Printf("took in total: %v\n", time.Since(start))
}

func scanInput(ctx context.Context, input *os.File) <-chan string {
	out := make(chan string)

	scanner := bufio.NewScanner(input)

	go func(ctx context.Context, scanner *bufio.Scanner, out chan<- string) {
		defer input.Close()
		defer close(out)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case out <- scanner.Text():
			}
		}
	}(ctx, scanner, out)

	return out
}

func parseNumber(ctx context.Context, input <-chan string) <-chan int64 {
	out := make(chan int64)

	go func(ctx context.Context, input <-chan string, out chan<- int64) {
		defer close(out)

		for in := range input {
			number, err := strconv.ParseInt(in, 10, 64)
			if err != nil {
				panic(err)
			}

			select {
			case <-ctx.Done():
				return
			case out <- number:
			}
		}
	}(ctx, input, out)

	return out
}
