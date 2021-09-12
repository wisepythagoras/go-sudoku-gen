package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type Sudoku struct {
	// Figure out the logic here. Ideally we want them to add
	// as many as they want, but currently the logic below is
	// just for 9 (the typical).
	N     uint8  `json:"n"` // The number of columns and rows.
	Seed  int64  `json:"seed"`
	Board []*Box `json:"board"`
}

// Init initializes the Sudoku instance. It's required before running `Fill`.
func (s *Sudoku) Init() {
	s.N = 9
	s.Board = make([]*Box, s.N)

	for i := range s.Board {
		s.Board[i] = &Box{N: s.N}
		s.Board[i].Init()
	}

	rand.Seed(s.Seed)
}

// Fill fills the Sudoku board with numbers.
func (s *Sudoku) Fill() {
	history := make([]int, 0, 8)

	for i := 0; i <= 8; i++ {
		if len(history) > 8 {
			history = history[1:]
		}

		history = append(history, i)
		occurances := 0

		// Count the number of times the current box index has been processed in the past
		// few loops.
		for _, v := range history {
			if v == i {
				occurances++
			}
		}

		// If the current board composition has been lead into an infinite loop, we want
		// to empty the current board and start over.
		if len(history) > 5 && occurances > len(history)/2 {
			for _, box := range s.Board {
				box.Empty()
			}

			i = 0
			s.Fill()
			break
		}

		box := s.Board[i]

		// Empty a box before it's processed.
		s.Board[i].Empty()

		var boxH1 *Box
		var boxH2 *Box
		var boxV1 *Box
		var boxV2 *Box

		// Get the adjacent vertical boxes.
		if i < 3 {
			boxV1 = s.Board[i+3]
			boxV2 = s.Board[i+6]
		} else if i < 6 {
			boxV1 = s.Board[i-3]
			boxV2 = s.Board[i+3]
		} else {
			boxV1 = s.Board[i-3]
			boxV2 = s.Board[i-6]
		}

		// Get the adjacent horizontal boxes.
		if i%3 == 0 {
			boxH1 = s.Board[i+1]
			boxH2 = s.Board[i+2]
		} else if i%3 == 1 {
			boxH1 = s.Board[i-1]
			boxH2 = s.Board[i+1]
		} else {
			boxH1 = s.Board[i-2]
			boxH2 = s.Board[i-1]
		}

		for j := 0; j < 9; j++ {
			row := make([]uint8, 9)

			// Get all the numbes in the current row.
			if j < 3 {
				copy(row, boxH1.GetRow(0))
				row = append(row, boxH2.GetRow(0)...)
			} else if j < 6 {
				copy(row, boxH1.GetRow(1))
				row = append(row, boxH2.GetRow(1)...)
			} else if j < 9 {
				copy(row, boxH1.GetRow(2))
				row = append(row, boxH2.GetRow(2)...)
			}

			// Get all the numbers in the current column.
			col := append(boxV1.GetCol(j%3), boxV2.GetCol(j%3)...)

			possible := make([]uint8, 0)

			// Get the list of all of the numbers that can possibly apply to
			// the current cell.
			for k := uint8(1); k <= 9; k++ {
				isPossible := true

				if box.Has(k) {
					isPossible = false
				} else {
					for _, v := range row {
						if v == k {
							isPossible = false
							break
						}
					}

					for _, v := range col {
						if v == k {
							isPossible = false
							break
						}
					}
				}

				if isPossible {
					possible = append(possible, k)
				}
			}

			// If the possible choices have been exahusted, we go to the previous box
			// and set it up from scratch. There is a possibility that this back-tracking
			// could lead to an infinite loop, but the circuit-breaker will prevent this
			// by resetting the entire board and starting again, without resetting the
			// seed that was used.
			if len(possible) == 0 {
				if i <= 2 {
					i = 0
				} else {
					i -= 2
				}

				break
			}

			// Get a random number from all the possibilities and insert it in the target
			// position of the current box.
			num := possible[rand.Intn(len(possible))]
			box.InsertPos(j, num)
		}
	}
}

// GeneratePuzzle needs to run after `Fill`. It generates a proper puzzle with some
// indecies which are hidden.
func (s *Sudoku) GeneratePuzzle() *Sudoku {
	rand.Seed(s.Seed)

	board := make([][]uint8, 9)
	missing := rand.Intn(59-57) + 57
	maxEmptyPerBox := 8
	minEmptyPerBox := 5

	for i, box := range s.Board {
		board[i] = box.GetNumbers()
	}

	for i := 0; i < 5; i++ {
		opposite := i
		emptyAmount := 0

		if i == 3 {
			opposite = i + 2
		} else {
			opposite = 8 - i
		}

		for j := 0; j < 9; j++ {
			if board[i][j] == 0 {
				emptyAmount++
				continue
			}

			if missing == 0 {
				break
			}

			indexIsEmpty := rand.Intn(2)
			oppositeIndex := 8 - j

			if indexIsEmpty == 1 {
				board[i][j] = 0
				board[opposite][oppositeIndex] = 0
				missing -= 2
				emptyAmount++
			}

			if emptyAmount >= maxEmptyPerBox {
				break
			}
		}

		// If not enough cells are empty, we need to go back and empty more.
		if emptyAmount < minEmptyPerBox {
			i -= 1
			continue
		}

		if missing == 0 {
			break
		}
	}

	puzzle := &Sudoku{
		N:    s.N,
		Seed: s.Seed,
	}
	puzzle.Init()

	for i := 0; i < 9; i++ {
		box := &Box{N: s.N}
		box.Init()
		box.SetNumbers(board[i])
		puzzle.Board[i] = box
	}

	return puzzle
}

// Save creates a JSON file for this board.
func (s *Sudoku) Save(fileName string) error {
	sudokuJson, err := json.Marshal(s)

	if err != nil {
		return err
	}

	if fileName == "@seed" {
		fileName = fmt.Sprintf("sudoku-%d.json", s.Seed)
	}

	// Write the file with 0644 permissions.
	err = os.WriteFile(fileName, sudokuJson, 0644)

	if err != nil {
		return err
	}

	return nil
}

// Print displays the board in stdout.
func (s *Sudoku) Print(showRich bool) {
	if showRich {
		printLine(0)
	}

	for i := 0; i < 3; i++ {
		box1 := s.Board[i*3]
		box2 := s.Board[i*3+1]
		box3 := s.Board[i*3+2]

		for j := 0; j < 3; j++ {
			row := make([]uint8, 3)
			copy(row, box1.GetRow(j))

			row = append(row, box2.GetRow(j)...)
			row = append(row, box3.GetRow(j)...)

			if showRich {
				fmt.Print("\xe2\x95\x91")
			}

			for k, v := range row {
				end := ""

				if showRich {
					end = "\xe2\x94\x82"

					if (k-2)%3 == 0 {
						end = "\xe2\x95\x91"
					}

					fmt.Print(" ")
				}

				if v != 0 {
					fmt.Print(v, " ", end)
				} else {
					fmt.Print("  ", end)
				}
			}

			fmt.Println()

			if showRich && (i != 2 || (i == 2 && j < 2)) {
				printLine(j + 3)
			}
		}
	}

	if showRich {
		printLine(9)
	}
}

func printLine(i int) {
	if i == 0 {
		fmt.Print("\xe2\x95\x94")
	} else if i == 9 {
		fmt.Print("\xe2\x95\x9a")
	} else {
		if (i)%5 == 0 {
			fmt.Print("\xe2\x95\xa0")
		} else {
			fmt.Print("\xe2\x95\x9f")
		}
	}

	for l := 0; l < 35; l++ {
		if l-11 >= 0 && (l-11)%12 == 0 {
			if i == 0 {
				fmt.Print("\xe2\x95\xa6")
			} else if i == 9 {
				fmt.Print("\xe2\x95\xa9")
			} else if (i-2)%3 == 0 {
				fmt.Print("\xe2\x95\xac")
			} else {
				fmt.Print("\xe2\x95\xab")
			}
		} else if (l-3)%4 == 0 {
			if i == 0 {
				fmt.Print("\xe2\x95\xa4")
			} else if i == 9 {
				fmt.Print("\xe2\x95\xa7")
			} else if (i-2)%3 == 0 {
				fmt.Print("\xe2\x95\xaa")
			} else {
				fmt.Print("\xe2\x94\xbc")
			}
		} else if i == 0 || i == 9 || (i-2)%3 == 0 {
			fmt.Print("\xe2\x95\x90")
		} else {
			fmt.Print("\xe2\x94\x80")
		}
	}

	if i == 0 {
		fmt.Print("\xe2\x95\x97")
	} else if i == 9 {
		fmt.Print("\xe2\x95\x9d")
	} else {
		if (i)%5 == 0 {
			fmt.Print("\xe2\x95\xa3")
		} else {
			fmt.Print("\xe2\x95\xa2")
		}
	}

	fmt.Println()
}
