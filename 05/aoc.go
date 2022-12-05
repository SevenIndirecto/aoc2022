package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type stack []string

type move struct {
	amount int
	from   int
	to     int
}

func PartOne(lines []string) string {
	stacks, moves := LoadSupplyStacksAndMoves(lines)

	for _, move := range moves {
		processMove9000(stacks, move)
	}

	output := ""
	for _, s := range stacks {
		output += s[len(s)-1]
	}
	return output
}

func processMove9000(stacks []stack, m move) []stack {
	var cargo string

	for i := 0; i < m.amount; i++ {
		from := m.from - 1
		to := m.to - 1
		// pop
		cargo, stacks[from] = stacks[from][len(stacks[from])-1], stacks[from][:len(stacks[from])-1]
		stacks[to] = append(stacks[to], cargo)
	}
	return stacks
}

func PartTwo(lines []string) string {
	stacks, moves := LoadSupplyStacksAndMoves(lines)

	for _, move := range moves {
		processMove9001(stacks, move)
	}

	output := ""
	for _, s := range stacks {
		output += s[len(s)-1]
	}
	return output
}

func processMove9001(stacks []stack, m move) []stack {
	from := m.from - 1
	to := m.to - 1
	for j := m.amount; j > 0; j-- {
		stacks[to] = append(stacks[to], stacks[from][len(stacks[from])-j])
	}
	stacks[from] = stacks[from][:len(stacks[from])-m.amount]
	return stacks
}

func LoadSupplyStacksAndMoves(lines []string) ([]stack, []move) {
	// find number of stacks
	startLine := 0
	numStacks := 0

	for {
		if string(lines[startLine][1]) == "1" {
			numStacks, _ = strconv.Atoi(string(lines[startLine][len(lines[startLine])-2]))
			break
		}
		startLine++
	}
	// build stacks
	stacks := make([]stack, numStacks)
	fmt.Println(numStacks)

	for i := startLine-1; i >= 0; i-- {
		// Read for each stack
		for s := 0; s < numStacks; s++ {
			cargoPos := s*4 + 1
			if cargoPos >= len(lines[i]) {
				continue
			}
			stackCargo := string(lines[i][cargoPos])
			if stackCargo == " " {
				continue
			}
			stacks[s] = append(stacks[s], stackCargo)
		}
	}

	moves := make([]move, 0)
	// load moves
	for i := startLine + 2; i < len(lines); i++ {
		s := strings.Split(lines[i], " ")
		if len(s) < 6 {
			break
		}
		amount, _ := strconv.Atoi(s[1])
		from, _ := strconv.Atoi(s[3])
		to, _ := strconv.Atoi(s[5])
		moves = append(moves, move{amount: amount, from: from, to: to})
	}

	return stacks, moves
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
