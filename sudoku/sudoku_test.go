package sudoku_test

import (
	"testing"

	"github.com/wisepythagoras/go-sudoku-gen/sudoku"
)

func TestSudokuSolve(t *testing.T) {
	s := sudoku.Sudoku{}
	s.Init()

	var arr [][]uint8 = [][]uint8{
		{7, 0, 2, 0, 0, 0, 1, 0, 0},
		{0, 5, 0, 0, 0, 3, 0, 0, 9},
		{6, 0, 0, 0, 0, 0, 5, 0, 0},
		{8, 0, 0, 0, 4, 3, 0, 9, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 9, 0, 7, 5, 0, 0, 0, 8},
		{0, 0, 9, 0, 0, 0, 0, 0, 7},
		{7, 0, 0, 2, 0, 0, 0, 4, 0},
		{0, 0, 5, 0, 0, 0, 2, 0, 3},
	}

	for i, numbers := range arr {
		s.Board[i].SetNumbers(numbers)
	}

	if !s.Solve() {
		t.Error("Unable to solve sudoku")
	}

	if s.CountSolutions() != 1 {
		t.Error("The board is supposed to have only one solution")
	}
}
