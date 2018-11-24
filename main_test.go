package main

import (
	"testing"
)

func TestEndToEnd(t *testing.T) {
	dictionary := NewDictionary("./test.txt")
	scheme := [][]string{
		[]string{"C", "M", "E"},
		[]string{"A", "T", "R"},
		[]string{"N", "B", "S"},
	}
	board := NewBoard(scheme, dictionary)

	if board == nil {
		t.Errorf("board should not be empty")
	}

	board.FindWords()
	results := board.results

	comparable := []string{"CAT", "MAN", "TREM", "NAME"}

	for i := 0; i < len(comparable); i++ {
		if !contain(results, comparable[i]) {
			t.Errorf("Result should contain: %s", comparable[i])
		}
	}
}

func contain(results []string, value string) bool {
	for _, result := range results {
		if result == value {
			return true
		}
	}
	return false
}
