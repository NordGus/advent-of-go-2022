package part1

import "errors"

var (
	InvalidResourceErr = errors.New("blueprint: invalid resource")
)

type Resource uint

const (
	Invalid Resource = iota
	Ore
	Clay
	Obsidian
	Geode
)

func getResource(resource string) (Resource, error) {
	switch resource {
	case "geode":
		return Geode, nil
	case "obsidian":
		return Obsidian, nil
	case "clay":
		return Clay, nil
	case "ore":
		return Ore, nil
	default:
		return Invalid, InvalidResourceErr
	}
}
