package structs

import (
	"errors"
	"fmt"
	"sort"
)

type Monkey struct {
	items               list
	operation           func(int) int
	test                func(int) bool
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

func (m *Monkey) SetTest(test func(int) bool) error {
	if m.test != nil {
		return errors.New("test already set")
	}

	m.test = test

	return nil
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

func PlayRound(monkeys []Monkey) {
	for i := 0; i < len(monkeys); i++ {
		for {
			item, err := monkeys[i].items.pop()
			if err != nil {
				break
			}
			monkeys[i].InspectedItemsCount++

			newItem := monkeys[i].operation(item) / 3

			if monkeys[i].test(newItem) {
				monkeys[monkeys[i].testSuccessTarget].items.push(newItem)
			} else {
				monkeys[monkeys[i].testFailTarget].items.push(newItem)
			}

		}
	}
}

func SortByItemInspectedCount(monkeys []Monkey) []Monkey {
	sort.Sort(sort.Reverse(byItemInspectedCount(monkeys)))

	return monkeys
}
