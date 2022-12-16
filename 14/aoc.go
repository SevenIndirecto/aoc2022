package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type node struct {
	x int
	y int
}

const (
	Air int = iota
	Rock
	Sand
	Source
)

func loadMap(lines []string, isPartTwo bool) ([][]int, int) {
	paths := make([][]node, 0)
	maxX := -1
	minX := -1
	maxY := -1

	// parse
	for _, l := range lines {
		s := strings.Split(l, " -> ")
		p := make([]node, 0)
		for _, coords := range s {
			c := strings.Split(coords, ",")
			x, _ := strconv.Atoi(c[0])
			y, _ := strconv.Atoi(c[1])
			p = append(p, node{x: x, y: y})

			if maxX == -1 || x > maxX {
				maxX = x
			}
			if minX == -1 || x < minX {
				minX = x
			}
			if maxY == -1 || maxY < y {
				maxY = y
			}
		}
		paths = append(paths, p)
	}

	width := maxX - minX + 1

	// Fill with air
	m := make([][]int, 0)
	for y := 0; y <= maxY; y++ {
		row := make([]int, 0)

		for x := 0; x < width; x++ {
			row = append(row, Air)
		}
		m = append(m, row)
	}

	fmt.Println("MinX", minX, "MaxX", maxX, "MaxY", maxY)

	// Create paths
	for _, p := range paths {
		for i, _ := range p {
			x := p[i].x - minX
			y := p[i].y

			m[y][x] = Rock

			if i+1 >= len(p) {
				break
			}

			nx := p[i+1].x - minX
			ny := p[i+1].y

			dx := 0
			dy := 0
			distance := 0

			// note: larger Y is further down
			// Fill lines
			if x == nx {
				dx = 0
				if y < ny {
					// vertical down
					dy = 1
					distance = ny - y
				} else {
					// vertical up
					dy = -1
					distance = y - ny
				}
			} else {
				dy = 0
				if x < nx {
					// right
					dx = 1
					distance = nx - x
				} else {
					// left
					dx = -1
					distance = x - nx
				}
			}

			for d := 0; d < distance; d++ {
				m[y+d*dy][x+d*dx] = Rock
			}
		}
	}

	// Create source
	sourceX := 500 - minX
	m[0][sourceX] = Source

	if isPartTwo {
		emptyRow := make([]int, len(m[0]))
		for x := 0; x < len(m[0]); x++ {
			emptyRow[x] = Air
		}
		floor := make([]int, len(m[0]))
		for x := 0; x < len(m[0]); x++ {
			floor[x] = Rock
		}
		m = append(m, emptyRow, floor)
	}

	return m, sourceX
}

func drawMap(m [][]int) {
	smap := map[int]string {
		Air: ".",
		Rock: "#",
		Source: "+",
		Sand: "o",
	}

	for y, _ := range m {
		for x, _ := range m[y] {
			fmt.Print(smap[m[y][x]])
		}
		fmt.Println()
	}
}

func dropSand(m [][]int, sourceX int) (isOffMap bool) {
	s := node{x: sourceX, y: 0}

	for {
		if s.y+1 >= len(m) {
			// Reached floor
			return true
		}
		if m[s.y+1][s.x] == Air {
			// Can move below?
			s.y++
		} else if s.x-1 < 0 {
			return true
		} else if m[s.y+1][s.x-1] == Air {
			// Can move diagonal left
			s.x--
			s.y++
		} else if s.x+1 >= len(m[0]) {
			return true
		} else if m[s.y+1][s.x+1] == Air {
			s.x++
			s.y++
		} else {
			if s.x-1 < 0 || s.x+1 >= len(m[0]) {
				return true
			}
			// Could not move anywhere else
			m[s.y][s.x] = Sand
			break
		}
	}
	return false
}

func expandMap(m [][]int, expandLeft bool) {
	for y := range m {
		fill := Air
		if y == len(m)-1 {
			fill = Rock
		}

		if expandLeft {
			m[y] = append([]int{fill}, m[y]...)
		} else {
			m[y] = append(m[y], fill)
		}
	}
}

func dropSandPartTwo(m [][]int, sourceX int) (newSourceX int) {
	// TODO: Does the source move when expanding map?
	s := node{x: sourceX, y: 0}
	newSourceX = sourceX

	for {
		if s.y+1 == len(m) || (
				s.x+1 < len(m[0]) &&
				s.x-1 >= 0 &&
				m[s.y+1][s.x] != Air &&
				m[s.y+1][s.x-1] != Air &&
				m[s.y+1][s.x+1] != Air) {
			m[s.y][s.x] = Sand
			break
		} else if m[s.y+1][s.x] == Air {
			// Can move below?
			s.y++
		} else if s.x-1 < 0 {
			expandMap(m, true)
			s.x++
			m[0][newSourceX] = Air
			newSourceX++
			m[0][newSourceX] = Source
			continue
		} else if m[s.y+1][s.x-1] == Air {
			// Can move diagonal left
			s.x--
			s.y++
		} else if s.x+1 >= len(m[0]) {
			expandMap(m, false)
			continue
		} else if m[s.y+1][s.x+1] == Air {
			// Can move diagonal right
			s.x++
			s.y++
		} else {
			// Could not move anywhere else
			m[s.y][s.x] = Sand
			break
		}
	}
	return newSourceX
}

func PartOne(lines []string) int {
	m, sourceX := loadMap(lines, false)
	drawMap(m)
	fmt.Println()
	fmt.Println()

	count := -1
	isOffMap := false
	for count = -1; !isOffMap ; count++ {
		isOffMap = dropSand(m, sourceX)
	}

	drawMap(m)
	return count
}

func PartTwo(lines []string) int {
	m, sourceX := loadMap(lines, true)
	drawMap(m)

	var count int
	for count = 0; m[0][sourceX] != Sand; count++ {
		sourceX = dropSandPartTwo(m, sourceX)
	}
	//drawMap(m)

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
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
