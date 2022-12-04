package main

import "testing"

func TestPartOne(t *testing.T) {
	lines, _ := LoadLines("aoc04_test.txt")

	got := PartOne(lines)
	if got != 2 {
		t.Errorf("Test failed, got %v, expected 2", got)
	}
}

func TestPartTwo(t *testing.T) {
	lines, _ := LoadLines("aoc04_test.txt")

	got := PartTwo(lines)
	expected := 4
	if got != expected {
		t.Errorf("Test failed, got %v, expected %v", got, expected)
	}
}
