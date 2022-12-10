package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type knot struct {
	x int
	y int
	visited map[string]bool
}

type moveInstruction struct {
	x int
	y int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (k *knot) markVisit() {
	key := strconv.Itoa(k.x) + "," + strconv.Itoa(k.y)
	k.visited[key] = true
}

func processMoves(move string, knots []knot) {
	moveInstructions := map[string]moveInstruction{
		"U": {x: 0, y: 1},
		"R": {x: 1, y: 0},
		"D": {x: 0, y: -1},
		"L": {x: -1, y: 0},
	}

	s := strings.Split(move, " ")
	m := moveInstructions[s[0]]
	distance, _ := strconv.Atoi(s[1])

	for i := 0; i < distance; i++ {
		// Move Head
		knots[0].x += m.x
		knots[0].y += m.y

		for j := 1; j < len(knots); j++ {
			snap(&knots[j-1], &knots[j])
		}
	}
}

func snap(h *knot, t *knot) {
	dx := abs(t.x-h.x)
	dy := abs(t.y-h.y)

	if dx == 2 && dy == 2 {
		// Can make this cleaner, but sleepy : D
		if t.x > h.x && t.y > h.y {
			t.x--
			t.y--
		} else if t.x > h.x && t.y < h.y {
			t.x--
			t.y++
		} else if t.x < h.x && t.y < h.y {
			t.x++
			t.y++
		} else {
			t.x++
			t.y--
		}
	} else if dx == 2 {
		if t.x < h.x {
			t.x++
		} else {
			t.x--
		}
		t.y = h.y
	} else if dy == 2 {
		if t.y < h.y {
			t.y++
		} else {
			t.y--
		}
		t.x = h.x
	}
	t.markVisit()
}

func PartOne(lines []string) int {
	t := knot{x: 0, y: 0, visited: map[string]bool{"0,0": true}}
	h := knot{x: 0, y: 0}
	knots := []knot{h, t}

	for _, move := range lines {
		processMoves(move, knots)
	}

	return len(t.visited)
}

func PartTwo(lines []string) int {
	knots := make([]knot, 0)
	for i := 0; i < 10; i++ {
		knots = append(knots, knot{x: 0, y: 0, visited: map[string]bool{"0,0": true}})
	}

	for _, move := range lines {
		processMoves(move, knots)
	}

	return len(knots[9].visited)
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
