package main

import (
	"encoding/json"
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
	var decoded FormattedResult

	err := json.Unmarshal([]byte(board.Result()), &decoded)
	if err != nil {
		t.Errorf("Problems with decoding board.Results")
	}

	if decoded.Score != 6 {
		t.Errorf("Expected score to be %d, but received %d", 6, decoded.Score)
	}

	for i := 0; i < len(comparable); i++ {
		if !contain(decoded.Words, comparable[i]) {
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
