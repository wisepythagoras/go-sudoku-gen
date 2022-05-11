package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/wisepythagoras/go-sudoku-gen/image"
	"github.com/wisepythagoras/go-sudoku-gen/sudoku"
)

func main() {
	curr := time.Now().UnixNano()
	seedPtr := flag.Int64("seed", curr, "The seed; defaults to current unix timestamp")
	simpleOutputPtr := flag.Bool("simple", false, "Shows a board without UTF-8 borders")
	outputPtr := flag.String("output", "", "The output path (@seed for auto naming)")
	saveImgPtr := flag.Bool("save-img", false, "Whether to save the image or not")
	saveSolutionImgPtr := flag.Bool("save-solution-img", false, "Whether to save the image of the solution or not")
	solvePtr := flag.String("solve", "", "A puzzle to solve")
	flag.Parse()

	var err error

	if *solvePtr != "" {
		board, err := sudoku.ParseBoard(*solvePtr)

		if err != nil {
			fmt.Println(err)
		}

		board.Print(true)

		if *saveImgPtr {
			err = createAndSaveImage(board, false)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Saved the printable image of the sudoku puzzle")
			}
		}

		if *saveSolutionImgPtr {
			err = createAndSaveImage(board, true)

			if err != nil {
				fmt.Println(err)
			}
		}

		numOfSolutions := board.CountSolutions()
		board.Solve()
		board.Print(true)

		fmt.Println("Possible solutions:", numOfSolutions)

		return
	}

	fmt.Println("Seed:", *seedPtr)

	board := sudoku.Sudoku{Seed: *seedPtr}
	board.Init()

	start := time.Now()
	board.Fill()
	puzzle := board.GeneratePuzzle()

	// Here we measure the time it took to run the sudokugeneration algorithm.
	duration := time.Since(start)

	board.Print(!*simpleOutputPtr)
	puzzle.Print(!*simpleOutputPtr)

	if *outputPtr != "" {
		err = board.Save(*outputPtr)
	}

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

	if *saveImgPtr {
		err = createAndSaveImage(puzzle, false)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Saved the printable image of the sudoku puzzle")
		}
	}

	if *saveSolutionImgPtr {
		err = createAndSaveImage(&board, true)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func createAndSaveImage(puzzle *sudoku.Sudoku, isSolution bool) error {
	img, _ := image.CreateImage(puzzle)
	label := ""

	if isSolution {
		label = "solution-"
	}

	fileName := fmt.Sprintf("sudoku-%s%d.png", label, puzzle.Seed)
	f, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
}
