package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type tree struct {
	processed bool
	visibility [4]bool
	height int
	x int
	y int
}

func loadForest(lines []string) [][]*tree {
	forest := make([][]*tree, len(lines))

	for y := 0; y < len(lines); y++ {
		row := make([]*tree, len(lines[y]))
		forest[y] = row

		for x := 0; x < len(lines[y]); x++ {
			height, _ := strconv.Atoi(string(lines[y][x]))
			tree := tree{processed: false, height: height, x: x, y: y}
			row[x] = &tree
		}
	}
	return forest
}

func (t *tree) processTree(forest [][]*tree, lenX int, lenY int) (bool, int) {
	// Check all 4 directions
	directions := [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	scenicScore := 1
	for i, d := range directions {
		nY := t.y
		nX := t.x
		scenicCount := 0

		for {
			nY = nY + d[1]
			nX = nX + d[0]

			scenicCount++
			if nY < 0 || nX < 0 || nY >= lenY || nX >= lenX {
				t.visibility[i] = true
				scenicCount--
				break
			} else {
				// Check until edge or block
				nextTree := forest[nY][nX]
				if nextTree.height >= t.height {
					t.visibility[i] = false
					break
				} else if nextTree.processed && nextTree.visibility[i] {
					// Next tree is visible, smaller and already processed can stop
					t.visibility[i] = true

					// Part 2
					if i == 0 {
						scenicCount += nY
					} else if i == 1 {
						scenicCount += lenX - 1 - nX
					} else if i == 2 {
						scenicCount += lenY - 1 - nY
					} else {
						scenicCount += nX
					}
					break
				}
			}
		}
		scenicScore *= scenicCount
	}

	t.processed = true
	for _, isVisibleSomewhere := range t.visibility {
		if isVisibleSomewhere {
			return true, scenicScore
		}
	}
	return false, scenicScore
}

func PartOneAndTwo(lines []string) (int, int) {
	forest := loadForest(lines)
	lenX := len(forest[0])
	lenY := len(forest)
	count := 0
	maxScore := 0

	for y := 0; y < len(forest); y++ {
		for x := 0; x < len(forest[y]); x++ {
			visible, score := forest[y][x].processTree(forest, lenX, lenY)
			if score > maxScore {
				maxScore = score
			}
			if visible {
				count++
			}
		}
	}
	return count, maxScore
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
	partOne, partTwo := PartOneAndTwo(lines);
	fmt.Printf("Part one %v\n", partOne)
	fmt.Printf("Part two %v\n", partTwo)
}
