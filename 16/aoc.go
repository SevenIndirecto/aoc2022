package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type valve struct {
	id string
	rate int
	open bool
	openedAt int
	tunnels []string
	distancesToValves map[string]int
}

func loadVales(lines []string) map[string]valve {
	valves := make(map[string]valve, 0)

	for _, l := range lines {
		v := valve{open: false, openedAt: -1, tunnels: make([]string, 0)}

		idSplit := strings.Split(l, " has ")
		idSplit = strings.Split(idSplit[0], "Valve ")
		v.id = idSplit[1]

		rateSplit := strings.Split(l, ";")
		rateSplit = strings.Split(rateSplit[0], "=")
		rate, _ := strconv.Atoi(rateSplit[1])
		v.rate = rate

		split := strings.Split(l, "to valve")
		valveStr := split[1][1:]
		if string(split[1][0]) == "s" {
			valveStr = split[1][2:]
		}
		for _, id := range strings.Split(valveStr, ", ") {
			v.tunnels = append(v.tunnels, id)
		}

		valves[v.id] = v
	}
	return valves
}

func copyMap(m map[string]valve) map[string]valve {
	newMap := make(map[string]valve)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

func getScore(state map[string]valve, maxMinutes int) int {
	score := 0
	for _, v := range state {
		if v.open {
			score += (maxMinutes - v.openedAt) * v.rate
		}
	}
	return score
}

func allValvesOpen(state map[string]valve) bool {
	for _, v := range state {
		if v.rate > 0 && !v.open {
			return false
		}
	}
	return true
}

type valveValue struct {
	id string
	value int
}

func shortestDistanceBFS(from string, to string, state map[string]valve) int {
	queue := []string{from}
	visited := map[string]bool{from: true}
	distances := map[string]int{from: 0}

	for len(queue) > 0 {
		currentId := queue[0]
		queue = queue[1:]

		for _, candidateId := range state[currentId].tunnels {
			if visited[candidateId] {
				continue
			}
			queue = append(queue, candidateId)
			distances[candidateId] = distances[currentId]+1
			visited[candidateId] = true

			if candidateId == to {
				return distances[candidateId]
			}
		}
	}

	panic("Should not happen")
	return -1
}

// Find shortest paths to all valves with non-zero flow rates
func addShortestPaths(state map[string]valve) map[string]valve {
	newState := copyMap(state)
	for fromId, v := range newState {
		if v.rate < 1 && fromId != "AA" {
			continue
		}
		v.distancesToValves = make(map[string]int, 0)

		for _, targetValve := range newState {
			if targetValve.rate < 1 || targetValve.id == fromId {
				continue
			}
			cost := shortestDistanceBFS(fromId, targetValve.id, state)
			v.distancesToValves[targetValve.id] = cost
		}
		newState[fromId] = v
	}

	return newState
}

func getValveValuesIfOpened(currentMinute int, currentLoc string, state map[string]valve, maxMinutes int) []valveValue {
	valveValuesIfOpened := make([]valveValue, 0)
	for _, v := range state {
		if v.open || v.rate < 1 {
			continue
		}

		minuteCostToMove := state[currentLoc].distancesToValves[v.id]
		valueAdded := (maxMinutes - currentMinute - minuteCostToMove - 1) * v.rate
		if valueAdded > 0 {
			valveValuesIfOpened = append(valveValuesIfOpened, valveValue{id: v.id, value: valueAdded})
		}
	}
	// Sort descending
	sort.Slice(valveValuesIfOpened, func(i int, j int) bool {
		return valveValuesIfOpened[i].value > valveValuesIfOpened[j].value
	})

	return valveValuesIfOpened
}

func findBestScore(currentMinute int, currentLoc string, state map[string]valve, currentBestScore int, maxMinutes int) int {
	if currentMinute > maxMinutes || allValvesOpen(state) {
		return getScore(state, maxMinutes)
	}

	valveValuesIfOpened := getValveValuesIfOpened(currentMinute, currentLoc, state, maxMinutes)

	if len(valveValuesIfOpened) < 1 {
		return getScore(state, maxMinutes)
	}

	bestScore := currentBestScore
	for _, valveToOpen := range valveValuesIfOpened {
		sc := copyMap(state)
		v := sc[valveToOpen.id]
		v.open = true
		v.openedAt = currentMinute + sc[currentLoc].distancesToValves[valveToOpen.id] + 1
		sc[valveToOpen.id] = v

		newScore := findBestScore(v.openedAt, valveToOpen.id, sc, bestScore, maxMinutes)
		if newScore > bestScore {
			bestScore = newScore
		}
	}
	return bestScore
}

type task struct {
	from string
	target string
	openAt int
}

type actor struct {
	loc string
	task task
}

func copyActors(actors []actor) []actor {
	newActors := make([]actor, 0)
	for _, a := range actors {
		newActors = append(newActors, a)
	}
	return newActors
}

func scoreWithElephant(currentMinute int, actors []actor, state map[string]valve, currentBestScore int, maxMinutes int) int {
	if currentMinute > maxMinutes || allValvesOpen(state) {
		return getScore(state, maxMinutes)
	}
	bestScore := currentBestScore
	sc := copyMap(state)
	actors = copyActors(actors)

	// 1. Execute tasks if applicable
	for i := range actors {
		t := actors[i].task
		if t.openAt == currentMinute  {
			v := sc[t.target]
			v.open = true
			v.openedAt = currentMinute
			sc[t.target] = v

			// Nullify task
			t.openAt = -1
			actors[i].task = t
			actors[i].loc = t.target
		}
	}

	// 2. Assign new tasks
	actorsThatNeedNewTask := make([]int, 0)
	for i := range actors {
		if actors[i].task.openAt == -1 {
			actorsThatNeedNewTask = append(actorsThatNeedNewTask, i)
		}
	}

	// TODO: Just for two for now... But anyway, this covers both needing new tasks
	if len(actorsThatNeedNewTask) == 2 {
		a1 := actors[0]
		a2 := actors[1]
		targetsForFirst := getValveValuesIfOpened(currentMinute, a1.loc, sc, maxMinutes)
		targetsForSecond := getValveValuesIfOpened(currentMinute, a2.loc, sc, maxMinutes)

		if len(targetsForFirst) == 0 && len(targetsForSecond) == 0 {
			return getScore(sc, maxMinutes)
		}
		// Let's see how often we come here
		//fmt.Println("Problematic land...", targetsForFirst, a1, targetsForSecond, a2)

		if len(targetsForFirst) < 1 {
			actorsThatNeedNewTask = []int{1}
		} else if len(targetsForSecond) < 1 {
			actorsThatNeedNewTask = []int{0}
		} else {
			// Both have at least one target
			for _, target1 := range targetsForFirst {
				for _, target2 := range targetsForSecond {
					xt1 := a1.task
					xt1.from = a1.loc
					xt1.target = target1.id
					xt1.openAt = currentMinute + sc[a1.loc].distancesToValves[target1.id] + 1
					actors[0].task = xt1

					xt2 := a2.task
					xt2.from = a2.loc
					xt2.target = target2.id
					xt2.openAt = currentMinute + sc[a2.loc].distancesToValves[target2.id] + 1
					actors[1].task = xt2

					// Process newly assigned tasks
					targetMinute := maxMinutes + 1
					for _, ax := range actors {
						if ax.task.openAt != -1 && ax.task.openAt < targetMinute {
							targetMinute = ax.task.openAt
						}
					}
					newScore := scoreWithElephant(targetMinute, actors, sc, bestScore, maxMinutes)
					if newScore > bestScore {
						fmt.Println("Best score", bestScore)
						bestScore = newScore
					}
				}
			}

			//for _, target2 := range targetsForSecond {
			//	for _, target1 := range targetsForFirst {
			//		xt1 := a1.task
			//		xt1.from = a1.loc
			//		xt1.target = target1.id
			//		xt1.openAt = currentMinute + sc[a1.loc].distancesToValves[target1.id] + 1
			//		actors[0].task = xt1
			//
			//		xt2 := a2.task
			//		xt2.from = a2.loc
			//		xt2.target = target2.id
			//		xt2.openAt = currentMinute + sc[a2.loc].distancesToValves[target2.id] + 1
			//		actors[1].task = xt2
			//
			//		// Process newly assigned tasks
			//		targetMinute := maxMinutes + 1
			//		for _, ax := range actors {
			//			if ax.task.openAt != -1 && ax.task.openAt < targetMinute {
			//				targetMinute = ax.task.openAt
			//			}
			//		}
			//		newScore := scoreWithElephant(targetMinute, actors, sc, bestScore, maxMinutes)
			//		if newScore > bestScore {
			//			bestScore = newScore
			//		}
			//	}
			//}
		}
	}

	if len(actorsThatNeedNewTask) == 1 {
		// Only assign task to one
		aIndex := actorsThatNeedNewTask[0]
		a := actors[aIndex]
		orderedValueTargets := getValveValuesIfOpened(currentMinute, a.loc, sc, maxMinutes)

		// No possible task to assign to current actor
		if len(orderedValueTargets) < 1 {
			// Should process other actor
			targetMinute := maxMinutes + 1
			for _, ax := range actors {
				if ax.task.openAt != -1 && ax.task.openAt < targetMinute {
					targetMinute = ax.task.openAt
				}
			}
			newScore := scoreWithElephant(targetMinute, actors, sc, bestScore, maxMinutes)
			if newScore > bestScore {
				bestScore = newScore
			}
			fmt.Println("2) Best score", bestScore)
			return bestScore
		}

		fmt.Println("Starting-------------------------------------")
		// Assign new tasks to actor
		for _, valveToOpen := range orderedValueTargets {
			fmt.Println("Ordered value targets", orderedValueTargets)
			fmt.Println("Valve to open", valveToOpen)
			actors = copyActors(actors)
			otherActorHasSameTarget := false
			otherActorOpensAt := -1

			for _, ax := range actors {
				if ax.task.openAt != -1 && ax.task.target == valveToOpen.id {
					otherActorHasSameTarget = true
					otherActorOpensAt = ax.task.openAt
					fmt.Println("hi", ax, ax.task.openAt, otherActorOpensAt)
					break
				}
			}
			// In theory... could check if it's cheaper for actor that needs to do task
			// to do it instead, but shouldn't be the case due to ordering...
			newOpenAt := currentMinute + sc[a.loc].distancesToValves[valveToOpen.id] + 1
			fmt.Println("yoo", otherActorOpensAt, "New Open At", newOpenAt, otherActorHasSameTarget)
			if otherActorHasSameTarget && newOpenAt >= otherActorOpensAt {
			//if otherActorHasSameTarget && newOpenAt >= otherActorOpensAt || !otherActorHasSameTarget {
				fmt.Println("continue")
				continue
			} else if otherActorHasSameTarget {
				// This does nothing lul... and wrong due to timing abuse

				fmt.Println("Actors that need new task", actorsThatNeedNewTask)
				fmt.Println(actors, "activeActor", a, newOpenAt, "otherOpensAt", otherActorOpensAt)
				fmt.Println(sc)
				fmt.Println(valveToOpen)

				otherActorIndex := 0
				if aIndex == 0 {
					otherActorIndex = 1
				}
				otherActor := actors[otherActorIndex]
				t := otherActor.task
				t.openAt = -1
				directions := getValveValuesIfOpened(currentMinute, otherActor.loc, sc, maxMinutes)

				for _, d := range directions {
					if d.id == valveToOpen.id {
						continue
					}

					t.openAt = currentMinute + sc[otherActor.loc].distancesToValves[d.id] + 1
					t.target = d.id
					actors[otherActorIndex].task = t
					break
				}
			}

			t := a.task
			t.target = valveToOpen.id
			t.openAt = currentMinute + sc[a.loc].distancesToValves[valveToOpen.id] + 1
			a.task = t
			actors[aIndex] = a

			// Find minute to jump to
			targetMinute := maxMinutes + 1
			for _, ax := range actors {
				if ax.task.openAt != -1 && ax.task.openAt < targetMinute {
					targetMinute = ax.task.openAt
				}
			}

			newScore := scoreWithElephant(targetMinute, actors, sc, bestScore, maxMinutes)
			if newScore > bestScore {
				fmt.Println("Best score", bestScore)
				bestScore = newScore
			}
		}
	}
	return bestScore
}

func PartOne(lines []string) int {
	valves := loadVales(lines)
	valves = addShortestPaths(valves)
	return findBestScore(0, "AA", valves, 0, 30)
}

func PartTwo(lines []string) int {
	valves := loadVales(lines)
	valves = addShortestPaths(valves)
	actors := []actor{
		{loc: "AA", task: task{openAt: -1}},
		{loc: "AA", task: task{openAt: -1}},
	}
	return scoreWithElephant(0, actors, valves, 0, 26)
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
