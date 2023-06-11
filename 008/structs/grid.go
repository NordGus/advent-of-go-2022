package structs

import (
	"errors"
	"sync"
)

type Grid struct {
	sync.Mutex
	xSize int
	ySize int
	trees [][]*Tree
	count int
	cap   int
}

func NewGrid(xSize int, ySize int) Grid {
	cap := xSize * ySize
	trees := make([][]*Tree, xSize)

	for i := 0; i < len(trees); i++ {
		trees[i] = make([]*Tree, ySize)
	}

	return Grid{
		xSize: xSize,
		ySize: ySize,
		trees: trees,
		cap:   cap,
	}
}

func (g *Grid) AddTree(x int, y int, height int) error {
	g.Lock()
	defer g.Unlock()

	if g.count >= g.cap {
		return errors.New("grid is full")
	}

	if g.trees[x][y] != nil {
		return errors.New("grid element overwrite")
	}

	tree := Tree{xCoordinate: x, yCoordinate: y, height: height}

	g.trees[x][y] = &tree
	g.count++

	return nil
}

func (g *Grid) Size() int {
	return g.count
}

func (g *Grid) DetectVisibleTrees(wg *sync.WaitGroup, out chan<- bool) {
	defer wg.Done()

	for _, row := range g.trees {
		for _, tree := range row {
			wg.Add(1)
			go g.treeIsVisible(wg, tree, out)
		}
	}
}

func (g *Grid) treeIsVisible(wg *sync.WaitGroup, tree *Tree, out chan<- bool) {
	defer wg.Done()

	if tree.xCoordinate == 0 || tree.xCoordinate == g.xSize-1 {
		out <- true
		return
	}

	if tree.yCoordinate == 0 || tree.yCoordinate == g.ySize-1 {
		out <- true
		return
	}

	signals := make(chan bool, 4)
	detected := false

	go func(signals chan<- bool) {
		var innerWg sync.WaitGroup

		innerWg.Add(1)
		// Detect West
		go func(wg *sync.WaitGroup, signals chan<- bool) {
			defer wg.Done()

			for i := tree.xCoordinate - 1; i >= 0; i-- {
				if !tree.IsTaller(*g.trees[i][tree.yCoordinate]) {
					signals <- false
					return
				}
			}

			signals <- true
		}(&innerWg, signals)

		innerWg.Add(1)
		// Detect East
		go func(wg *sync.WaitGroup, signals chan<- bool) {
			defer wg.Done()

			for i := tree.xCoordinate + 1; i < g.xSize; i++ {
				if tree.IsTaller(*g.trees[i][tree.yCoordinate]) {
					continue
				}

				signals <- false
				return
			}

			signals <- true
		}(&innerWg, signals)

		innerWg.Add(1)
		// Detect North
		go func(wg *sync.WaitGroup, signals chan<- bool) {
			defer wg.Done()

			for i := tree.yCoordinate - 1; i >= 0; i-- {
				if tree.IsTaller(*g.trees[tree.xCoordinate][i]) {
					continue
				}

				signals <- false
				return
			}

			signals <- true
		}(&innerWg, signals)

		innerWg.Add(1)
		// Detect South
		go func(wg *sync.WaitGroup, signals chan<- bool) {
			defer wg.Done()

			for i := tree.yCoordinate + 1; i < g.ySize; i++ {
				if tree.IsTaller(*g.trees[tree.xCoordinate][i]) {
					continue
				}

				signals <- false
				return
			}

			signals <- true
		}(&innerWg, signals)

		innerWg.Wait()
		close(signals)
	}(signals)

	for signal := range signals {
		if signal {
			detected = true
		}
	}

	out <- detected
}

func (g *Grid) CalculateTreesScenicScore(wg *sync.WaitGroup, out chan<- int) {
	defer wg.Done()

	for _, row := range g.trees {
		for _, tree := range row {
			wg.Add(1)
			go g.treeScenicScore(wg, tree, out)
		}
	}
}

func (g *Grid) treeScenicScore(wg *sync.WaitGroup, tree *Tree, out chan<- int) {
	defer wg.Done()

	if tree.xCoordinate == 0 || tree.xCoordinate == g.xSize-1 {
		out <- 0
		return
	}

	if tree.yCoordinate == 0 || tree.yCoordinate == g.ySize-1 {
		out <- 0
		return
	}

	scores := make(chan int, 4)
	score := 1

	go func(scores chan<- int) {
		var innerWg sync.WaitGroup

		innerWg.Add(1)
		// Detect West
		go func(wg *sync.WaitGroup, scores chan<- int) {
			defer wg.Done()
			distance := 0

			for i := tree.xCoordinate - 1; i >= 0; i-- {
				distance++

				if !tree.IsTaller(*g.trees[i][tree.yCoordinate]) {
					break
				}
			}

			scores <- distance
		}(&innerWg, scores)

		innerWg.Add(1)
		// Detect East
		go func(wg *sync.WaitGroup, scores chan<- int) {
			defer wg.Done()
			distance := 0

			for i := tree.xCoordinate + 1; i < g.xSize; i++ {
				distance++

				if !tree.IsTaller(*g.trees[i][tree.yCoordinate]) {
					break
				}
			}

			scores <- distance
		}(&innerWg, scores)

		innerWg.Add(1)
		// Detect North
		go func(wg *sync.WaitGroup, scores chan<- int) {
			defer wg.Done()
			distance := 0

			for i := tree.yCoordinate - 1; i >= 0; i-- {
				distance++

				if !tree.IsTaller(*g.trees[tree.xCoordinate][i]) {
					break
				}
			}

			scores <- distance
		}(&innerWg, scores)

		innerWg.Add(1)
		// Detect South
		go func(wg *sync.WaitGroup, scores chan<- int) {
			defer wg.Done()
			distance := 0

			for i := tree.yCoordinate + 1; i < g.ySize; i++ {
				distance++

				if !tree.IsTaller(*g.trees[tree.xCoordinate][i]) {
					break
				}
			}

			scores <- distance
		}(&innerWg, scores)

		innerWg.Wait()
		close(scores)
	}(scores)

	for ss := range scores {
		score *= ss
	}

	tree.scenicScore = score

	out <- score
}
