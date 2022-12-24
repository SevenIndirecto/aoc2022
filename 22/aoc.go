package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	Right int = iota
	Down
	Left
	Up
)

const (
	Top = Up
	Bottom = Down
)

const (
	Clockwise int = 1
	CounterClockwise int = -1
)

const (
	Wall int = 10
	Empty int = 11
	Void int = 12
)

const (
	Rotate int = iota
	Move
)

type person struct {
	side int
	x int
	y int
	direction int
}

type command struct {
	value int
	op int
}

type transition struct {
	side int
	edge int
	//orient int
}

type logEntry struct {
	orientation int
	p point
}

type side struct {
	id int
	transitions map[int]transition
	grid [][]int
	xOffset int
	yOffset int
	log []logEntry
}

func main() {
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines, false))
}

func copyGridSection(startY int, endY int, startX int, endX int, grid [][]int, heightCount int) [][]int {
	size := len(grid) / heightCount
	newGrid := make([][]int, size)

	for y, yNewGrid := startY, 0; y < endY; y, yNewGrid = y+1, yNewGrid+1 {
		newGrid[yNewGrid] = make([]int, size)

		for x, xNewGrid := startX, 0; x < endX; x, xNewGrid = x+1, xNewGrid+1 {
			newGrid[yNewGrid][xNewGrid] = grid[y][x]
		}
	}
	return newGrid
}

// TODO: Only works for my input...
func initSides(grid [][]int) map[int]side {
	m := make(map[int]side)
	size := len(grid) / 4

	m[1] = side{
		id: 1,
		transitions: map[int]transition{
			Top: {6, Left},
			Right: {2, Left},
			Bottom: {3, Top},
			Left: {4, Left},
		},
		grid: copyGridSection(0, size, size, 2*size, grid, 4),
		log: make([]logEntry, 0),
		xOffset: size,
		yOffset: 0,
	}

	m[2] = side{
		id: 2,
		transitions: map[int]transition{
			Top: {6, Bottom},
			Right: {5, Right},
			Bottom: {3, Right},
			Left: {1, Right},
		},
		grid: copyGridSection(0, size, 2*size, 3*size, grid, 4),
		log: make([]logEntry, 0),
		xOffset: 2*size,
		yOffset: 0,
	}
	m[3] = side{
		id: 3,
		transitions: map[int]transition{
			Top: {1, Bottom},
			Right: {2, Bottom},
			Bottom: {5, Top},
			Left: {4, Top},
		},
		grid: copyGridSection(size, 2*size, size, 2*size, grid, 4),
		log: make([]logEntry, 0),
		xOffset: size,
		yOffset: size,
	}
	m[4] = side{
		id: 4,
		transitions: map[int]transition{
			Top: {3, Left},
			Right: {5, Left},
			Bottom: {6, Top},
			Left: {1, Left},
		},
		grid: copyGridSection(2*size, 3*size, 0, size, grid, 4),
		log: make([]logEntry, 0),
		xOffset: 0,
		yOffset: 2*size,
	}
	m[5] = side{
		id: 5,
		transitions: map[int]transition{
			Top: {3, Bottom},
			Right: {2, Right},
			Bottom: {6, Right},
			Left: {4, Right},
		},
		grid: copyGridSection(2*size, 3*size, size, 2*size, grid, 4),
		log: make([]logEntry, 0),
		xOffset: size,
		yOffset: 2 * size,
	}
	m[6] = side{
		id: 6,
		transitions: map[int]transition{
			Top:    {4, Bottom},
			Right:  {5, Bottom},
			Bottom: {2, Top},
			Left:   {1, Top},
		},
		grid: copyGridSection(3*size, 4*size, 0, size, grid, 4),
		log: make([]logEntry, 0),
		xOffset: 0,
		yOffset: 3 * size,
	}

	return m
}

// Derp, assumed all cubes have the same layout as the test example -.-
func initSidesTest(grid [][]int) map[int]side {
	m := make(map[int]side)
	size := len(grid) / 3

	m[1] = side{
		id: 1,
		transitions: map[int]transition{
			Top: {2, Bottom},
			Right: {6, Right},
			Bottom: {4, Top},
			Left: {3, Top},
		},
		grid: copyGridSection(0, size, 2*size, 3*size, grid, 3),
		xOffset: 2 * size,
		yOffset: 0,
	}

	m[2] = side{
		id: 2,
		transitions: map[int]transition{
			Top: {1, Top},
			Right: {3, Left},
			Bottom: {5, Bottom},
			Left: {6, Bottom},
		},
		grid: copyGridSection(size, 2*size, 0, size, grid, 3),
		xOffset: 0,
		yOffset: size,
	}
	m[3] = side{
		id: 3,
		transitions: map[int]transition{
			Top: {1, Left},
			Right: {4, Left},
			Bottom: {5, Left},
			Left: {2, Right},
		},
		grid: copyGridSection(size, 2*size, size, 2*size, grid, 3),
		xOffset: size,
		yOffset: size,
	}
	m[4] = side{
		id: 4,
		transitions: map[int]transition{
			Top: {1, Bottom},
			Right: {6, Top},
			Bottom: {5, Top},
			Left: {3, Right},
		},
		grid: copyGridSection(size, 2*size, 2*size, 3*size, grid, 3),
		xOffset: 2 * size,
		yOffset: size,
	}
	m[5] = side{
		id: 5,
		transitions: map[int]transition{
			Top: {4, Bottom},
			Right: {6, Left},
			Bottom: {2, Bottom},
			Left: {3, Bottom},
		},
		grid: copyGridSection(2*size, 3*size, 2*size, 3*size, grid, 3),
		xOffset: 2 * size,
		yOffset: 2 * size,
	}
	m[6] = side{
		id: 6,
		transitions: map[int]transition{
			Top:    {4, Right},
			Right:  {1, Right},
			Bottom: {2, Left},
			Left:   {5, Right},
		},
		grid: copyGridSection(2*size, 3*size, 3*size, 4*size, grid, 3),
		xOffset: 3 * size,
		yOffset: 2 * size,
	}

	return m
}

func copySidesToMap(sides map[int]side) [][]int {
	size := len(sides[1].grid)
	height := 0
	width := 0
	for _, s := range sides {
		cy := s.yOffset + size
		if cy > height {
			height = cy
		}
		cx := s.xOffset + size
		if cx > width {
			width = cx
		}
	}
	m := make([][]int, 0)

	for y := 0; y < height; y++ {
		row := make([]int, width)
		for x := range row {
			row[x] = Void
		}
		m = append(m, row)
	}

	for _, s := range sides {

		for y, row := range s.grid {
			for x, value := range row {
				realY := s.yOffset + y
				realX := s.xOffset + x

				m[realY][realX] = value
			}
		}

		for _, l := range s.log {
			realY := s.yOffset + l.p.y
			realX := s.xOffset + l.p.x

			m[realY][realX] = l.orientation
		}
	}

	return m
}


func drawMap(m [][]int, p person, sides map[int]side) {
	symbol := map[int]string{
		Void: " ",
		Empty: ".",
		Wall: "#",
	}

	directionMap := map[int]string{
		Right: ">",
		Down: "v",
		Left: "<",
		Up: "^",
	}

	personLoc := point{
		x: sides[p.side].xOffset + p.x,
		y: sides[p.side].yOffset + p.y,
	}
	for y, row := range m {
		for x, tile := range row {
			if personLoc.x == x && personLoc.y == y {
				fmt.Printf("O")
			} else if tile < 10 {
				// TODO: This breaks part one...
				fmt.Printf(directionMap[tile])
			} else {
				fmt.Printf(symbol[tile])
			}
		}
		fmt.Println()
	}
}

func (p *person) rotate(direction int) {
	p.direction += direction
	if direction == Clockwise {
		p.direction %= 4
	} else if p.direction == -1 {
		p.direction = Up
	}
}

type point struct {
	x int
	y int
}

func mirrorY(p point, size int) point {
	return point{p.x, size - 1 - p.y}
}

func mirrorX(p point, size int) point {
	return point{size - 1 - p.x, p.y}
}

func (p *person) cubeMove(sides map[int]side, distance int) {
	deltas := map[int][2]int{
		Right: {1, 0},
		Down: {0, 1},
		Left: {-1, 0},
		Up: {0, -1},
	}
	size := len(sides[p.side].grid)

	tm := make(map[[2]int]func(p point, size int)point)
	tm[[2]int{Top, Bottom}] = func(p point, size int) point {
		return point{p.x,size-1}
	}
	tm[[2]int{Right, Right}] = func(p point, size int) point {
		return mirrorY(point{x: size-1, y: p.y}, size)
	}
	tm[[2]int{Bottom, Top}] = func(p point, size int) point {
		return point{p.x,0}
	}
	tm[[2]int{Left, Top}] = func(p point, size int) point {
		return point{x: p.y, y: 0}
	}
	tm[[2]int{Top, Top}] = func(p point, size int) point {
		return mirrorY(point{x: p.x, y: 0}, size)
	}
	tm[[2]int{Right, Left}] = func(p point, size int) point {
		return point{x: 0, y: p.y}
	}
	tm[[2]int{Bottom, Bottom}] = func(p point, size int) point {
		return mirrorX(point{x: p.x, y: size-1}, size)
	}
	tm[[2]int{Left, Bottom}] = func(p point, size int) point {
		return mirrorX(point{x: p.y, y: size-1}, size)
	}
	tm[[2]int{Top, Left}] = func(p point, size int) point {
		return point{x: 0, y: p.x}
	}
	tm[[2]int{Bottom, Left}] = func(p point, size int) point {
		return mirrorY(point{x: 0, y: p.x}, size)
	}
	tm[[2]int{Left, Right}] = func(p point, size int) point {
		return point{x: size-1, y: p.y}
	}
	tm[[2]int{Right, Top}] = func(p point, size int) point {
		return mirrorX(point{x: p.y, y: 0}, size)
	}
	tm[[2]int{Top, Right}] = func(p point, size int) point {
		return mirrorY(point{x: size-1, y: p.x}, size)
	}
	tm[[2]int{Left, Left}] = func(p point, size int) point {
		return mirrorY(point{x: 0, y: p.y}, size)
	}
	tm[[2]int{Bottom, Right}] = func(p point, size int) point {
		return point{x: size-1, y: p.x}
	}
	tm[[2]int{Right, Bottom}] = func(p point, size int) point {
		return point{x: p.y, y: size-1}
	}

	x := p.x
	y := p.y

	for i := 0; i < distance; i++ {
		dx := deltas[p.direction][0]
		dy := deltas[p.direction][1]

		sideGrid := sides[p.side].grid
		nx := x + dx
		ny := y + dy

		// Handle cube transitions
		transitionSide := -1
		if ny < 0 {
			transitionSide = Top
		} else if nx < 0 {
			transitionSide = Left
		} else if ny >= len(sideGrid) {
			transitionSide = Bottom
		} else if nx >= len(sideGrid[ny]) {
			transitionSide = Right
		}

		if transitionSide != -1 {
			trans := sides[p.side].transitions[transitionSide]

			newEdge := trans.edge
			newOrientation := (trans.edge + 2) % 4
			newSide := trans.side
			newPoint := tm[[2]int{transitionSide, newEdge}](point{nx, ny}, size)

			if sides[newSide].grid[newPoint.y][newPoint.x] == Wall {
				// Stay at current cube and end movement
				break
			} else {
				// Moving to new cube
				p.direction = newOrientation
				p.side = newSide
				nx = newPoint.x
				ny = newPoint.y
			}
		}

		if sides[p.side].grid[ny][nx] == Wall {
			// Stay where we were before
			break
		}
		x = nx
		y = ny

		s := sides[p.side]
		s.log = append(sides[p.side].log, logEntry{
			orientation: p.direction,
			p:           point{x, y},
		})
		sides[p.side] = s
	}

	p.x = x
	p.y = y
}

func (p *person) move(m [][]int, distance int) {
	deltas := map[int][2]int{
		Right: {1, 0},
		Down: {0, 1},
		Left: {-1, 0},
		Up: {0, -1},
	}
	dx := deltas[p.direction][0]
	dy := deltas[p.direction][1]

	x := p.x
	y := p.y
	lastGoodX := x
	lastGoodY := y
	for i := 0; i < distance; i++ {
		nx := x + dx
		ny := y + dy

		// Handle out of bounds movement
		if ny < 0 {
			// Moving up
			ny = len(m)-1
		} else if nx < 0 {
			// Moving left
			nx = len(m[ny])-1
		} else if ny >= len(m) {
			// Moving down
			ny = 0
		} else if nx >= len(m[ny]) {
			// Moving right
			nx = 0
		}

		// Did we enter the void? Skip until we hit wall or empty
		if m[ny][nx] == Void {
			i--
			x = nx
			y = ny
			continue
		}

		if m[ny][nx] == Wall {
			// Stay where we were before
			break
		}

		// Else we're on an empty space
		if m[ny][nx] != Empty {
			panic("Logic error in movement")
		}

		lastGoodX = nx
		lastGoodY = ny
		x = nx
		y = ny
	}
	p.x = lastGoodX
	p.y = lastGoodY
}

func executeCommands(m [][]int, commands []command, sides map[int]side, isPartTwo bool) int {
	p := person{side: 1, y: 0, x: -1, direction: Right}
	for x, tile := range m[0] {
		if tile == Empty {
			p.x = x
			break
		}
	}

	if isPartTwo {
		p.x = 0
	}

	for _, cmd := range commands {
		if cmd.op == Rotate {
			p.rotate(cmd.value)
		} else {
			if isPartTwo {
				p.cubeMove(sides, cmd.value)
			} else {
				p.move(m, cmd.value)
			}
		}
	}

	//if isPartTwo {
	//	m = copySidesToMap(sides)
	//}
	//drawMap(m, p, sides)

	if isPartTwo {
		y := p.y + sides[p.side].yOffset + 1
		x := p.x + sides[p.side].xOffset + 1
		return 1000 * y + 4 * x + p.direction
	}
	return (p.y+1) * 1000 + (p.x+1) * 4 + p.direction
}

func loadMap(lines []string) [][]int {
	m := make([][]int, 0)

	width := 0
	for _, l := range lines {
		 xLen := len(l)
		 if xLen > width {
		 	width = xLen
		 }
	}

	for _, l := range lines {
		row := make([]int, width)

		for x := range row {
			if x < len(l) {
				if l[x] == ' ' {
					row[x] = Void
				} else if l[x] == '.' {
					row[x] = Empty
				} else if l[x] == '#' {
					row[x] = Wall
				} else {
					panic("Logic error")
				}
			} else {
				row[x] = Void
			}
		}
		m = append(m, row)
	}
	return m
}

func loadCommands(commandStream string) []command {
	cmds := make([]command, 0)

	moveStart := 0
	for i, chr := range commandStream {
		c := string(chr)
		if c == "R" || c == "L" {
			// Process previous move
			dStr := commandStream[moveStart:i]
			distance, _ := strconv.Atoi(dStr)
			cmds = append(cmds, command{op: Move, value: distance})
			moveStart = i+1
		}
		if c == "R" {
			cmds = append(cmds, command{op: Rotate, value: Clockwise})
		} else if c == "L" {
			cmds = append(cmds, command{op: Rotate, value: CounterClockwise})
		}

		if i == len(commandStream)-1 && moveStart <= i {
			dStr := commandStream[moveStart:i+1]
			distance, _ := strconv.Atoi(dStr)
			cmds = append(cmds, command{op: Move, value: distance})
		}
	}

	return cmds
}


func PartOne(lines []string) int {
	cmds := loadCommands(lines[len(lines)-1])
	m := loadMap(lines[0:len(lines)-2])
	return executeCommands(m, cmds, make(map[int]side), false)
}

func PartTwo(lines []string, isTest bool) int {
	cmds := loadCommands(lines[len(lines)-1])
	m := loadMap(lines[0:len(lines)-2])
	var sides map[int]side
	if isTest {
		sides = initSidesTest(m)
	} else {
		sides = initSides(m)
	}
	return executeCommands(m, cmds, sides, true)
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

