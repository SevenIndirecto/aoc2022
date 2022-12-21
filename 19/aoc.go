package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	Ore int = iota
	Clay
	Obsidian
	Geode
)

type cost [3]int

type blueprint [4]cost

type state struct {
	resources [4]int
	robots [4]int
	minute int
}

func newState() state {
	s := state{
		minute: 0,
		resources: [4]int{0, 0, 0, 0},
		robots: [4]int{0, 0, 0, 0},
	}
	s.robots[Ore] = 1
	return s
}

func (s *state) copy() state {
	ns := state{}
	ns.resources = [4]int{}
	ns.robots = [4]int{}
	ns.minute = s.minute
	for i := 0; i < 4; i++ {
		ns.resources = s.resources
		ns.robots = s.robots
	}
	return ns
}

func (s *state) collect() {
	for i := range s.resources {
		s.resources[i] += s.robots[i]
	}
}

func (s *state) buildRobot(robot int, b blueprint) {
	for i, resourceCost := range b[robot] {
		s.resources[i] -= resourceCost
	}
	s.robots[robot]++
}

func (s *state) canBuild(robot int, b blueprint) bool {
	for i, resourceCost := range b[robot] {
		if s.resources[i] < resourceCost {
			return false
		}
	}
	return true
}

func (s *state) score() int {
	return s.resources[Geode]
}

func maximizeGeodes(s state, b blueprint, highScore int, maxMinutes int) int {
	robotsWeCanBuild := make(map[int]bool)

	bestScore := highScore
	s.minute++

	if s.minute > maxMinutes {
		return s.score()
	}

	for i := len(s.robots)-1; i >= 0; i-- {
		robotsWeCanBuild[i] = s.canBuild(i, b)
	}
	s.collect()

	// Branch - Robots
	for robotType := len(s.robots)-1; robotType >= 0; robotType-- {
		if !robotsWeCanBuild[robotType] {
			continue
		}

		ns := s.copy()
		ns.buildRobot(robotType, b)
		newScore := maximizeGeodes(ns, b, bestScore, maxMinutes)
		if newScore > bestScore {
			bestScore = newScore
		}
		if  (robotType == Clay && robotsWeCanBuild[Obsidian]) || (robotType == Obsidian && !(robotsWeCanBuild[Clay] && s.robots[Obsidian] < 1)) || robotType == Geode {
			// Well, these are some heuristics I had to massage to get this done...
			break
		}
	}

	if s.resources[Ore] < 9 && !robotsWeCanBuild[Obsidian] && !robotsWeCanBuild[Geode] {
		newScore := maximizeGeodes(s.copy(), b, bestScore, maxMinutes)
		if newScore > bestScore {
			bestScore = newScore
		}
	}

	return bestScore
}

func loadBlueprints(lines []string) []blueprint {
	bps := make([]blueprint, 0)

	for _, l := range lines {
		b := blueprint{}

		s := strings.Split(l, "Each ore robot costs ")
		s = strings.Split(s[1], " ore.")
		oreRobotCost, _ := strconv.Atoi(s[0])
		b[Ore] = cost{oreRobotCost, 0, 0}

		s = strings.Split(l, "Each clay robot costs ")
		s = strings.Split(s[1], " ore.")
		clayRobotCost, _ := strconv.Atoi(s[0])
		b[Clay] = cost{clayRobotCost, 0, 0}

		s = strings.Split(l, "Each obsidian robot costs ")
		s = strings.Split(s[1], " clay. Each geode")
		s = strings.Split(s[0], " ore and ")
		oreCost, _ := strconv.Atoi(s[0])
		clayCost, _ := strconv.Atoi(s[1])
		b[Obsidian] = cost{oreCost, clayCost, 0}

		s = strings.Split(l, "Each geode robot costs ")
		s = strings.Split(s[1], " obsidian.")
		s = strings.Split(s[0], " ore and ")
		oreCost, _ = strconv.Atoi(s[0])
		obsidianCost, _ := strconv.Atoi(s[1])
		b[Geode] = cost{oreCost, 0, obsidianCost}

		bps = append(bps, b)
	}

	return bps
}

func PartOne(lines []string) int {
	bps := loadBlueprints(lines)
	scores := make([]int, len(bps))

	sum := 0
	for i, b := range bps {
		id := i+1
		s := newState()
		score := maximizeGeodes(s, b, 0, 24)
		scores[i] = id * score
		sum += scores[i]
	}
	fmt.Println(scores)
	return sum
}

func PartTwo(lines []string) int {
	bps := loadBlueprints(lines)
	scores := make([]int, len(bps))

	mul := 1
	for i, b := range bps {
		if i >= 3 {
			break
		}
		s := newState()
		score := maximizeGeodes(s, b, 0, 32)
		scores[i] = score
		mul *= scores[i]
	}
	fmt.Println(scores)
	return mul
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
