package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	Noop int = iota
	Add
)

type cpu struct {
	x int
	cycle int
	pipeline []instruction
	log []int
}

type instruction struct {
	variant int
	param int
	cycles int
}

func (c *cpu) processPipeline() {
	c.log = append(c.log, c.x)
	c.cycle++

	newPipeline := make([]instruction, 0)

	for _, i := range c.pipeline {
		if i.cycles == 2 {
			i.cycles--
			newPipeline = append(newPipeline, i)
			continue
		}

		if i.variant == Add {
			c.x += i.param
		}
	}
	c.pipeline = newPipeline
}

func (c *cpu) addInstruction(i instruction) {
	c.pipeline = append(c.pipeline, i)
}

func executeProgram(lines []string) cpu {
	c := cpu{x: 1, cycle: 1, pipeline: make([]instruction, 0), log: make([]int, 0)}

	for _, l := range lines {
		s := strings.Split(l, " ")
		if len(s) == 1 {
			c.addInstruction(instruction{cycles: 1, variant: Noop})
		} else {
			p, _ := strconv.Atoi(s[1])
			c.addInstruction(instruction{cycles: 2, param: p, variant: Add})
		}

		for {
			if len(c.pipeline) < 1 {
				break
			}
			c.processPipeline()
		}
	}

	return c
}

func PartOne(lines []string) int {
	c := executeProgram(lines)

	sum := 0
	for i := 20; i <= 220; i+=40 {
		sum += i * c.log[i-1]
	}
	return sum
}

func PartTwo(lines []string) int {
	c := executeProgram(lines)
	w := 40
	h := 6

	fmt.Println(c.log[:16])

	for i := 0; i < h * w; i++ {
		x := c.log[i]
		p := i % 40

		if p == x-1 || p == x || p == x+1 {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}

		if (i+1) % w == 0 {
			fmt.Println()
		}
	}
	return 0
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
