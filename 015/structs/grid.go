package structs

import (
	"context"
	"fmt"
	"log"
	"math"
	"runtime"
	"sort"
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
		ds := device.taxiCapDistanceTo(start)

		// Ignore sensors which exclusion zone is not dissected by the line defined by the parameter y
		if ds > device.distanceToClosestBeacon {
			continue
		}

		// catch first iteration
		if minX == -math.MaxInt {
			minX = start.x
		}

		if maxX == math.MaxInt {
			maxX = start.x
		}

		remainder := int(math.Abs(float64(ds - device.distanceToClosestBeacon)))

		left := start.x - remainder
		right := start.x + remainder

		if minX > left {
			minX = left
		}

		if maxX < right {
			maxX = right
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
				defer func(wg *sync.WaitGroup, sem chan bool) {
					wg.Done()
					<-sem
				}(wg, sem)

				location := coordinates{x: -1, y: y}
				ranges := make([]rang, 0, len(sg.sensors))

				for _, device := range sg.sensors {
					start := coordinates{x: device.location.x, y: y}
					ds := device.taxiCapDistanceTo(start)

					// Ignore sensors which exclusion zone is not dissected by the line defined by the parameter y
					if ds > device.distanceToClosestBeacon {
						continue
					}

					remainder := int(math.Abs(float64(ds - device.distanceToClosestBeacon)))

					left := start.x - remainder
					right := start.x + remainder

					if left < lower {
						left = lower
					}

					if right > upper {
						right = upper
					}

					ranges = append(ranges, rang{start: left, end: right})
				}

				sort.Slice(ranges, func(i, j int) bool {
					return (ranges[i].start < ranges[j].start)
				})

				fusions := make([]rang, 0, len(ranges))
				fusion := ranges[0]

				for i := 1; i < len(ranges); i++ {
					if fusion.canMerge(ranges[i]) {
						fusion = fusion.merge(ranges[i])
						continue
					}

					fusions = append(fusions, fusion)
					fusion = ranges[i]
				}

				fusions = append(fusions, fusion)

				sort.Slice(fusions, func(i, j int) bool {
					return (fusions[i].start < fusions[j].start)
				})

				if len(fusions) > 1 {
					if len(fusions) > 2 {
						fmt.Println(fusions)
					}

					location.x = fusions[0].end + 1
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
		sem <- true
		wg.Wait()
		close(sem)
		close(out)
		<-sem
	}(wg, sem, results)

	for result := range results {
		if result.x > -1 && sg.isDistressBeacon(result) {
			finished.Lock()
			finished.stopped = true
			finished.Unlock()
			return uint(result.x*frequencyMultiplier + result.y), nil
		}
	}

	return 0, fmt.Errorf("location not found for distress beacon in area where x and y between %v and %v", lower, upper)
}

func (sg *Grid) isDistressBeacon(location coordinates) bool {
	for _, device := range sg.sensors {
		if device.taxiCapDistanceTo(location) <= device.distanceToClosestBeacon {
			return false
		}
	}

	return true
}
