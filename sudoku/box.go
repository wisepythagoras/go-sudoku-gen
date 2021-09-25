package sudoku

import (
	"fmt"
)

// Box defines the strucure of each box in a Sudoku puzzle.
type Box struct {
	N          uint8
	numbers    []uint8
	numbersMap map[uint8]uint8
}

// Init initializes the numbers array and numbers map. This needs to be
// run before anything else runs.
func (b *Box) Init() {
	b.numbers = make([]uint8, b.N)
	b.numbersMap = make(map[uint8]uint8)

	for i := range b.numbers {
		b.numbers[i] = 0
	}
}

// Has returns whether this box has a number.
func (b *Box) Has(n uint8) bool {
	if _, ok := b.numbersMap[n]; ok {
		return true
	}

	return false
}

// RowHas returns whether a row has a number.
func (b *Box) RowHas(r, n uint8) bool {
	pos, ok := b.numbersMap[n]

	if !ok {
		return false
	}

	// TODO: This shouldn't be 2.
	if pos-3*r <= 2 {
		return true
	}

	return false
}

// ColHas returns whether a column has a number.
func (b *Box) ColHas(c, n uint8) bool {
	pos, ok := b.numbersMap[n]

	if !ok {
		return false
	}

	if pos-c == 0 || pos-c == 3 || pos-c == 6 {
		return true
	}

	return false
}

// GetPos returns the number in a specific absolute position of the box.
func (b *Box) GetPos(i int) uint8 {
	if i >= int(b.N) || i < 0 {
		return 0
	}

	return b.numbers[i]
}

// Insert places a number in a specific cell.
func (b *Box) Insert(c, r, n uint8) bool {
	if c >= b.N/3 || r >= b.N/3 || n > b.N {
		return false
	}

	if n != 0 && b.Has(n) {
		return false
	}

	pos := c + r*3

	if _, ok := b.numbersMap[b.numbers[pos]]; ok {
		delete(b.numbersMap, b.numbers[pos])
	}

	b.numbers[pos] = n

	if n != 0 {
		b.numbersMap[n] = pos
	}

	return true
}

// InsertPos places a number in a specific position.
func (b *Box) InsertPos(pos int, n uint8) bool {
	if pos >= int(b.N) || n > b.N {
		return false
	}

	if n != 0 && b.Has(n) {
		return false
	}

	if _, ok := b.numbersMap[b.numbers[pos]]; ok {
		delete(b.numbersMap, b.numbers[pos])
	}

	b.numbers[pos] = n
	b.numbersMap[n] = uint8(pos)

	return true
}

// GetRow returns the numbers in a specific row.
func (b *Box) GetRow(r int) []uint8 {
	if r > 3 {
		return []uint8{}
	}

	return b.numbers[r*3 : r*3+3]
}

// GetCol returns the umbers in a specific column.
func (b *Box) GetCol(c int) []uint8 {
	if c > 3 {
		return []uint8{}
	}

	return []uint8{b.numbers[c], b.numbers[c+3], b.numbers[c+6]}
}

// GetNumbers returns a copy of the numbers in the box.
func (b *Box) GetNumbers() []uint8 {
	numbers := make([]uint8, 9)
	copy(numbers, b.numbers)

	return numbers
}

// CountEmpty returns the amount of empty cells.
func (b *Box) CountEmpty() int {
	count := 0

	for _, n := range b.numbers {
		if n == 0 {
			count++
		}
	}

	return count
}

// SetNumbers sets the numbers.
func (b *Box) SetNumbers(numbers []uint8) {
	b.numbers = numbers

	for i, num := range numbers {
		b.numbersMap[num] = uint8(i)
	}
}

func (b *Box) Empty() {
	b.Init()
}

// MarshalJSON is used by the JSON module to tell it how the JSON format looks
// like.
func (b *Box) MarshalJSON() ([]byte, error) {
	array := "["

	for i, v := range b.numbers {
		if i == 0 {
			array = fmt.Sprintf("%s%d", array, v)
		} else {
			array = fmt.Sprintf("%s,%d", array, v)
		}
	}

	array = array + "]"

	return []byte(array), nil
}
