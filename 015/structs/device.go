package structs

import (
	"math"
)

type device struct {
	location                coordinates
	id                      deviceType
	closestBeacon           *device
	distanceToClosestBeacon int
}

func (d *device) taxiCapDistanceTo(c coordinates) int {
	return int((math.Abs(float64(d.location.x-c.x)) + math.Abs(float64(d.location.y-c.y))))
}
