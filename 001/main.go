package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const inputFileName = "001/1/input.txt"

type Elf struct {
	num   int64
	foods []int64
	total int64
}

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	inputs := scanInput(file)
	elves := createElves(inputs)
	var maxElf Elf

	for elf := range elves {
		if elf.total > maxElf.total {
			maxElf = elf
		}
	}

	fmt.Println(maxElf)
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

func createElves(input <-chan string) <-chan Elf {
	out := make(chan Elf)

	go func(in <-chan string, out chan<- Elf) {
		var elf Elf

		for in := range input {
			if in == "" {
				out <- elf

				elf = Elf{
					num:   elf.num + 1,
					foods: []int64{},
					total: 0,
				}

				continue
			}

			parsedIn, err := strconv.ParseInt(in, 10, 0)
			if err != nil {
				log.Fatal(err)
			}

			elf.total += parsedIn
			elf.foods = append(elf.foods, parsedIn)
		}

		if elf.total > 0 {
			out <- elf
		}

		close(out)
	}(input, out)

	return out
}
