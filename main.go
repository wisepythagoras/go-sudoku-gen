package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		return
	}

	seed, err := strconv.Atoi(args[0])

	if err != nil {
		fmt.Println(err)
		return
	}

	sudoku := Sudoku{
		N:             9,
		MissingDigits: 10,
		Seed:          int64(seed),
	}

	sudoku.Init()
	sudoku.Fill()
	sudoku.Print()
}
