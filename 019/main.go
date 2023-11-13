package main

import (
	"bufio"
	"fmt"
	"github.com/NordGus/advent-of-go-2022/019/factory"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	inputFileName = "019/input.txt"

	part1TimeLimit = 24
	part2TimeLimit = 32
)

type InputBlueprint struct {
	RawData string
	ID      int
	Robots  []InputRobotCost
}

type InputRobotCost struct {
	RawData   string
	Type      string
	Materials map[string]int
}

func main() {
	start := time.Now()
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	factories := make([]factory.Factory, 0, 10)

	input := scanInput(file)
	blueprints := initBlueprint(input)
	blueprintsWithIDs := parseBlueprintID(blueprints)
	blueprintsWithRobots := parseBlueprintRobots(blueprintsWithIDs)
	blueprintsCompleted := parseBlueprintRobotsCosts(blueprintsWithRobots)

	for b := range blueprintsCompleted {
		blueprint := factory.NewBlueprint(b.ID)

		for _, robot := range b.Robots {
			err := blueprint.AddRobotRecipe(robot.Type, robot.Materials)
			if err != nil {
				panic(err)
			}
		}

		factories = append(factories, factory.NewFactory(blueprint))
	}

	start1 := time.Now()
	p1 := executePart1Simulation(factories, part1TimeLimit)
	fmt.Printf("Part 1: What do you get if you add up the quality level of all of the blueprints in your list? %v (took %v)\n", p1, time.Since(start1))

	start2 := time.Now()
	p2 := executePart2Simulation(factories, part2TimeLimit)
	fmt.Printf("Part 2: What do you get if you multiply these numbers together? %v (took %v)\n", p2, time.Since(start2))

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

func initBlueprint(input <-chan string) <-chan InputBlueprint {
	out := make(chan InputBlueprint, 5)

	go func(input <-chan string, out chan<- InputBlueprint) {
		for in := range input {
			out <- InputBlueprint{RawData: in}
		}

		close(out)
	}(input, out)

	return out
}

func parseBlueprintID(blueprints <-chan InputBlueprint) <-chan InputBlueprint {
	out := make(chan InputBlueprint, 5)

	go func(input <-chan InputBlueprint, out chan<- InputBlueprint) {
		for in := range input {
			var (
				rawData       []string
				blueprintInfo []string
				id            int64
			)

			rawData = strings.Split(in.RawData, ": ")
			blueprintInfo = strings.Split(rawData[0], " ")
			id, _ = strconv.ParseInt(blueprintInfo[1], 10, 64)

			in.RawData = rawData[1]
			in.ID = int(id)

			out <- in
		}

		close(out)
	}(blueprints, out)

	return out
}

func parseBlueprintRobots(blueprints <-chan InputBlueprint) <-chan InputBlueprint {
	out := make(chan InputBlueprint, 5)

	go func(input <-chan InputBlueprint, out chan<- InputBlueprint) {
		for in := range input {
			var rawData []string

			rawData = strings.Split(in.RawData, ". ")

			for i := 0; i < len(rawData); i++ {
				rawData[i] = strings.ReplaceAll(rawData[i], ".", "")
			}

			in.Robots = make([]InputRobotCost, len(rawData))

			for i := 0; i < len(in.Robots); i++ {
				switch {
				case strings.Contains(rawData[i], "ore robot"):
					in.Robots[i].Type = "ore"
				case strings.Contains(rawData[i], "clay robot"):
					in.Robots[i].Type = "clay"
				case strings.Contains(rawData[i], "obsidian robot"):
					in.Robots[i].Type = "obsidian"
				default:
					in.Robots[i].Type = "geode"
				}

				in.Robots[i].RawData = rawData[i]
			}

			in.RawData = ""

			out <- in
		}

		close(out)
	}(blueprints, out)

	return out
}

func parseBlueprintRobotsCosts(blueprints <-chan InputBlueprint) <-chan InputBlueprint {
	out := make(chan InputBlueprint, 5)

	go func(input <-chan InputBlueprint, out chan<- InputBlueprint) {
		for in := range input {
			for i, robot := range in.Robots {
				var rawData []string

				robot.RawData = strings.ReplaceAll(
					robot.RawData,
					fmt.Sprintf("Each %s robot costs ", robot.Type),
					"",
				)

				rawData = strings.Split(robot.RawData, " and ")

				in.Robots[i].Materials = make(map[string]int, len(rawData))

				for j := 0; j < len(rawData); j++ {
					data := strings.Split(rawData[j], " ")

					cost, _ := strconv.ParseInt(data[0], 10, 64)

					in.Robots[i].Materials[data[1]] = int(cost)
				}

				in.Robots[i].RawData = ""
			}

			out <- in
		}

		close(out)
	}(blueprints, out)

	return out
}

func executePart1Simulation(factories []factory.Factory, duration int) int {
	var result int

	runtime.GC()

	for i := 0; i < len(factories); i++ {
		result += factories[i].QualityScoreDuring(duration) * factories[i].BlueprintID()
	}

	return result
}

func executePart2Simulation(factories []factory.Factory, duration int) int {
	var result int = 1

	runtime.GC()

	for i := 0; i < 3; i++ {
		result *= factories[i].QualityScoreDuring(duration)
		runtime.GC()
	}

	return result
}
