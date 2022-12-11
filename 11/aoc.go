package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type monkey struct {
	id int
	items []int
	inspect func(x int) int
	test func(x int) int
	divBy int
	inspectCount int
}

func loadMonkeys(lines []string) []monkey {
	monkeys := make([]monkey, 0)

	i := 0
	m := monkey{}
	monkeyId := 0

	for idx := 0; idx < len(lines); idx++ {
		l := lines[idx]
		i++

		if len(l) < 2 {
			monkeys = append(monkeys, m)
			i = 0
			continue
		}

		if i == 1 {
			m = monkey{
				id: monkeyId,
				items: make([]int, 0),
				inspectCount: 0,
			}
			monkeyId++
		} else if i == 2 {
			s := strings.Split(l, ": ")
			s = strings.Split(s[1], ", ")

			for _, item := range s {
				worryLevel, _ := strconv.Atoi(item)
				m.items = append(m.items, worryLevel)
			}
		} else if i == 3 {
			if len(strings.Split(l, "old * old")) > 1 {
				m.inspect = func(x int) int {
					return x * x
				}
			} else {
				s := strings.Split(l, "* ")
				if len(s) > 1 {
					p, _ := strconv.Atoi(s[len(s)-1])
					m.inspect = func(x int) int {
						return x * p
					}
				} else {
					s = strings.Split(l, "+ ")
					p, _ := strconv.Atoi(s[len(s)-1])
					m.inspect = func(x int) int {
						return x + p
					}
				}
			}
		} else if i == 4 {
			// Read next 3 lines
			// get divisible by
			s := strings.Split(l, " ")
			divisibleBy, _ := strconv.Atoi(s[len(s)-1])
			m.divBy = divisibleBy

			s = strings.Split(lines[idx+1], " ")
			successTarget, _ := strconv.Atoi(s[len(s)-1])

			s = strings.Split(lines[idx+2], " ")
			failTarget, _ := strconv.Atoi(s[len(s)-1])

			m.test = func(x int) int {
				if x % divisibleBy == 0 {
					return successTarget
				}
				return failTarget
			}
			idx+=2
		}
	}
	monkeys = append(monkeys, m)
	return monkeys
}

func executeRounds(rounds int, monkeys []monkey, divByThree bool) {
	worryLevelContainer := 1
	for _, m := range monkeys {
		worryLevelContainer *= m.divBy
	}

	for idx := 0; idx < rounds; idx++ {
		for i := range monkeys {
			m := &monkeys[i]
			for _, worryLevel := range m.items {
				worryLevel = m.inspect(worryLevel)
				m.inspectCount++

				if divByThree {
					worryLevel /= 3
				} else {
					worryLevel %= worryLevelContainer
				}

				targetMonkey := m.test(worryLevel)
				monkeys[targetMonkey].items = append(monkeys[targetMonkey].items, worryLevel)
			}
			m.items = make([]int, 0)
		}
	}
}

func PartOne(lines []string) int {
	monkeys := loadMonkeys(lines)
	executeRounds(20, monkeys, true)

	ic := make([]int, 0)
	for _, m := range monkeys {
		ic = append(ic, m.inspectCount)
	}
	sort.Ints(ic)
	return ic[len(ic)-1] * ic[len(ic)-2]
}

func PartTwo(lines []string) int {
	monkeys := loadMonkeys(lines)
	executeRounds(10000, monkeys, false)
	ic := make([]int, 0)
	for _, m := range monkeys {
		ic = append(ic, m.inspectCount)
	}
	sort.Ints(ic)
	return ic[len(ic)-1] * ic[len(ic)-2]
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
