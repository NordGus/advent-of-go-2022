package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/NordGus/advent-of-go-2022/009/structs"
)

const (
	inputFileName = "009/input.txt"
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	input := scanInput(file)
	moves := parseMovements(input)
	rope := makeMoves(moves)

	r := <-rope

	fmt.Println(r.CountTailUniqueLocations())
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

func parseMovements(input <-chan string) <-chan structs.Movement {
	out := make(chan structs.Movement)

	go func(input <-chan string, out chan<- structs.Movement) {
		for in := range input {
			move := strings.Split(in, " ")
			amount, err := strconv.ParseInt(move[1], 10, 0)
			if err != nil {
				log.Fatalf("error Parsing Movements: %v", err.Error())
			}

			movement := structs.Movement{
				Amount: int(amount),
			}

			switch move[0] {
			case "R":
				movement.Direction = structs.RightMovement
			case "L":
				movement.Direction = structs.LeftMovement
			case "U":
				movement.Direction = structs.UpMovement
			case "D":
				movement.Direction = structs.DownMovement
			default:
				log.Fatalf("error Parsing Movements: unsupported movement direction")
			}

			out <- movement
		}

		close(out)
	}(input, out)

	return out
}

func makeMoves(input <-chan structs.Movement) <-chan *structs.Rope {
	out := make(chan *structs.Rope, 1)

	go func(moves <-chan structs.Movement, out chan<- *structs.Rope) {
		rope := structs.NewRope()

		for move := range moves {
			rope.Move(move)
		}

		out <- &rope

		close(out)
	}(input, out)

	return out
}
