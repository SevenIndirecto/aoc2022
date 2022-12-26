package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	N int = iota
	S
	W
	E
	NE
	NW
	SE
	SW
)

type point struct {
	x int
	y int
}

type grid [][]bool

type elf struct {
	p point
	nextPos point
	ignoreNextMove bool
}

func isOutOfBounds(x int, y int, g grid) bool {
	return y >= len(g) || x >= len(g[0]) || y < 0 || x < 0
}

func proposeDirections(g grid, elves []elf, offset int) (map[point]*elf, int) {
	directions := map[int]point {
		NE: {1,-1},
		N: {0, -1},
		NW: {-1, -1},
		W: {-1, 0},
		SW: {-1,1},
		S: {0, 1},
		SE: {1, 1},
		E: {1,0},
	}

	possibleMoves := [4][3]int{
		{N, NE, NW},
		{S, SE, SW},
		{W, NW, SW},
		{E, NE, SE},
	}

	firstElfToMoveToNextPoint := make(map[point]*elf)
	allClearCount := 0

	for elfIndex, e := range elves {
		e.ignoreNextMove = false

		allClear := true
		for dy := -1; dy <= 1; dy++ {
			if !allClear {
				break
			}
			for dx := -1; dx <= 1; dx++ {
				nx := e.p.x+dx
				ny := e.p.y+dy
				if (dx == 0 && dy == 0) || isOutOfBounds(nx, ny, g) {
					continue
				}
				if g[ny][nx] {
					allClear = false
					break
				}
			}
		}
		if allClear {
			allClearCount++
			continue
		}

		for i := 0; i < 4; i++ {
			m := possibleMoves[(i + offset) % 4]

			blockedMove := false
			for _, d := range m {
				ny := e.p.y + directions[d].y
				nx := e.p.x + directions[d].x

				if !isOutOfBounds(nx, ny, g) && g[ny][nx] {
					blockedMove = true
					break
				}
			}
			if blockedMove {
				continue
			}

			delta := directions[m[0]]
			nextPos := point{
				e.p.x + delta.x,
				e.p.y + delta.y,
			}

			_, anElfAlreadyMovingHere := firstElfToMoveToNextPoint[nextPos]
			if anElfAlreadyMovingHere {
				firstElfToMoveToNextPoint[nextPos].ignoreNextMove = true
				break
			}

			firstElfToMoveToNextPoint[nextPos] = &elves[elfIndex]
			firstElfToMoveToNextPoint[nextPos].nextPos = nextPos
			break
		}
	}

	return firstElfToMoveToNextPoint, allClearCount
}

func expandGrid(nx int, ny int, g grid, elves []elf) grid {
	if nx < 0 {
		for y, row := range g {
			g[y] = append([]bool{false}, row...)
		}
		// Shift all elves right
		for i := range elves {
			elves[i].p.x++
			elves[i].nextPos.x++
		}
	}

	if ny < 0 {
		newRow := make([]bool, len(g[0]))
		g = append(grid{newRow}, g...)

		for i := range elves {
			elves[i].p.y++
			elves[i].nextPos.y++
		}
	}

	if nx >= len(g[0]) {
		for y := range g {
			g[y] = append(g[y], false)
		}
	}

	if ny >= len(g) {
		g = append(g, make([]bool, len(g[0])))
	}
	return g
}

func processRounds(g grid, elves []elf, rounds int) (grid, int) {
	for i := 0; i < rounds; i++ {
		elvesThatMightMove, allClearCount := proposeDirections(g, elves, i)

		if allClearCount == len(elves) {
			return g, i+1
		}

		for _, e := range elvesThatMightMove {
			if e.ignoreNextMove {
				continue
			}

			if isOutOfBounds(e.nextPos.x, e.nextPos.y, g) {
				g = expandGrid(e.nextPos.x, e.nextPos.y, g, elves)
			}
		}

		for _, e := range elvesThatMightMove {
			if e.ignoreNextMove {
				e.ignoreNextMove = false
				continue
			}

			g[e.p.y][e.p.x] = false
			g[e.nextPos.y][e.nextPos.x] = true
			e.p = e.nextPos
		}
	}

	return g, -1
}

func getScore(g grid, elves []elf) int {
	minX := len(g[0])
	maxX := -1
	minY := len(g)
	maxY := -1

	for _, e := range elves {
		if e.p.x < minX {
			minX = e.p.x
		}
		if e.p.x > maxX {
			maxX = e.p.x
		}
		if e.p.y < minY {
			minY = e.p.y
		}
		if e.p.y > maxY {
			maxY = e.p.y
		}
	}

	count := 0
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if !g[y][x] {
				count++
			}
		}
	}
	return count
}

func drawMap(g grid) {
	for y := range g {
		for x := range g[y] {
			if g[y][x] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}

func PartOne(lines []string) int {
	g, elves := loadMap(lines)
	g, _ = processRounds(g, elves, 10)
	return getScore(g, elves)
}

func PartTwo(lines []string) int {
	g, elves := loadMap(lines)
	_, rounds := processRounds(g, elves, 100000000000)
	return rounds
}

func loadMap(lines []string) (grid, []elf) {
	g := make(grid, len(lines))
	elves := make([]elf, 0)

	for y, l := range lines {
		row := make([]bool, len(l))
		for x, c := range l {
			if c == '#' {
				elves = append(elves, elf{p: point{x, y}})
				row[x] = true
			}
		}
		g[y] = row
	}
	return g, elves
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
