package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type assignment struct {
	start int
	end   int
}

type pair [2]assignment

func parseElves(lines []string) []pair {
	pairs := make([]pair, 0)

	for _, p := range lines {
		elves := strings.Split(p, ",")

		pair := pair{}
		for i := 0; i < 2; i++ {
			split := strings.Split(elves[i], "-")
			start, _ := strconv.Atoi(split[0])
			end, _ := strconv.Atoi(split[1])

			pair[i] = assignment{start: start, end: end}
		}
		pairs = append(pairs, pair)
	}
	return pairs
}

func PartOne(lines []string) int {
	count := 0
	pairs := parseElves(lines)

	for _, asigns := range pairs {
		var long = 0
		var short = 1
		if asigns[0].end-asigns[0].start < asigns[1].end-asigns[1].start {
			long = 1
			short = 0
		}

		if asigns[long].start <= asigns[short].start && asigns[long].end >= asigns[short].end {
			count++
		}
	}
	return count
}

func PartTwo(lines []string) int {
	pairs := parseElves(lines)
	count := 0
	for _, elves := range pairs {
		flips := [2]bool{false, true}

		for _, flip := range flips {
			a := 0
			b := 1
			if flip {
				a = 1
				b = 0
			}

			if (elves[a].start <= elves[b].start && elves[b].start <= elves[a].end) ||
				(elves[a].start <= elves[b].end && elves[b].end <= elves[a].end) {
				count++
				break
			}
		}
	}
	return count
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
	lines, _ := LoadLines("aoc04.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
