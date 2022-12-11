package main

import "testing"

func TestPartOne(t *testing.T) {
	lines, _ := LoadLines("test.txt")

	expected := 13140
	got := PartOne(lines)
	if got != expected {
		t.Errorf("Part one failed, got %v, expected %v", got, expected)
	}
}
