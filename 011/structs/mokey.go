package structs

import (
	"errors"
	"fmt"
	"sort"
)

type Monkey struct {
	items               list
	operation           func(int) int
	divisibleBy         int
	testSuccessTarget   int
	testFailTarget      int
	InspectedItemsCount uint
}

func NewMonkey() Monkey {
	return Monkey{}
}

func (m *Monkey) AddItems(items ...int) {
	for _, item := range items {
		m.items.push(item)
	}
}

func (m *Monkey) SetOperation(operation func(int) int) error {
	if m.operation != nil {
		return errors.New("operation already set")
	}

	m.operation = operation

	return nil
}

func (m *Monkey) SetTest(test int) {
	m.divisibleBy = test
}

func (m *Monkey) SetTestSuccessTarget(target int) {
	m.testSuccessTarget = target
}

func (m *Monkey) SetTestFailTarget(target int) {
	m.testFailTarget = target
}

func (m *Monkey) Items() []int {
	return m.items.toSlice()
}

func (m *Monkey) Print(number int) {
	fmt.Println("Monkey", number, "items:", m.items.toSlice(), "items inspected", m.InspectedItemsCount)
}

func PlayRound(monkeys []Monkey, modulo int) {
	for i := 0; i < len(monkeys); i++ {
		for {
			item, err := monkeys[i].items.pop()
			if err != nil {
				break
			}
			monkeys[i].InspectedItemsCount++

			newItem := monkeys[i].operation(item) % modulo

			if newItem%monkeys[i].divisibleBy == 0 {
				monkeys[monkeys[i].testSuccessTarget].items.push(newItem)
			} else {
				monkeys[monkeys[i].testFailTarget].items.push(newItem)
			}

		}
	}
}

func GetLimiter(monkeys []Monkey) int {
	limiter := 1

	for _, monkey := range monkeys {
		limiter *= monkey.divisibleBy
	}

	return limiter
}

func SortByItemInspectedCount(monkeys []Monkey) []Monkey {
	sort.Sort(sort.Reverse(byItemInspectedCount(monkeys)))

	return monkeys
}
