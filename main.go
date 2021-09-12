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

	start := time.Now()
	sudoku.Fill()
	duration := time.Since(start)

	sudoku.Print(true)
	err := sudoku.Save()

	if err != nil {
		fmt.Println(err)
	}

	ms := duration.Milliseconds()

	fmt.Print("Execution time: ")

	if ms > 0 {
		fmt.Printf("%dms\n", ms)
	} else {
		fmt.Printf("0.%dms\n", duration.Microseconds())
	}
}
