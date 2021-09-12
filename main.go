package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	curr := time.Now().UnixNano()
	seedPtr := flag.Int64("seed", curr, "The seed; defaults to current ts")
	flag.Parse()

	sudoku := Sudoku{
		N:    9,
		Seed: *seedPtr,
	}

	sudoku.Init()
	sudoku.Fill()
	sudoku.Print(true)
	err := sudoku.Save()

	fmt.Println(err)
}
