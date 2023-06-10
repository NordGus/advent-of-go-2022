package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NordGus/advent-of-go-2022/007/structs"
)

const (
	inputFileName = "007/input.txt"
	maxDirSize    = 100_000
)

func main() {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fs := structs.NewFilesystem()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " ")

		if data[0] == "$" && data[1] == "ls" {
			continue
		}

		if data[0] == "$" && data[1] == "cd" {
			fs.CD(data[2])
			continue
		}

		fs.Stream(data[1], data[0])
	}

	var totalSize uint64

	for _, size := range fs.DirectorySizes() {
		if size <= maxDirSize {
			totalSize += size
		}
	}

	fmt.Println("Total size:", totalSize)
}
