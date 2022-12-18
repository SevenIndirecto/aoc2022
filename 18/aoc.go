package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type point struct {
	x, y, z int
}

type grid [][][]bool

func PartOne(lines []string) int {
	g := loadCubes(lines)
	return countUnconnectedSides(g)
}

func PartTwo(lines []string) int {
	return 0
}

func loadCubes(lines []string) grid {
	min := point{-1, -1, -1}
	max := point{-1, -1, -1}
	points := make([]point, 0)

	for _, l := range lines {
		s := strings.Split(l, ",")
		x, _ := strconv.Atoi(s[0])
		y, _ := strconv.Atoi(s[1])
		z, _ := strconv.Atoi(s[2])
		points = append(points, point{x, y, z})

		if min.x == -1 || x < min.x {
			min.x = x
		}
		if max.x == -1 || x > max.x {
			max.x = x
		}
		if min.y == -1 || y < min.y {
			min.y = y
		}
		if max.y == -1 || y > max.y {
			max.y = y
		}
		if min.z == -1 || z < min.z {
			min.z = z
		}
		if max.z == -1 || z > max.z {
			max.z = z
		}
	}

	// TODO: Won't do offsets for now, but can do later if needed
	g := make(grid, 0)
	for x := 0; x <= max.x; x++ {
		xPlane := make([][]bool, 0)
		for y := 0; y <= max.y; y++ {
			yRow := make([]bool, 0)
			for z := 0; z <= max.z; z++ {
				yRow = append(yRow, false)
			}
			xPlane = append(xPlane, yRow)
		}
		g = append(g, xPlane)
	}
	//fmt.Println(g)

	for _, p := range points {
		g[p.x][p.y][p.z] = true
	}
	//fmt.Println(g)
	return g
}

func countUnconnectedSides(g grid) int {
	deltas := []point{
		{0, 0, 1},
		{0, 0, -1},
		{0, 1, 0},
		{0, -1, 0},
		{1, 0, 0},
		{-1, 0, 0},
	}

	count := 0
	for x := 0; x < len(g); x++ {
		for y := 0; y < len(g[x]); y++ {
			for z := 0; z < len(g[x][y]); z++ {
				if !g[x][y][z] {
					continue
				}

				for _, d := range deltas {
					nx := x + d.x
					ny := y + d.y
					nz := z + d.z

					if nx < 0 || ny < 0 || nz < 0 || nx >= len(g) || ny >= len(g[x]) || nz >= len(g[x][y]) {
						// Out of bounds
						count++
					} else {
						// No neighbor
						if !g[nx][ny][nz] {
							count++
						}
					}
				}
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
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
