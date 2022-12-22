package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	Mul int = iota
	Sub
	Add
	Div
)

type monkey struct {
	id string
	hasValue bool
	op int
	value int
	originalValue int
	originalHasValue bool
	leftOperand string
	rightOperand string
}

func (m *monkey) solve(monkeys map[string]*monkey) bool {
	if m.hasValue {
		return true
	}
	if monkeys[m.leftOperand].hasValue && monkeys[m.rightOperand].hasValue {
		r := monkeys[m.leftOperand].value
		l := monkeys[m.rightOperand].value
		switch m.op {
		case Mul: m.value = r * l
		case Div: m.value = r / l
		case Sub: m.value = r - l
		case Add: m.value = r + l
		}
		m.hasValue = true
		return true
	}
	return false
}

func (m *monkey) toString(monkeys map[string]*monkey) string {
	return fmt.Sprintf(
		"[%v] hasValue: %v, value: %v, op: %v, left: [%v, %v], right: [%v, %v]",
		m.id,
		m.hasValue,
		m.value,
		m.op,
		monkeys[m.leftOperand].hasValue,
		monkeys[m.leftOperand].value,
		monkeys[m.rightOperand].hasValue,
		monkeys[m.rightOperand].value,
	)
}

func getMonkeyValue(target string, monkeys map[string]*monkey, getDiff bool) int {
	if monkeys[target].hasValue {
		return monkeys[target].value
	}
	//fmt.Println("Human monkey", monkeys["humn"].hasValue, monkeys["humn"].value)

	unsolvedMonkeys := make(map[string]bool)
	for monkeyId, m := range monkeys {
		if !m.hasValue && !(getDiff && monkeyId == "root") {
			unsolvedMonkeys[monkeyId] = true
		}
	}

	root := monkeys["root"]

	for len(unsolvedMonkeys) > 0 || getDiff {
		if !getDiff && monkeys[target].hasValue {
			return monkeys[target].value
		}

		if getDiff && monkeys[root.leftOperand].hasValue && monkeys[root.rightOperand].hasValue {
			return monkeys[root.leftOperand].value - monkeys[root.rightOperand].value
		}

		for id := range unsolvedMonkeys {
			if monkeys[id].solve(monkeys) {
				delete(unsolvedMonkeys, id)
			}
		}
	}

	return monkeys[target].value
}

func findRequiredNumber(monkeys map[string]*monkey) int {
	low := 0
	high := 379578518396784
	diff := -1
	testValue := -1

	monkeys["humn"].hasValue = true
	monkeys["humn"].value = low
	valueAtLow := getMonkeyValue("root", monkeys, true)
	valueAtLowPositive := valueAtLow > 0

	for diff != 0 {
		testValue = (low + high) / 2
		for i := range monkeys {
			monkeys[i].value = monkeys[i].originalValue
			monkeys[i].hasValue = monkeys[i].originalHasValue
		}
		monkeys["humn"].hasValue = true
		monkeys["humn"].value = testValue
		diff = getMonkeyValue("root", monkeys, true)

		if diff == 0 {
			return testValue
		}
		if valueAtLowPositive {
			if diff > 0 {
				low = testValue
			} else {
				high = testValue
			}
		} else {
			if diff < 0 {
				low = testValue
			} else {
				high = testValue
			}
		}
	}

	return testValue
}

func loadMonkeys(lines []string) map[string]*monkey {
	monkeys := make(map[string]*monkey)
	opMap := map[string]int{
		"*": Mul,
		"/": Div,
		"-": Sub,
		"+": Add,
	}

	for _, l := range lines {
		s := strings.Split(l, ": ")
		monkeyId := s[0]
		m := monkey{hasValue: false, originalHasValue: false, id: monkeyId}

		if len(s[1]) == 11 {
			m.leftOperand = s[1][:4]
			m.rightOperand = s[1][7:]
			m.op = opMap[string(s[1][5])]
		} else {
			v, _ := strconv.Atoi(s[1])
			m.value = v
			m.originalValue = v
			m.hasValue = true
			m.originalHasValue = true
		}
		monkeys[monkeyId] = &m
	}

	return monkeys
}

func PartOne(lines []string) int {
	monkeys := loadMonkeys(lines)
	return getMonkeyValue("root", monkeys, false)
}

func PartTwo(lines []string) int {
	monkeys := loadMonkeys(lines)
	return findRequiredNumber(monkeys)
}

func LoadLines(path string) ([]string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	txt := string(dat)
	lines := strings.Split(txt, "\n")
	return lines[:len(lines)-1], nil
}

func main() {
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
