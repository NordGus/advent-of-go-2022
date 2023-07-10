package structs

import (
	"fmt"
	"log"
	"math"
)

type deviceType int

const (
	sensor deviceType = iota
	beacon
)

type Grid struct {
	devices map[coordinates]*device
}

func NewGrid() Grid {
	return Grid{
		devices: make(map[coordinates]*device, 1_000),
	}
}

func (sg *Grid) AddSensor(sensorX int, sensorY int, beaconX int, beaconY int) {
	sensorCoordinates := coordinates{x: sensorX, y: sensorY}
	beaconCoordinates := coordinates{x: beaconX, y: beaconY}

	if sg.devices[sensorCoordinates] != nil {
		log.Fatalf("can't add sensor to location %v, there's another device there", sensorCoordinates)
	}

	sen := device{
		location: sensorCoordinates,
		id:       sensor,
	}

	bea := sg.devices[beaconCoordinates]

	if bea != nil && bea.id != beacon {
		log.Fatalf("mismatched device %v in location %v", &bea, beaconCoordinates)
	}

	if bea == nil || bea.id != beacon {
		bea = &device{
			location: beaconCoordinates,
			id:       beacon,
		}

		sg.devices[beaconCoordinates] = bea
	}

	sen.closestBeacon = bea
	sen.distanceToClosestBeacon = sen.taxiCapDistanceTo(bea.location)

	sg.devices[sensorCoordinates] = &sen
}

func (sg *Grid) HowManyPositionsCannotContainABeaconAt(y int) (uint, error) {
	minX := -math.MaxInt
	maxX := math.MaxInt

	for _, device := range sg.devices {
		// Ignore beacons, because they don't provide the necessary information
		if device.id == beacon {
			continue
		}

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

func (sg *Grid) Print() {
	for location, device := range sg.devices {
		if device.closestBeacon != nil {
			fmt.Println(location, device, *device.closestBeacon)
		} else {
			fmt.Println(location, device)
		}
	}
}
