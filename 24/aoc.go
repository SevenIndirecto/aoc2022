package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	Empty int = 0
	Up = 1
	Right = 2
	Down = 4
	Left = 8
	Wall = 16
)

type grid [][]int

type point struct {
	x int
	y int
}

type pointInTime struct {
	x, y, shift int
	backToStart bool
	visitedEnd bool
}

func (g grid) copy(zeroFilled bool) grid {
	ng := make(grid, len(g))
	for y := range g {
		row := make([]int, len(g[y]))
		for x, value := range g[y] {
			if zeroFilled && value != Wall {
				row[x] = Empty
			} else {
				row[x] = value
			}
		}
		ng[y] = row
	}
	return ng
}

func moveWinds(g grid) grid {
	ng := g.copy(true)

	deltaMap := map[int]point {
		Up: {0, -1},
		Right: {1, 0},
		Down: {0, 1},
		Left: {-1, 0},
	}

	width := len(g[0])
	height := len(g)

	for y := range g {
		for x, tileContent := range g[y] {
			if tileContent == Empty || tileContent == Wall {
				continue
			}

			for i := 0; i <= 3; i++ {
				blizzardDirection := intPow(2, i)
				if (blizzardDirection & tileContent) != blizzardDirection {
					continue
				}

				nx := x + deltaMap[blizzardDirection].x
				ny := y + deltaMap[blizzardDirection].y

				if nx <= 0 {
					nx = width - 2
				} else if nx >= width - 1 {
					nx = 1
				} else if ny <= 0 {
					ny = height - 2
				} else if ny >= height - 1 {
					ny = 1
				}
				ng[ny][nx] |= blizzardDirection
			}
		}
	}

	return ng
}

func createAllPossibleBlizzardStates(g grid) []grid {
	blizzardStates := make([]grid, 0)

	variants := lcm(len(g)-2, len(g[0])-2)
	ng := g
	for i := 0; i < variants; i++ {
		ng = moveWinds(ng)
		blizzardStates = append(blizzardStates, ng)
	}
	return blizzardStates
}

func useBfs(g grid, loopBack bool) int {
	blizzardStates := createAllPossibleBlizzardStates(g)
	numStates := len(blizzardStates)
	moves := [5]point{
		{0, 1},
		{1, 0},
		{0, 0},
		{-1,0},
		{0, -1},
	}
	startPoint := point{1, 0}
	endPoint := point{len(g[0])-2, len(g)-1}

	start := pointInTime{1, 0, 0, false, false}
	distances := map[pointInTime]int{start: 0}
	queue := []pointInTime{start}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		step := distances[p]+1
		shift := step % numStates
		gridState := blizzardStates[shift]

		for _, delta := range moves {
			nx := p.x + delta.x
			ny := p.y + delta.y

			if ny < 0 || ny >= len(g) || gridState[ny][nx] != Empty {
				continue
			}

			visitedEnd := false
			if loopBack && (p.visitedEnd || (nx == endPoint.x && ny == endPoint.y)) {
				visitedEnd = true
			}

			backToStart := false
			if loopBack && p.backToStart || (p.visitedEnd && nx == startPoint.x && ny == startPoint.y) {
				backToStart = true
			}

			// Have we already visited this point?
			npt := pointInTime{nx, ny, shift, backToStart, visitedEnd}
			_, alreadyVisited := distances[npt]

			if !alreadyVisited {
				queue = append(queue, npt)
				distances[npt] = step
			}
		}
	}

	m := -1
	for key, d := range distances {
		if key.x == len(g[0])-2 && key.y == len(g)-1 {
			if loopBack && (!key.backToStart || !key.visitedEnd) {
				continue
			}
			if m == -1 || d < m {
				m = d
			}
		}
	}

	return m+1
}

func PartOne(lines []string) int {
	g := loadMap(lines)
	return useBfs(g, false)
}

func PartTwo(lines []string) int {
	g := loadMap(lines)
	return useBfs(g, true)
}

func drawGrid(g grid, p point) {
	m := map[int]string {
		0: ".",
		Up: "^",
		Right: ">",
		Down: "v",
		Left: "<",
		Wall: "#",
		15: "4",
		14: "3",
		13: "3",
		12: "2",
		11: "3",
		10: "2",
		9: "2",
		7: "3",
		6: "2",
		5: "2",
		3: "2",
	}

	for y := range g {
		for x := range g[y] {
			if x == p.x && y == p.y {
				fmt.Print("E")
			} else {
				fmt.Print(m[g[y][x]])
			}
		}
		fmt.Println()
	}
}

func loadMap(lines []string) grid {
	g := make(grid, len(lines))
	m := map[rune]int{
		'#': Wall,
		'.': Empty,
		'^': Up,
		'>': Right,
		'v': Down,
		'<': Left,
	}

	for y, l := range lines {
		row := make([]int, len(l))
		for x, c := range l {
			row[x] = m[c]
		}
		g[y] = row
	}
	return g
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

func intPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

func main() {
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
