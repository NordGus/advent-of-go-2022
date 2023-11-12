package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	blprnt "github.com/NordGus/advent-of-go-2022/019/shared/blueprint"
)

const (
	inputFileName = "019/input.txt"
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

	blprnts := make([]blprnt.Blueprint, 0, 10)

	input := scanInput(file)
	blueprints := initBlueprint(input)
	blueprintsWithIDs := parseBlueprintID(blueprints)
	blueprintsWithRobots := parseBlueprintRobots(blueprintsWithIDs)
	blueprintsCompleted := parseBlueprintRobotsCosts(blueprintsWithRobots)

	for blueprint := range blueprintsCompleted {
		b := blprnt.New(blueprint.ID)

		for _, robot := range blueprint.Robots {
			err := b.AddRobotRecipe(robot.Type, robot.Materials)
			if err != nil {
				panic(err)
			}
		}

		blprnts = append(blprnts, b)

		fmt.Printf("%v\n", blueprint)
		fmt.Printf("%+v\n", b)
	}

	fmt.Printf("took in total: %v\n", time.Since(start))
}

func scanInput(input *os.File) <-chan string {
	out := make(chan string)

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
	out := make(chan InputBlueprint)

	go func(input <-chan string, out chan<- InputBlueprint) {
		for in := range input {
			out <- InputBlueprint{RawData: in}
		}

		close(out)
	}(input, out)

	return out
}

func parseBlueprintID(blueprints <-chan InputBlueprint) <-chan InputBlueprint {
	out := make(chan InputBlueprint)

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
	out := make(chan InputBlueprint)

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
	out := make(chan InputBlueprint)

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
