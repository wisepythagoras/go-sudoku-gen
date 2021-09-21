# Sudoku Generator

This little program, which has no external dependencies, is a quite fast procedural [Sudoku](https://en.wikipedia.org/wiki/Sudoku) puzzle generator.

## Usage

```
Usage of ./go-sudoku-gen:
  -output string
        The output path (@seed for auto naming)
  -save-img
        Whether to save the image or not
  -seed int
        The seed; defaults to current unix timestamp (default 1631573683595299425)
  -simple
        Shows a board without UTF-8 borders
```

## How it works

### Generating valid boards

The board is created and filled by randomizing the number that goes in each cell. If a point is reached where it runs out of possible placements, then it attempts to re-fill the board from scratch.

### Generating a puzzle

In order to generate a valid puzzle, the algorithm randomly chooses which cells to empty. At the end, it will verify that there is only one possible solution, otherwise it will attempt to re-generate a puzzle.

## Sample output

``` sh
./go-sudoku-gen -seed 1034875
```

```
Seed: 1034875
╔═══╤═══╤═══╦═══╤═══╤═══╦═══╤═══╤═══╗
║ 8 │ 6 │ 1 ║ 4 │ 5 │ 7 ║ 3 │ 9 │ 2 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 2 │ 7 │ 3 ║ 1 │ 6 │ 9 ║ 4 │ 5 │ 8 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 4 │ 5 │ 9 ║ 8 │ 3 │ 2 ║ 1 │ 7 │ 6 ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║ 3 │ 1 │ 5 ║ 2 │ 7 │ 6 ║ 9 │ 8 │ 4 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 9 │ 2 │ 8 ║ 3 │ 4 │ 5 ║ 6 │ 1 │ 7 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 7 │ 4 │ 6 ║ 9 │ 1 │ 8 ║ 2 │ 3 │ 5 ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║ 1 │ 8 │ 2 ║ 5 │ 9 │ 4 ║ 7 │ 6 │ 3 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 6 │ 9 │ 4 ║ 7 │ 8 │ 3 ║ 5 │ 2 │ 1 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 5 │ 3 │ 7 ║ 6 │ 2 │ 1 ║ 8 │ 4 │ 9 ║
╚═══╧═══╧═══╩═══╧═══╧═══╩═══╧═══╧═══╝
╔═══╤═══╤═══╦═══╤═══╤═══╦═══╤═══╤═══╗
║   │   │ 1 ║   │   │ 7 ║ 3 │   │   ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 2 │   │ 3 ║   │   │   ║ 4 │ 5 │ 8 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 4 │   │ 9 ║   │   │ 2 ║   │   │   ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║ 3 │ 1 │   ║   │   │   ║ 9 │   │   ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║   │   │ 8 ║   │ 4 │   ║ 6 │   │   ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║   │   │ 6 ║   │   │   ║   │ 3 │ 5 ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║   │   │   ║ 5 │   │   ║ 7 │   │ 3 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 6 │ 9 │ 4 ║   │   │   ║ 5 │   │ 1 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║   │   │ 7 ║ 6 │   │   ║ 8 │   │   ║
╚═══╧═══╧═══╩═══╧═══╧═══╩═══╧═══╧═══╝
Execution time: 342ms
```

## Generating a printable board

You only need to supply the `-save-img` flag. The result looks like this:

![Sudoku puzzle](sample/sudoku.png "Sudoku puzzle")
