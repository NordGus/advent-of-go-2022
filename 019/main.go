package main

import (
	"bufio"
	"fmt"
	"github.com/NordGus/advent-of-go-2022/019/part1"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	inputFileName = "019/input.txt"

	part1TimeLimit = 24
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

	factories1 := make([]part1.Factory, 0, 10)

	input := scanInput(file)
	blueprints := initBlueprint(input)
	blueprintsWithIDs := parseBlueprintID(blueprints)
	blueprintsWithRobots := parseBlueprintRobots(blueprintsWithIDs)
	blueprintsCompleted := parseBlueprintRobotsCosts(blueprintsWithRobots)

	for b := range blueprintsCompleted {
		blueprint := part1.NewBlueprint(b.ID)

		for _, robot := range b.Robots {
			err := blueprint.AddRobotRecipe(robot.Type, robot.Materials)
			if err != nil {
				panic(err)
			}
		}

		factories1 = append(factories1, part1.NewFactory(blueprint))
	}

	start1 := time.Now()
	p1 := executePart1Simulation(factories1, part1TimeLimit)
	fmt.Printf("Part 1: What do you get if you add up the quality level of all of the blueprints in your list? %v (took %v)\n", p1, time.Since(start1))

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

func executePart1Simulation(factories []part1.Factory, duration int) int {
	var (
		result int
		wg     sync.WaitGroup

		results = make(chan int, len(factories))
		sem     = make(chan bool, runtime.GOMAXPROCS(0)-2)
	)

	wg.Add(len(factories))

	for _, factory := range factories {
		go func(wg *sync.WaitGroup, sem chan bool, out chan<- int, factory part1.Factory, duration int) {
			defer func() {
				<-sem
				wg.Done()
			}()

			sem <- true
			out <- factory.QualityScoreDuring(duration)
		}(&wg, sem, results, factory, duration)
	}

	wg.Wait()
	close(results)

	for r := range results {
		result += r
	}

	return result
}
