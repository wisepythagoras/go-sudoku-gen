package sudoku

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
	count int64
}

// Init initializes the Sudoku instance. It's required before running `Fill`.
func (s *Sudoku) Init() {
	s.N = 9
	s.count = 0
	s.Board = make([]*Box, s.N)

	for i := range s.Board {
		s.Board[i] = &Box{N: s.N}
		s.Board[i].Init()
	}

	rand.Seed(s.Seed)
}

// Fill fills the Sudoku board with numbers.
func (s *Sudoku) Fill() {
	for i := 0; i <= 8; i++ {
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
				for _, box := range s.Board {
					box.Empty()
				}

				i = 0
				s.Fill()
				return
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
// TODO: Start from scratch.
func (s *Sudoku) GeneratePuzzle() *Sudoku {
	const targetMissing = 58
	const maxEmptyPerBox = 9
	const minEmptyPerBox = 4

	rand.Seed(s.Seed + s.count)

	// This will hold the raw values of our board.
	board := make([][]uint8, 9)

	// Get all the numbers of our board into the array.
	for i, box := range s.Board {
		board[i] = box.GetNumbers()
	}

	// This variable will count the number of cells we empty out.
	totalRemoved := 0

	// In order for a puzzle to be valid, it needs to to have all numbers present, otherwise it's likely
	// a puzzle will be unsolvable.
	numMap := make(map[uint8]int)
	numMap[1] = 9
	numMap[2] = 9
	numMap[3] = 9
	numMap[4] = 9
	numMap[5] = 9
	numMap[6] = 9
	numMap[7] = 9
	numMap[8] = 9
	numMap[9] = 9

	for i, box := range s.Board {
		board[i] = box.GetNumbers()
	}

	for i := 0; i < 4; i++ {
		opposite := 8 - i

		amountToEmpty := rand.Intn(maxEmptyPerBox-minEmptyPerBox) + minEmptyPerBox

		for j := amountToEmpty; j > 0; {
			index := rand.Intn(9)

			if board[i][index] == 0 {
				continue
			}

			available := numMap[board[i][index]]
			oppositeIndex := 8 - index
			oppositeAvailable := numMap[board[opposite][oppositeIndex]]

			if (board[opposite][oppositeIndex] == board[i][index] &&
				available < 2) ||
				oppositeAvailable < 2 {
				continue
			}

			totalRemoved += 2

			if board[opposite][oppositeIndex] == board[i][index] {
				available -= 1
			} else {
				numMap[board[opposite][oppositeIndex]] = oppositeAvailable - 1
			}

			numMap[board[i][index]] = available - 1

			board[i][index] = 0
			board[opposite][oppositeIndex] = 0

			j--
		}
	}

	// For the fifth box (the center one).
	for j := 0; j < 4; j++ {
		shouldEmpty := rand.Intn(4) >= 1

		if !shouldEmpty {
			continue
		}

		available := numMap[board[4][j]]
		oppositeIndex := 8 - j
		oppositeAvailable := numMap[board[4][oppositeIndex]]

		if (board[4][oppositeIndex] == board[4][j] &&
			available < 2) ||
			oppositeAvailable < 2 {
			continue
		}

		if j == 4 {
			totalRemoved++
		} else {
			totalRemoved += 2
		}

		if board[4][oppositeIndex] == board[4][j] {
			available -= 1
		} else {
			numMap[board[4][oppositeIndex]] = oppositeAvailable - 1
		}

		numMap[board[4][j]] = available - 1

		board[4][j] = 0
		board[4][oppositeIndex] = 0
	}

	// if totalRemoved < targetMissing-4 {
	// 	s.count++
	// 	return s.GeneratePuzzle()
	// }

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

	solvedPuzzle := &Sudoku{}
	solvedPuzzle.Copy(puzzle)

	if puzzle.HasMultipleSolutions() {
		s.count++
		return s.GeneratePuzzle()
	}

	return puzzle
}

// IsEqual checks whether two Sudoku boards are equal.
func (s *Sudoku) IsEqual(sudoku *Sudoku) bool {
	for i, box := range s.Board {
		selfNumbers := box.GetNumbers()
		otherNumbers := sudoku.Board[i].GetNumbers()

		for j, num := range selfNumbers {
			if num != otherNumbers[j] {
				return false
			}
		}
	}

	return true
}

// GetCounter returns the number of itterations on the specific board.
func (s *Sudoku) GetCounter() int64 {
	return s.count
}

// GetRow returns all the numbers in a specific row.
func (s *Sudoku) GetRow(row int) []uint8 {
	numbers := make([]uint8, 0)
	boxIdx := 2

	if row < 3 {
		boxIdx = 0
	} else if row < 6 {
		boxIdx = 1
	}

	for i := boxIdx * 3; i < boxIdx*3+3; i++ {
		numbers = append(numbers, s.Board[i].GetRow(row%3)...)
	}

	return numbers
}

// GetCol gets all the numbers in a given column.
func (s *Sudoku) GetCol(col int) []uint8 {
	numbers := make([]uint8, 0)
	boxIdx := 2

	if col < 3 {
		boxIdx = 0
	} else if col < 6 {
		boxIdx = 1
	}

	for i := 0 + boxIdx; i <= 6+boxIdx; i += 3 {
		numbers = append(numbers, s.Board[i].GetCol(col%3)...)
	}

	return numbers
}

// GetBoxFromRowCol returns the box from a specific point in the board.
func (s *Sudoku) GetBoxFromRowCol(row int, col int) *Box {
	boxRow := 2
	boxCol := 2

	if row < 3 {
		boxRow = 0
	} else if row < 6 {
		boxRow = 1
	}

	if col < 3 {
		boxCol = 0
	} else if col < 6 {
		boxCol = 1
	}

	return s.Board[(boxRow*3)+boxCol]
}

// Solve tries to solve the puzzle and returns the first possible solution.
func (s *Sudoku) Solve() bool {
	for i := 0; i < 9; i++ {
		row := s.GetRow(i)

		for j, num := range row {
			if num != 0 {
				continue
			}

			column := s.GetCol(j)
			box := s.GetBoxFromRowCol(i, j)
			possibilites := getVHPossibilities(row, column, box)

			for _, possibility := range possibilites {
				box.Insert(uint8(j%3), uint8(i%3), possibility)

				if s.Solve() {
					return true
				} else {
					box.Insert(uint8(j%3), uint8(i%3), 0)
				}
			}

			return false
		}
	}

	return true
}

// CountEmpty returns the total number of empty cells in the puzzle.
func (s *Sudoku) CountEmpty() int {
	count := 0

	for _, box := range s.Board {
		count += box.CountEmpty()
	}

	return count
}

// CountSolutions returns the total amount of solutions for this board.
func (s *Sudoku) CountSolutions() int64 {
	sudoku := &Sudoku{}
	sudoku.Copy(s)

	return s.internalCount(sudoku, 0, int64(sudoku.CountEmpty()), false)
}

// HasMultipleSolutions returns true if there are multiple solutions, or false if there
// is only one.
func (s *Sudoku) HasMultipleSolutions() bool {
	sudoku := &Sudoku{}
	sudoku.Copy(s)

	return s.internalCount(sudoku, 0, int64(sudoku.CountEmpty()), true) > 1
}

// internalCount will loop through all possible combinations and count the number of solutions it finds.
// By supplying true to `breakIfMultiple` it will break if it finds more than one solutions.
func (s *Sudoku) internalCount(sudoku *Sudoku, count, totalEmpty int64, breakIfMultiple bool) int64 {
	// If we've filled all of the initially empty slots, we increment the counter.
	if totalEmpty == 0 {
		count++
	}

	if breakIfMultiple && count > 1 {
		return count
	}

	for i := 0; i < 9; i++ {
		row := sudoku.GetRow(i)

		for j, num := range row {
			if num != 0 {
				continue
			}

			column := sudoku.GetCol(j)
			box := sudoku.GetBoxFromRowCol(i, j)
			possibilites := getVHPossibilities(row, column, box)

			// Recursively run through all of the possibilities and try to fill in all of
			// the empty cells.
			for _, possibility := range possibilites {
				box.Insert(uint8(j%3), uint8(i%3), possibility)
				count = s.internalCount(sudoku, count, totalEmpty-1, breakIfMultiple)
				box.Insert(uint8(j%3), uint8(i%3), 0)

				if breakIfMultiple && count > 1 {
					return count
				}
			}

			return count
		}
	}

	return count
}

// Copy copies a sudoku board into this instance.
func (s *Sudoku) Copy(board *Sudoku) {
	s.Init()

	for i, box := range board.Board {
		s.Board[i].SetNumbers(box.GetNumbers())
	}

	s.N = board.N
	s.Seed = board.Seed
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

// GetBox returns a box in a specific position.
func (s *Sudoku) GetBox(idx int) *Box {
	if idx >= 9 {
		return nil
	}

	return s.Board[idx]
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

// getVHPossibilities gets the vertical and horizontal possibilities.
func getVHPossibilities(row, col []uint8, box *Box) []uint8 {
	possibilites := make([]uint8, 0)

	for i := 1; i <= 9; i++ {
		if box.Has(uint8(i)) {
			continue
		}

		found := false

		for _, num := range row {
			if num == 0 {
				continue
			}

			if num == uint8(i) {
				found = true
				break
			}
		}

		if found {
			continue
		}

		for _, num := range col {
			if num == 0 {
				continue
			}

			if num == uint8(i) {
				found = true
				break
			}
		}

		if !found {
			possibilites = append(possibilites, uint8(i))
		}
	}

	return possibilites
}
