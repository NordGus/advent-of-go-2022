package structs

import (
	"context"
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
)

type deviceType int

const (
	sensor deviceType = iota
	beacon

	frequencyMultiplier = 4_000_000
)

type processStatus struct {
	sync.Mutex
	stopped bool
}

type Grid struct {
	sensors map[coordinates]*device
	beacons map[coordinates]*device
}

func NewGrid() Grid {
	return Grid{
		sensors: make(map[coordinates]*device, 1_000),
		beacons: make(map[coordinates]*device, 1_000),
	}
}

func (sg *Grid) AddSensor(sensorX int, sensorY int, beaconX int, beaconY int) {
	sensorCoordinates := coordinates{x: sensorX, y: sensorY}
	beaconCoordinates := coordinates{x: beaconX, y: beaconY}

	if sg.sensors[sensorCoordinates] != nil {
		log.Fatalf("can't add sensor to location %v, there's another device there", sensorCoordinates)
	}

	sen := device{
		location: sensorCoordinates,
		id:       sensor,
	}

	bea := sg.beacons[beaconCoordinates]

	if bea == nil {
		bea = &device{
			location: beaconCoordinates,
			id:       beacon,
		}

		sg.beacons[beaconCoordinates] = bea
	}

	sen.closestBeacon = bea
	sen.distanceToClosestBeacon = sen.taxiCapDistanceTo(bea.location)

	sg.sensors[sensorCoordinates] = &sen
}

func (sg *Grid) HowManyPositionsCannotContainABeaconAt(y int) (uint, error) {
	minX := -math.MaxInt
	maxX := math.MaxInt

	for _, device := range sg.sensors {
		start := coordinates{x: device.location.x, y: y}

		// Ignore sensors which exclusion zone is not dissected by the line defined by the parameter y
		if device.taxiCapDistanceTo(start) > device.distanceToClosestBeacon {
			continue
		}

		// catch first iteration
		if minX == -math.MaxInt {
			minX = start.x
		}

		if maxX == math.MaxInt {
			maxX = start.x
		}

		// find left dissection point point
		left := start

		for device.distanceToClosestBeacon >= device.taxiCapDistanceTo(left) {
			if minX > left.x {
				minX = left.x
			}

			left = coordinates{x: left.x - 1, y: left.y}
		}

		// find right dissection point point
		right := start

		for device.distanceToClosestBeacon >= device.taxiCapDistanceTo(right) {
			if maxX < right.x {
				maxX = right.x
			}

			right = coordinates{x: right.x + 1, y: right.y}
		}
	}

	if minX == -math.MaxInt && maxX == math.MaxInt {
		return 0, fmt.Errorf("the line define by y=%v doesn't dissect any sensor's exclusion zones", y)
	}

	return uint(maxX - minX), nil
}

func (sg *Grid) TuningFrequencyOfOfDistressBeacon(lower int, upper int) (uint, error) {
	workers := runtime.GOMAXPROCS(0) * 10
	wg := new(sync.WaitGroup)
	sem := make(chan bool, workers)
	results := make(chan coordinates, workers)
	ctx, cancel := context.WithCancel(context.Background())
	finished := new(processStatus)

	defer cancel()

	wg.Add(1)

	// work scheduler
	go func(ctx context.Context, wg *sync.WaitGroup, sem chan bool, out chan<- coordinates, workers int, lower int, upper int, finished *processStatus) {
		defer wg.Done()

		for y := lower; y <= upper; y++ {
			finished.Lock()
			if finished.stopped {
				break
			}
			finished.Unlock()

			wg.Add(1)
			sem <- true

			go func(ctx context.Context, wg *sync.WaitGroup, sem chan bool, results chan<- coordinates, y int, lower int, upper int) {
				defer wg.Done()
				defer func() {
					<-sem
				}()

				location := coordinates{x: -1, y: y}
				line := make([]bool, upper+1)

				for _, device := range sg.sensors {
					start := coordinates{x: device.location.x, y: y}

					// Ignore sensors which exclusion zone is not dissected by the line defined by the parameter y
					if device.taxiCapDistanceTo(start) > device.distanceToClosestBeacon {
						continue
					}

					// find left dissection point point
					left := start

					for device.distanceToClosestBeacon > device.taxiCapDistanceTo(left) && left.x > lower {
						left = coordinates{x: left.x - 1, y: left.y}
					}

					// find right dissection point point
					right := start

					for device.distanceToClosestBeacon > device.taxiCapDistanceTo(right) && right.x < upper {
						right = coordinates{x: right.x + 1, y: right.y}
					}

					if left.x < lower {
						left.x = lower
					}

					if right.x > upper {
						right.x = upper
					}

					for i := left.x; i <= right.x; i++ {
						if line[i] {
							continue
						}

						line[i] = true
					}
				}

				for i := 0; i < len(line); i++ {
					if !line[i] {
						location.x = i
						break
					}
				}

				select {
				case <-ctx.Done():
					return
				case results <- location:
				}
			}(ctx, wg, sem, out, y, lower, upper)
		}
	}(ctx, wg, sem, results, workers, lower, upper, finished)

	go func(wg *sync.WaitGroup, sem chan bool, out chan coordinates) {
		wg.Wait()
		close(sem)
		close(out)
	}(wg, sem, results)

	for result := range results {
		if result.x > -1 {
			finished.Lock()
			finished.stopped = true
			finished.Unlock()
			return uint(result.x*frequencyMultiplier + result.y), nil
		}
	}

	return 0, fmt.Errorf("location not found for distress beacon in area where x and y between %v and %v", lower, upper)
}
