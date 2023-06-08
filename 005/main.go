package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/NordGus/advent-of-go-2022/005/structs"
)

const (
	inputFileName = "005/input.txt"
)

type Instruction struct {
	Amount uint64
	From   int64
	To     int64
}

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	rawStacks, rawInstructions := scanInput(file)
	stacks := initStacks(rawStacks)
	instructions := parseInstructions(rawInstructions)
	finalStacks := moveCreates(stacks, instructions)

	for _, stack := range finalStacks {
		fmt.Print(stack.Peek())
	}

	fmt.Println("")
}

func scanInput(input *os.File) (<-chan string, <-chan string) {
	stacks := make(chan string)
	instructions := make(chan string)

	scanner := bufio.NewScanner(input)

	go func(scanner *bufio.Scanner, stacks chan<- string, instructions chan<- string) {
		active := stacks

		for scanner.Scan() {
			if scanner.Text() == "" {
				close(active)
				active = instructions
				continue
			}

			active <- scanner.Text()
		}

		close(active)
	}(scanner, stacks, instructions)

	return stacks, instructions
}

func initStacks(input <-chan string) <-chan []structs.Stack {
	out := make(chan []structs.Stack)

	go func(input <-chan string, out chan<- []structs.Stack) {
		created := false
		var tmp []structs.Stack
		for in := range input {
			if !strings.ContainsAny(in, "[]") {
				continue
			}

			parsed := parseStacksInput(in)

			if !created {
				tmp = make([]structs.Stack, len(parsed))
				created = true
			}

			for i, crate := range parsed {
				if !strings.ContainsAny(crate, "[]") {
					continue
				}
				tmp[i].Push(crate[1:2])
			}
		}

		initial := make([]structs.Stack, len(tmp))

		for i, stack := range tmp {
			for stack.Peek() != "" {
				value, err := stack.Pop()
				if err != nil {
					panic(err)
				}
				initial[i].Push(value)
			}
		}

		out <- initial

		close(out)
	}(input, out)

	return out
}

func parseStacksInput(raw string) []string {
	out := []string{}

	for i := 0; i < len(raw); i += 4 {
		out = append(out, raw[i:i+3])
	}

	return out
}

func parseInstructions(input <-chan string) <-chan Instruction {
	out := make(chan Instruction)

	go func(input <-chan string, out chan<- Instruction) {
		for in := range input {
			parsed := strings.Split(in, " ")

			amount, err := strconv.ParseUint(parsed[1], 10, 0)
			if err != nil {
				panic(err)
			}

			from, err := strconv.ParseInt(parsed[3], 10, 0)
			if err != nil {
				panic(err)
			}

			to, err := strconv.ParseInt(parsed[5], 10, 0)
			if err != nil {
				panic(err)
			}

			out <- Instruction{
				Amount: amount,
				From:   from,
				To:     to,
			}
		}

		close(out)
	}(input, out)

	return out
}

func moveCreates(stacks <-chan []structs.Stack, instructions <-chan Instruction) []structs.Stack {
	out := <-stacks

	for instruction := range instructions {
		for i := 0; i < int(instruction.Amount); i++ {
			crate, err := out[instruction.From-1].Pop()
			if err != nil {
				panic(err)
			}
			out[instruction.To-1].Push(crate)
		}
	}

	return out
}
