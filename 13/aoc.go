package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type pair struct {
	ListA *packet
	ListB *packet
}

func (p *pair) toString() string {
	return "{ " + p.ListA.toString() + " / " + p.ListB.toString() + " }"
}

type packet struct {
	Nodes []packet
	Value int
	OriginalString string
}

func (n *packet) isInteger() bool {
	return n.Value > -1
}

func (n *packet) toString() string {
	if len(n.Nodes) < 1 {
		if n.Value == -1 {
			return "[]"
		}
		return strconv.Itoa(n.Value)
	}

	packets := make([]string, 0)
	for _, n := range n.Nodes {
		packets = append(packets, n.toString())
	}
	return "[" + strings.Join(packets, ",") + "]"
}

func compare(left *packet, right *packet, pairNum int) (isRightOrder bool, shouldCheckNext bool) {
	if left.isInteger() && right.isInteger() {
		if left.Value < right.Value {
			return true, false
		} else if left.Value > right.Value {
			return false, false
		}
		return true, true
	}

	// Convert Integer -> [Integer]
	if left.isInteger() && !right.isInteger() {
		left = &packet{Nodes: []packet{{Value: left.Value}}, Value: -1}
	} else if !left.isInteger() && right.isInteger() {
		right = &packet{Nodes: []packet{{Value: right.Value}}, Value: -1}
	}

	for i := 0; i < len(left.Nodes); i++ {
		// Right list ran out before left list -> INCORRECT Order
		if i >= len(right.Nodes) {
			return false, false
		}

		inOrder, checkNext := compare(&left.Nodes[i], &right.Nodes[i], pairNum)
		if checkNext {
			continue
		}
		return inOrder, false
	}

	// Left list ran out before the right -> CORRECT Order
	if len(left.Nodes) < len(right.Nodes) {
		return true, false
	}

	if len(left.Nodes) != len(right.Nodes) {
		panic("Logic failure!")
	}
	// Could not determine and lists are of same length -> CHECK NEXT
	return true, true
}

func parseList(s string) packet {
	packets := make([]packet, 0)
	newNode := packet{Nodes: packets, Value: -1, OriginalString: s}
	if len(s) == 0 {
		return newNode
	}

	split := make([]string, 0)
	bracketBalance := 0
	itemStart := 0
	// example for s: [1,2],[1,[2]],3
	for i := 0; i < len(s); i++ {
		c := string(s[i])
		if c == "," {
			if bracketBalance == 0 {
				split = append(split, s[itemStart:i])
				itemStart = i+1
			}
		} else if c == "[" {
			bracketBalance++
		} else if c == "]" {
			bracketBalance--
		}

		if i == len(s) - 1 {
			split = append(split, s[itemStart:i+1])
		}
	}

	for _, n := range split {
		// List [...]
		if string(n[0]) == "[" {
			packets = append(packets, parseList(n[1:len(n)-1]))
		} else {
			// num
			v, _ := strconv.Atoi(n)
			packets = append(packets, packet{Value: v})
		}
	}

	newNode.Nodes = packets
	return newNode
}

func loadPairs(lines []string) []pair {
	pairs := make([]pair, 0)

	for i := 0; i < len(lines); i+=3 {
		pair := pair{}
		for j := 0; j < 2; j++ {
			// strip [ ]
			list := parseList(lines[i+j][1:len(lines[i+j])-1])

			if j == 0 {
				pair.ListA = &list
			} else {
				pair.ListB = &list
			}
		}
		pairs = append(pairs, pair)
	}

	return pairs
}

func loadPackets(lines []string) []packet {
	packets := make([]packet, 0)
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			continue
		}
		packets = append(packets, parseList(lines[i][1:len(lines[i])-1]))
	}
	return packets
}

func PartOne(lines []string) int {
	pairs := loadPairs(lines)

	sum := 0
	for i := 0; i < len(pairs); i++ {
		inOrder, _ := compare(pairs[i].ListA, pairs[i].ListB, i+1)
		if inOrder {
			sum += i+1
		}
	}
	return sum
}

func PartTwo(lines []string) int {
	lines = append (lines, "[[2]]")
	lines = append (lines, "[[6]]")
	packets := loadPackets(lines)
	sort.Slice(packets, func(a int, b int) bool {
		isOrdered, _ := compare(&packets[a], &packets[b], 0)
		return isOrdered
	})

	mul := 1
	for i := 0; i < len(packets); i++ {
		if packets[i].OriginalString == "[2]" || packets[i].OriginalString == "[6]" {
			mul *= i+1
		}
	}
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
