package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/wisepythagoras/go-sudoku-gen/sudoku"
)

func main() {
	curr := time.Now().UnixNano()
	seedPtr := flag.Int64("seed", curr, "The seed; defaults to current unix timestamp")
	simpleOutputPtr := flag.Bool("simple", false, "Shows a board without UTF-8 borders")
	outputPtr := flag.String("output", "", "The output path (@seed for auto naming)")
	flag.Parse()

	var err error

	sudoku := sudoku.Sudoku{Seed: *seedPtr}

	sudoku.Init()

	start := time.Now()
	sudoku.Fill()
	duration := time.Since(start)

	sudoku.Print(!*simpleOutputPtr)
	sudoku.GeneratePuzzle().Print(!*simpleOutputPtr)

	if *outputPtr != "" {
		err = sudoku.Save(*outputPtr)
	}

	if err != nil {
		fmt.Println(err)
	}

	ms := duration.Milliseconds()

	fmt.Println("Seed:", *seedPtr)
	fmt.Print("Execution time: ")

	if ms > 0 {
		fmt.Printf("%dms\n", ms)
	} else {
		fmt.Printf("0.%dms\n", duration.Microseconds())
	}
}
