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

	opponentPickRock    = "A"
	opponentPickPaper   = "B"
	opponentPickScissor = "C"

	outcomeDefeat  = "X"
	outcomeDraw    = "Y"
	outcomeVictory = "Z"
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
			round := Round{}
			round.opponentPick = getOpponentPick(picks[0])
			round.myPick = getMyPick(round.opponentPick, picks[1])

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
			switch getRoundOutcome(in) {
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

func getOpponentPick(pick string) Pick {
	switch pick {
	case opponentPickRock:
		return Pick{
			pick:  rock,
			score: rockScore,
		}
	case opponentPickPaper:
		return Pick{
			pick:  paper,
			score: paperScore,
		}
	case opponentPickScissor:
		return Pick{
			pick:  scissors,
			score: scissorsScore,
		}
	default:
		log.Println("unsupported type")
		return Pick{}
	}
}

func getMyPick(opponentPick Pick, outcome string) Pick {
	if outcome == outcomeDefeat {
		return getLosingPick(opponentPick)
	}

	if outcome == outcomeVictory {
		return getWinningPick(opponentPick)
	}

	return Pick{
		pick:  opponentPick.pick,
		score: opponentPick.score,
	}
}

func getLosingPick(opponentPick Pick) Pick {
	switch opponentPick.pick {
	case rock:
		return Pick{
			pick:  scissors,
			score: scissorsScore,
		}
	case paper:
		return Pick{
			pick:  rock,
			score: rockScore,
		}
	case scissors:
		return Pick{
			pick:  paper,
			score: paperScore,
		}
	default:
		log.Println("unsupported pick")
		return Pick{}
	}
}

func getWinningPick(opponentPick Pick) Pick {
	switch opponentPick.pick {
	case rock:
		return Pick{
			pick:  paper,
			score: paperScore,
		}
	case paper:
		return Pick{
			pick:  scissors,
			score: scissorsScore,
		}
	case scissors:
		return Pick{
			pick:  rock,
			score: rockScore,
		}
	default:
		log.Println("unsupported pick")
		return Pick{}
	}
}

func getRoundOutcome(round Round) string {
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
