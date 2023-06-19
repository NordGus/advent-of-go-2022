package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/NordGus/advent-of-go-2022/011/structs"
)

const (
	inputFileName = "011/input.txt"

	numberOfRounds = 10_000
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	input := scanInput(file)
	monkeys := buildMonkeys(input)
	list := prepareMonkeys(monkeys)
	final := play(list)

	structs.SortByItemInspectedCount(final)
	monkeyBusiness := 1

	for i := 0; i < 2; i++ {
		monkeyBusiness *= int(final[i].InspectedItemsCount)
	}

	fmt.Println("Level of monkey business:", monkeyBusiness)
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

func buildMonkeys(input <-chan string) <-chan structs.Monkey {
	out := make(chan structs.Monkey)

	go func(input <-chan string, out chan<- structs.Monkey) {
		monkey := structs.NewMonkey()

		for in := range input {
			str := strings.TrimSpace(in)

			if str == "" {
				out <- monkey
				monkey = structs.NewMonkey()
				continue
			}

			if strings.Contains(str, "Starting items: ") {
				monkey.AddItems(parseInitialItems(str)...)
				continue
			}

			if strings.Contains(str, "Operation: ") {
				err := monkey.SetOperation(parseOperation(str))
				if err != nil {
					panic(err)
				}
				continue
			}

			if strings.Contains(str, "Test: ") {
				monkey.SetTest(parseTest(str))
				continue
			}

			if strings.Contains(str, "If true: throw to monkey ") {
				monkey.SetTestSuccessTarget(parseTestTarget(str))
				continue
			}

			if strings.Contains(str, "If false: throw to monkey ") {
				monkey.SetTestFailTarget(parseTestTarget(str))
				continue
			}
		}

		out <- monkey
		close(out)
	}(input, out)

	return out
}

func parseInitialItems(input string) []int {
	items := strings.Split(input, " ")
	items = items[2:]
	out := make([]int, len(items))

	for i := 0; i < len(items); i++ {
		item, err := strconv.ParseInt(strings.Trim(items[i], ","), 10, 0)
		if err != nil {
			panic(err)
		}

		out[i] = int(item)
	}

	return out
}

func parseOperation(input string) func(int) int {
	operation := strings.Split(input, " ")
	operation = operation[3:]

	if operation[0] == operation[2] {
		return parseOldOldOperation(operation[1])
	}

	return parseOldConstOperation(operation[1], operation[2])
}

func parseOldOldOperation(operation string) func(int) int {
	switch operation {
	case "+":
		return func(old int) int {
			return old + old
		}
	case "-":
		return func(old int) int {
			return 0
		}
	case "*":
		return func(old int) int {
			return old * old
		}
	case "/":
		return func(old int) int {
			return 1
		}
	default:
		panic("unsupported operation type")
	}
}

func parseOldConstOperation(operation string, constant string) func(int) int {
	con, err := strconv.ParseInt(constant, 10, 0)
	if err != nil {
		panic(err)
	}

	switch operation {
	case "+":
		return func(old int) int {
			return old + int(con)
		}
	case "-":
		return func(old int) int {
			return old - int(con)
		}
	case "*":
		return func(old int) int {
			return old * int(con)
		}
	case "/":
		return func(old int) int {
			return old / int(con)
		}
	default:
		panic("unsupported operation type")
	}
}

func parseTest(input string) int {
	str := strings.Split(input, " ")
	constant, err := strconv.ParseInt(str[len(str)-1], 10, 0)
	if err != nil {
		panic(err)
	}

	return int(constant)
}

func parseTestTarget(input string) int {
	str := strings.Split(input, " ")
	target, err := strconv.ParseInt(str[len(str)-1], 10, 0)
	if err != nil {
		panic(err)
	}

	return int(target)
}

func prepareMonkeys(input <-chan structs.Monkey) []structs.Monkey {
	out := make([]structs.Monkey, 0)

	for monkey := range input {
		out = append(out, monkey)
	}

	return out
}

func play(monkeys []structs.Monkey) []structs.Monkey {
	modulo := structs.GetLimiter(monkeys)

	for round := 0; round < numberOfRounds; round++ {
		structs.PlayRound(monkeys, modulo)
	}

	return monkeys
}
