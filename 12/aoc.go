package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type point struct {
	shortestPrev *point
	shortestDistance int
	height int32
	x int
	y int
	isEnd bool
}

func (p *point) visitPoint(from *point, m [][]point, reverse bool) {
	directions := [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	newDistance := from.shortestDistance + 1
	if p.shortestDistance != -1 && newDistance >= p.shortestDistance {
		// Do not visit point, a shorter path already set
		return
	}

	// Mark as visited
	if p.shortestDistance == -1 || newDistance < p.shortestDistance {
		p.shortestDistance = newDistance
		p.shortestPrev = from
	}

	if !reverse && p.isEnd {
		return
	}

	candidates := make([]*point, 0)

	for _, d := range directions {
		nx := p.x + d[0]
		ny := p.y + d[1]

		if from.x == nx && from.y == ny {
			continue
		}

		// Edges
		if nx < 0 || ny < 0 || nx >= len(m[0]) || ny >= len(m) {
			continue
		}

		if reverse {
			// Part 2
			if m[ny][nx].height >= p.height - 1 {
				candidates = append(candidates, &m[ny][nx])
			}
		} else {
			// Part 1
			if m[ny][nx].height <= p.height + 1 {
				candidates = append(candidates, &m[ny][nx])
			}
		}
	}

	sort.Slice(candidates, func(i int, j int) bool {
		return candidates[i].height < candidates[j].height
	})

	for _, c := range candidates {
		c.visitPoint(p, m, reverse)
	}
}

func loadMap(lines []string) ([][]point, *point, *point) {
	m := make([][]point, 0)
	sx := -1
	sy := -1
	ex := -1
	ey := 1

	for y, l := range lines {
		m = append(m, make([]point, len(l)))
		for x, c := range l {
			isEnd := false
			height := c

			if c == 69 {
				isEnd = true
				height = 122
			} else if c == 83 {
				height = 97
			}

			p := point{
				height: height,
				x: x,
				y: y,
				isEnd: isEnd,
				shortestDistance: -1,
			}

			if c == 83 {
				sx = p.x
				sy = p.y
			} else if c == 69 {
				ex = p.x
				ey = p.y
			}
			m[y][x] = p
		}
	}

	return m, &m[sy][sx], &m[ey][ex]
}

func PartOne(lines []string) int {
	m, start, end := loadMap(lines)
	start.visitPoint(start, m, false)
	return end.shortestDistance
}

func PartTwo(lines []string) int {
	m, _, end := loadMap(lines)

	end.visitPoint(end, m, true)
	shortestDistance := -2
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[0]); x++ {
			if m[y][x].height != 97 {
				continue
			}

			if m[y][x].shortestDistance > -1 && (shortestDistance == -2 || m[y][x].shortestDistance < shortestDistance) {
				shortestDistance = m[y][x].shortestDistance
			}
		}
	}
	return shortestDistance
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
