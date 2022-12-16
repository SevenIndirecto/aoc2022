package main

import "testing"

func TestPartOne(t *testing.T) {
	lines, _ := LoadLines("test.txt")
	sensors := LoadData(lines)

	expected := 26
	got := GetCoverCountInRow(10, sensors)
	if got != expected {
		t.Errorf("Part one failed, got %v, expected %v", got, expected)
	}
}

func TestPartTwo(t *testing.T) {
	lines, _ := LoadLines("test.txt")
	sensors := LoadData(lines)

	expected := 56000011
	got := FindBeacon(20, 20, sensors)
	if got != expected {
		t.Errorf("Test failed, got %v, expected %v", got, expected)
	}
}
