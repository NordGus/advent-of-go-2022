package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/NordGus/advent-of-go-2022/008/structs"
)

const (
	inputFileName = "008/input.txt"
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	input := scanInput(file)
	grid := parseGrid(input)
	visibleTrees := countVisibleTrees(grid)

	var visibleTreesCount uint
	var treeCount uint

	for visible := range visibleTrees {
		if visible {
			visibleTreesCount++
		}
		treeCount++
	}

	fmt.Println("Visible Trees Count:", visibleTreesCount)
	fmt.Println("Trees Count:", treeCount)
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

func parseGrid(input <-chan string) <-chan *structs.Grid {
	out := make(chan *structs.Grid, 1)

	go func(rows <-chan string, out chan<- *structs.Grid) {
		var wg sync.WaitGroup
		var grid structs.Grid
		created := false
		currentRow := 0

		for row := range rows {
			if !created {
				size := len(row)
				grid = structs.NewGrid(size, size)
				created = true
			}

			wg.Add(1)
			go parseGridRow(&wg, &grid, currentRow, row)

			currentRow++
		}

		wg.Wait()

		out <- &grid

		close(out)
	}(input, out)

	return out
}

func parseGridRow(wg *sync.WaitGroup, grid *structs.Grid, currentRow int, row string) {
	defer wg.Done()

	for i := 0; i < len(row); i++ {
		height, err := strconv.ParseUint(string(row[i]), 10, 0)
		if err != nil {
			e := fmt.Errorf("failed to parse char '%v' for row %v on position %v. (raw: %v)", string(row[i]), currentRow, i, row)
			panic(errors.Join(e, err))
		}

		err = grid.AddTree(currentRow, i, int(height))
		if err != nil {
			e := fmt.Errorf("failed to parse char '%v' for row %v on position %v. (raw: %v)", string(row[i]), currentRow, i, row)
			panic(errors.Join(e, err))
		}
	}
}

func countVisibleTrees(input <-chan *structs.Grid) <-chan bool {
	grid := <-input
	out := make(chan bool, grid.Size())

	go func(grid *structs.Grid, out chan<- bool) {
		var wg sync.WaitGroup

		wg.Add(1)
		go grid.DetectVisibleTrees(&wg, out)

		wg.Wait()

		close(out)
	}(grid, out)

	return out
}
