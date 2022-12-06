package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func detectMarker(stream string, length int) int {
	for i := length-1; i < len(stream); i++ {
		set := make(map[uint8]bool)
		for j := length-1; j >= 0; j-- {
			set[stream[i-j]] = true
		}
		if len(set) == length {
			return i+1
		}
	}
	return 0
}

func PartOne(lines []string) int {
	return detectMarker(lines[0], 4)
}

func PartTwo(lines []string) int {
	return detectMarker(lines[0], 14)
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
