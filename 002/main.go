package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	inputFileName = "002/input.txt"

	rock     = "Rock"
	paper    = "Paper"
	scissors = "Scissors"

	victoryStatus = "victory"
	drawStatus    = "draw"
	defeatStatus  = "defeat"

	victoryScore = 6
	drawScore    = 3
	defeatScore  = 0

	rockScore     = 1
	paperScore    = 2
	scissorsScore = 3
)

type Pick struct {
	pick  string
	score int
}

type Round struct {
	opponentPick Pick
	myPick       Pick
	score        int
}

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	inputs := scanInput(file)
	rounds := buildRounds(inputs)
	scoredRound := scoreRounds(rounds)

	var totalScore int

	for round := range scoredRound {
		totalScore += round.score
	}

	fmt.Println(totalScore)
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

func buildRounds(input <-chan string) <-chan Round {
	out := make(chan Round)

	go func(input <-chan string, out chan<- Round) {
		for in := range input {
			picks := strings.Split(in, " ")
			round := Round{
				opponentPick: getPick(picks[0]),
				myPick:       getPick(picks[1]),
			}

			out <- round
		}

		close(out)
	}(input, out)

	return out
}

func scoreRounds(input <-chan Round) <-chan Round {
	out := make(chan Round)

	go func(input <-chan Round, out chan<- Round) {
		for in := range input {
			switch getRoundResult(in) {
			case victoryStatus:
				in.score += victoryScore + in.myPick.score
			case drawStatus:
				in.score += drawScore + in.myPick.score
			case defeatStatus:
				in.score += defeatScore + in.myPick.score
			}

			out <- in
		}

		close(out)
	}(input, out)

	return out
}

func getPick(pick string) Pick {
	switch pick {
	case "A", "X":
		return Pick{
			pick:  rock,
			score: rockScore,
		}
	case "B", "Y":
		return Pick{
			pick:  paper,
			score: paperScore,
		}
	case "C", "Z":
		return Pick{
			pick:  scissors,
			score: scissorsScore,
		}
	default:
		fmt.Println("unsupported type")
		return Pick{}
	}
}

func getRoundResult(round Round) string {
	if round.opponentPick.pick == round.myPick.pick {
		return drawStatus
	}

	if round.opponentPick.pick == rock && round.myPick.pick == paper {
		return victoryStatus
	}

	if round.opponentPick.pick == paper && round.myPick.pick == scissors {
		return victoryStatus
	}

	if round.opponentPick.pick == scissors && round.myPick.pick == rock {
		return victoryStatus
	}

	return defeatStatus
}
