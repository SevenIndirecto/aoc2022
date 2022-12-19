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

const (
	Empty int = iota
	Cube
	Colored
)

type grid [][][]int

func PartOne(lines []string) int {
	g := loadCubes(lines)
	return countUnconnectedSides(g, false)
}

func PartTwo(lines []string) int {
	g := loadCubes(lines)
	coloredGrid := paintOutsideCubes(g)
	return countUnconnectedSides(coloredGrid, true)
}

func paintOutsideCubes(g grid) grid {
	xLen := len(g)
	yLen := len(g[0])
	zLen := len(g[0][0])

	// NOTE: Could also ignore starting from different corners
	// and only start from {0, 0, 0}, but this might catch some
	// edge cases
	startPoints := [8]point{
		{0, 0, 0},
		{0, 0, zLen-1},
		{0, yLen-1, 0},
		{0, yLen-1, zLen-1},
		{xLen-1, 0, 0},
		{xLen-1, 0, zLen-1},
		{xLen-1, yLen-1, 0},
		{xLen-1, yLen-1, zLen-1},
	}

	neighbors := []point{
		{0, 0, 1},
		{0, 0, -1},
		{0, 1, 0},
		{0, -1, 0},
		{1, 0, 0},
		{-1, 0, 0},
	}

	// Paint from all corners of the cube grid
	for _, start := range startPoints {
		// Obviously could have these as a [3]array but w/e
		dx := -1
		xStop := -1
		if start.x == 0 {
			dx = 1
			xStop = xLen
		}
		dy := -1
		yStop := -1
		if start.y == 0 {
			dy = 1
			yStop = yLen
		}
		dz := -1
		zStop := -1
		if start.z == 0 {
			dz = 1
			zStop = zLen
		}

		for x := start.x; x != xStop; x += dx {
			for y := start.y; y != yStop; y += dy {
				for z := start.z; z != zStop; z += dz {
					if g[x][y][z] == Colored || g[x][y][z] == Cube {
						continue
					}

					for _, n := range neighbors {
						nx := x + n.x
						ny := y + n.y
						nz := z + n.z

						if nx < 0 || ny < 0 || nz < 0 || nx >= xLen || ny >= yLen || nz >= zLen || g[nx][ny][nz] == Colored {
							g[x][y][z] = Colored
							break
						}
					}
				}
			}
		}
	}

	return g
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

	g := make(grid, 0)
	for x := 0; x <= max.x; x++ {
		xPlane := make([][]int, 0)
		for y := 0; y <= max.y; y++ {
			yRow := make([]int, 0)
			for z := 0; z <= max.z; z++ {
				yRow = append(yRow, Empty)
			}
			xPlane = append(xPlane, yRow)
		}
		g = append(g, xPlane)
	}

	for _, p := range points {
		g[p.x][p.y][p.z] = Cube
	}
	return g
}

func countUnconnectedSides(g grid, ignoreAirPockets bool) int {
	target := Empty
	if ignoreAirPockets {
		target = Colored
	}

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
				if g[x][y][z] == Empty || g[x][y][z] == target {
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
						if g[nx][ny][nz] == target {
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
