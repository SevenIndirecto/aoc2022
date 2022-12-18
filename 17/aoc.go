package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	minus int = iota
	plus
	lShape
	vertical
	square
)

const (
	Air int = iota
	Rock
)

const (
	Right string = ">"
	Left string = "<"
)

type point struct {
	x int
	y int
}

type brick struct {
	points []point
	height int
}

func createBrick(s point, sequence int) brick {
	switch sequence {
	case minus:
		return brick{
			points: []point{
				{s.x, s.y},
				{s.x+1, s.y},
				{s.x+2, s.y},
				{s.x+3, s.y},
			},
			height: 1,
		}
	case plus:
		return brick{
			points: []point{
				{s.x+1, s.y},
				{s.x, s.y+1},
				{s.x+1, s.y+1},
				{s.x+2, s.y+1},
				{s.x+1, s.y+2},
			},
			height: 3,
		}
	case lShape:
		return brick{
			points: []point{
				{s.x, s.y},
				{s.x+1, s.y},
				{s.x+2, s.y},
				{s.x+2, s.y+1},
				{s.x+2, s.y+2},
			},
			height: 3,
		}
	case vertical:
		return brick{
			points: []point{
				{s.x, s.y},
				{s.x, s.y+1},
				{s.x, s.y+2},
				{s.x, s.y+3},
			},
			height: 4,
		}
	default:
		return brick{
			points: []point{
				{s.x, s.y},
				{s.x+1, s.y},
				{s.x, s.y+1},
				{s.x+1, s.y+1},
			},
			height: 2,
		}
	}
}

//func dropBrick(grid [][7]string, pattern *string, patternPos int, sequence int) patter

func expandGrid(grid [][7]int, b brick) ([][7]int, int) {
	// Find highest point
	highestRock := 0
	for y := len(grid)-1; y >= 0; y-- {
		foundY := false
		for _, cell := range grid[y] {
			if cell == Rock {
				foundY = true
				highestRock = y
				break
			}
		}
		if foundY {
			break
		}
	}

	if len(grid) == 0 {
		highestRock = -1
	}
	for i := len(grid); i <= highestRock + 3 + b.height; i++ {
		row := [7]int{}
		for idx := range row {
			row[idx] = Air
		}
		grid = append(grid, row)
	}
	return grid, highestRock
}

func printMap(grid [][7]int, b brick, skipBrick bool) {
	for y := len(grid)-1; y >= 0; y-- {
		for x := range grid[y] {
			drewBrick := false

			if !skipBrick {
				for _, p := range b.points {
					if p.x == x && p.y == y {
						fmt.Print("@")
						drewBrick = true
						break
					}
				}
			}

			if !drewBrick {
				if grid[y][x] == Rock {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
}

func PartOne(lines []string) int {
	grid := make([][7]int, 0)
	highestY := -1

	offset := -1
	pattern := lines[0]

	seq := 0
	for i := 0; i < 2022; i++ {
		// Just create a fake brick to allow expanding
		fauxBrick := createBrick(point{0, 0}, seq)
		grid, highestY = expandGrid(grid, fauxBrick)
		b := createBrick(point{x: 2, y: highestY + 4}, seq)
		//fmt.Println("SPAWN ----------")
		//printMap(grid, b, false)

		for {
			// Process wind
			offset++
			if offset >= len(pattern) {
				offset = 0
			}

			dx := 1
			if string(pattern[offset]) == Left {
				dx = -1
			}
			candidatePoints := make([]point, 0)
			for _, p := range b.points {
				nx := p.x + dx
				if nx < 0 || nx >= len(grid[p.y]) || grid[p.y][nx] == Rock {
					candidatePoints = make([]point, 0)
					break
				}
				candidatePoints = append(candidatePoints, point{nx, p.y})
			}
			if len(candidatePoints) > 0 {
				b.points = candidatePoints
			}

			// Move down by one
			candidatePoints = make([]point, 0)
			for _, p := range b.points {
				ny := p.y-1
				if ny < 0 || grid[ny][p.x] == Rock {
					candidatePoints = make([]point, 0)
					break
				}
				candidatePoints = append(candidatePoints, point{p.x, ny})
			}

			if len(candidatePoints) > 0 {
				b.points = candidatePoints
				//fmt.Println(i, string(pattern[offset]))
				//printMap(grid, b, false)
				//fmt.Println("----------")
			} else {
				// Reached bottom or rock
				for _, p := range b.points {
					grid[p.y][p.x] = Rock
				}
				//fmt.Println(i, string(pattern[offset]))
				//printMap(grid, b, true)
				//fmt.Println("----------")
				break
			}
		}

		seq++
		seq %= 5
	}


	highestY = -1
	foundY := false
	for y := len(grid)-1; y >= 0; y-- {
		for _, cell := range grid[y] {
			if cell == Rock {
				foundY = true
				highestY = y
				break
			}
		}
		if foundY {
			break
		}
	}

	return highestY+1
}

func PartTwo(lines []string) int {
	return 0
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
	// 3195 too high
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
