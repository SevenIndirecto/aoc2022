package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	Min int = 0
	Max int = 1
)

type point struct {
	x int
	y int
}

type sensor struct {
	loc point
	beacon point
	radius int
}

func LoadData(lines []string) []sensor {
	sensors := make([]sensor, 0)

	for _, l := range lines {
		sens := sensor{radius: -1}

		s := strings.Split(l, ":")
		s = strings.Split(s[0], "Sensor at ")
		s = strings.Split(s[1], ", ")
		xSplit := strings.Split(s[0], "=")
		x, _ := strconv.Atoi(xSplit[1])
		ySplit := strings.Split(s[1], "=")
		y, _ := strconv.Atoi(ySplit[1])
		sens.loc = point{x: x, y: y}

		s = strings.Split(l, ":")
		s = strings.Split(s[1], " closest beacon is at ")
		s = strings.Split(s[1], ", ")
		xSplit = strings.Split(s[0], "=")
		x, _ = strconv.Atoi(xSplit[1])
		ySplit = strings.Split(s[1], "=")
		y, _ = strconv.Atoi(ySplit[1])
		sens.beacon = point{x: x, y: y}
		sens.radius = distance(sens.loc, sens.beacon)

		sensors = append(sensors, sens)
	}
	return sensors
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// manhattan distance
func distance(a point, b point) int {
	return abs(a.x - b.x) + abs(a.y - b.y)
}

func canPointContainBeacon(p point, sensors []sensor) bool {
	for _, s := range sensors {
		if distance(p, s.loc) <= s.radius {
			// Only the current sensor's beacon could be within it's range
			// That's the only point where a beacon can be, otherwise we're sure it's not
			return s.beacon.x == p.x && s.beacon.y == p.y
		}
	}
	return true
}

func pointContainsUnknownBeacon(p point, sensors []sensor) bool {
	for _, s := range sensors {
		if distance(p, s.loc) <= s.radius {
			return false
		}
	}
	return true
}

func mergeOverlappingAreas(areas [][2]int) [][2]int {
	mergedAreas := make([][2]int, 0)
	for i := range areas {
		currentArea := &areas[i]
		for j := i+1; j < len(areas); j++ {
			// Try to merge current area into other area
			nextArea := &areas[j]

			// c: [---------]
			// n:     [-------]
			//
			// c: [-----------]
			// n:     [-------]
			//
			if currentArea[Max] <= nextArea[Max] && currentArea[Max] >= nextArea[Min]-1 {
				// Can merge current into next
				if currentArea[Min] < nextArea[Min] {
					nextArea[Min] = currentArea[Min]
				}

				// c: [---------]
				// n: [-------]
				//
				// c:       [---------]
				// n: [-------]
			} else if currentArea[Min] >= nextArea[Min] && currentArea[Min] <= nextArea[Max]+1 {
				if currentArea[Max] > nextArea[Max] {
					nextArea[Max] = currentArea[Max]
				}

				// c:  [-----------]
				// n:    [-------]
			} else if currentArea[Min] <= nextArea[Min] && currentArea[Max] >= nextArea[Max] {
				nextArea[Min] = currentArea[Min]
				nextArea[Max] = currentArea[Max]

				// c:    [-------]
				// n:  [-----------]
			} else if currentArea[Min] >= nextArea[Min] && currentArea[Max] <= nextArea[Max] {
				continue // this area is automatically included in next area
			} else {
				mergedAreas = append(mergedAreas, *currentArea)
				break
			}
		}

		if i >= len(areas)-1 {
			// Iterating last area, not covered by for j loop
			mergedAreas = append(mergedAreas, *currentArea)
		}
	}
	return mergedAreas
}

func FindBeacon(maxX int, maxY int, sensors []sensor) int {
	for y := 0; y <= maxY; y++ {
		areas := make([][2]int, 0)

		for _, s := range sensors {
			sMinY := s.loc.y - s.radius
			sMaxY := s.loc.y + s.radius
			if y < sMinY || y > sMaxY {
				// Sensor does not cover this y
				continue
			}

			newAreaMin := s.loc.x - s.radius + abs(s.loc.y - y)
			newAreaMax := s.loc.x + s.radius - abs(s.loc.y - y)
			if newAreaMin < 0 {
				newAreaMin = 0
			}
			if newAreaMax > maxX {
				newAreaMax = maxX
			}

			if len(areas) < 1 {
				areas = append(areas, [2]int{newAreaMin, newAreaMax})
			} else {
				// See if we can expand any area using this
				merged := false
				for i := range areas {
					if newAreaMin >= areas[i][Min] && newAreaMax <= areas[i][Max] {
						// New area completely absorbed, ignore
						merged = true
						//fmt.Println("Absorbed, break")
						break
					}

					if newAreaMax >= areas[i][Min]-1 && newAreaMin <= areas[i][Min] {
						// Could potentially update to new min
						areas[i][Min] = newAreaMin
						merged = true
					}

					if newAreaMin <= areas[i][Max]+1 && newAreaMax >= areas[i][Max] {
						areas[i][Max] = newAreaMax
						merged = true
					}
				}

				// Could not append into any existing area, new entry
				if !merged {
					areas = append(areas, [2]int{newAreaMin, newAreaMax})
				}

				// Did processing this new area cause any existing areas to overlap?
				areas = mergeOverlappingAreas(areas)
			}
		}
		if len(areas) > 1 {
			// Areas were not merged, so there's an X missing - it should be length 2
			if len(areas) != 2 {
				fmt.Println("y", y, areas)
				panic("Logic failure")
			}
			// Overkill to loop, but ..
			for x := 0; x <= maxX; x++ {
				foundAreaForX := false
				for _, a := range areas {
					if x < a[Min] || x > a[Max] {
						continue
					} else {
						foundAreaForX = true
					}
				}
				if !foundAreaForX {
					//fmt.Println("X=", x)
					return x * 4000000 + y
				}
			}
			return -2
		}
	}
	return -1
}

func getMinMaxX(sensors []sensor) (minX int, maxX int) {
	minX = sensors[0].loc.x
	maxX = minX

	for _, s := range sensors {
		if s.loc.x - s.radius < minX {
			minX = s.loc.x - s.radius
		}
		if s.loc.x + s.radius > maxX {
			maxX = s.loc.x + s.radius
		}
	}
	return minX, maxX
}

func GetCoverCountInRow(y int, sensors []sensor) int {
	minX, maxX := getMinMaxX(sensors)

	count := 0
	for x := minX; x <= maxX; x++ {
		if !canPointContainBeacon(point{x: x, y: y}, sensors) {
			count++
		}
	}
	return count
}

func PartOne(lines []string) int {
	sensors := LoadData(lines)
	t1 := time.Now()
	cc := GetCoverCountInRow(2000000, sensors)
	fmt.Printf("Time for part one %v\n", time.Now().Sub(t1))
	return cc
}

func PartTwo(lines []string) int {
	sensors := LoadData(lines)
	return FindBeacon(4000000, 4000000, sensors)
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
