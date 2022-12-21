package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type node struct {
	prev *node
	next *node
	value int
	key string
}

func loadData(lines []string, decryptionKey int) ([]int, map[string]*node) {
	original := make([]int, len(lines))
	mixer := make(map[string]*node)

	var prev *node = nil

	firstKey := ""
	lastKey := ""

	for i, l := range lines {
		value, _ := strconv.Atoi(l)
		value *= decryptionKey
		key := getKey(i, value)
		original[i] = value
		n := node{value: value, key: key}
		if prev != nil {
			n.prev = prev
			prev.next = &n
		}
		prev = &n
		mixer[key] = &n

		if firstKey == "" {
			firstKey = key
		}
		lastKey = key
	}
	mixer[firstKey].prev = mixer[lastKey]
	mixer[lastKey].next = mixer[firstKey]

	return original, mixer
}

func printMixer(original []int, mixer map[string]*node, startAt string) {
	i := 0
	for i = 0; i < len(original); i++ {
		if original[i] == 0 {
			break
		}
	}
	startAt = getKey(i, 0)
	start := mixer[startAt]

	order := make([]int, len(mixer))

	for i := 0; i < len(mixer); i++ {
		order[i] = start.value
		start = start.next
	}
	fmt.Println(order)
}

func getKey(i int, value int) string {
	return fmt.Sprintf("%v-%v", i, value)
}

func decode(original []int, mixer map[string]*node) {
	for i, value := range original {
		key := getKey(i, value)
		moveBy := value
		if moveBy == 0 {
			continue
		}

		moveBy %= len(mixer) - 1

		elementToMove := mixer[key]
		// Connect current neighbors
		elementToMove.prev.next = elementToMove.next
		elementToMove.next.prev = elementToMove.prev

		target := elementToMove
		if moveBy > 0 {
			for i := 0; i < moveBy; i++ {
				target = target.next
			}
		} else {
			for i := moveBy; i < 0; i++ {
				target = target.prev
			}
			// One extra, so we can have the same logic below
			target = target.prev
		}

		elementToMove.prev = target
		elementToMove.next = target.next

		target.next.prev = elementToMove
		target.next = elementToMove
	}
}

func getNumAfterZero(mixer map[string]*node, n int) int {
	n = n % len(mixer)
	// find target
	var target *node = nil
	for _, v := range mixer {
		if v.value == 0 {
			target = v
			break
		}
	}

	if target == nil {
		panic("Logic error")
	}

	for i := 0; i < n; i++ {
		target = target.next
	}
	return target.value
}

func PartOne(lines []string) int {
	original, mixer := loadData(lines, 1)
	decode(original, mixer)
	return getNumAfterZero(mixer, 1000) + getNumAfterZero(mixer, 2000) + getNumAfterZero(mixer, 3000)
}

func PartTwo(lines []string) int {
	original, mixer := loadData(lines, 811589153)
	for i := 0; i < 10; i++ {
		decode(original, mixer)
	}
	return getNumAfterZero(mixer, 1000) + getNumAfterZero(mixer, 2000) + getNumAfterZero(mixer, 3000)
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
