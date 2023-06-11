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
	inputFileName        = "007/input.txt"
	maxDirSize    uint64 = 100_000
	maxDiskSpace  uint64 = 70_000_000
	unusedTarget  uint64 = 30_000_000
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

	results := fs.DirectorySizes() // results are sorted by size from smallest to biggest

	currentUnusedSpace := maxDiskSpace - results[len(results)-1].Size // Filesystem Root Dir Size (AKA used space)
	var smallestSingleDirectory structs.DirectorySize

	for i := len(results) - 1; i >= 0; i-- {
		unusedSpace := currentUnusedSpace + results[i].Size

		// the moment the unused space can't satisfy the unusedTarget return the previous one
		if unusedSpace < unusedTarget {
			smallestSingleDirectory = results[i+1]
			break
		}
	}

	fmt.Println(smallestSingleDirectory)
}
