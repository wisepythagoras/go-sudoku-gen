package main

import (
	"fmt"
)

type Box struct {
	N          uint8
	numbers    []uint8
	numbersMap map[uint8]uint8
}

func (b *Box) Init() {
	b.numbers = make([]uint8, b.N)
	b.numbersMap = make(map[uint8]uint8)

	for i := range b.numbers {
		b.numbers[i] = 0
	}
}

func (b *Box) Has(n uint8) bool {
	if _, ok := b.numbersMap[n]; ok {
		return true
	}

	return false
}

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

func (b *Box) Insert(c, r, n uint8) bool {
	if c >= b.N/3 || r >= b.N/3 || n > b.N {
		return false
	}

	if b.Has(n) {
		return false
	}

	pos := c + r*3
	b.numbers[pos] = n
	b.numbersMap[n] = pos

	return true
}

func (b *Box) InsertPos(pos int, n uint8) bool {
	if pos >= int(b.N) || n > b.N {
		return false
	}

	if n != 0 && b.Has(n) {
		return false
	}

	b.numbers[pos] = n
	b.numbersMap[n] = uint8(pos)

	return true
}

func (b *Box) GetRow(r int) []uint8 {
	if r > 3 {
		return []uint8{}
	}

	return b.numbers[r*3 : r*3+3]
}

func (b *Box) GetCol(c int) []uint8 {
	if c > 3 {
		return []uint8{}
	}

	return []uint8{b.numbers[c], b.numbers[c+3], b.numbers[c+6]}
}

func (b *Box) GetNumbers() []uint8 {
	numbers := make([]uint8, 9)
	copy(numbers, b.numbers)

	return numbers
}

func (b *Box) SetNumbers(numbers []uint8) {
	b.numbers = numbers
}

func (b *Box) Empty() {
	b.Init()
}

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
