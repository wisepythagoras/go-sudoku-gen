package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/wisepythagoras/go-sudoku-gen/sudoku"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Create the image for a Sudoku puzzle.
func CreateImage(puzzle *sudoku.Sudoku) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, 1029, 1029))

	drawGrid(img)

	// Loop through the rows.
	for i := 0; i < 3; i++ {
		box1 := puzzle.GetBox(i * 3)
		box2 := puzzle.GetBox(i*3 + 1)
		box3 := puzzle.GetBox(i*3 + 2)
		startVPos := 85 + (110 * i * 3)

		if i == 2 {
			startVPos += 9
		} else if i == 1 {
			startVPos += 3
		}

		for j := 0; j < 3; j++ {
			row := make([]uint8, 3)
			vPos := startVPos + (110 * j)

			copy(row, box1.GetRow(j))

			row = append(row, box2.GetRow(j)...)
			row = append(row, box3.GetRow(j)...)

			for k, n := range row {
				if n == 0 {
					continue
				}

				hPos := 55 + 110*k

				if k >= 6 {
					hPos += 9
				} else if k >= 3 {
					hPos += 3
				}

				label := fmt.Sprintf("%d", n)
				addLabel(img, hPos, vPos, label)
			}
		}
	}

	return img, nil
}

func drawGrid(img *image.RGBA) {
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	drawRectangle(img, color.Black, 9, 9, 1020, 1020, 6)
	drawRectangle(img, color.Black, 346, 9, 346, 1020, 3)
	drawRectangle(img, color.Black, 680, 9, 680, 1020, 3)
	drawRectangle(img, color.Black, 9, 346, 1020, 346, 3)
	drawRectangle(img, color.Black, 9, 680, 1020, 680, 3)

	for i := 0; i < 9; i++ {
		pos := 10 + 111*(i+1) + i
		drawRectangle(img, color.Black, pos, 10, pos, 1020, 1)
		drawRectangle(img, color.Black, 10, pos, 1020, pos, 1)
	}
}

func drawRectangle(img draw.Image, color color.Color, x1, y1, x2, y2, thickness int) {
	for i := x1; i < x2; i++ {
		for j := 0; j < thickness; j++ {
			img.Set(i, y1+j, color)
			img.Set(i, y2-j, color)
		}
	}

	for i := y1; i <= y2; i++ {
		for j := 0; j < thickness; j++ {
			img.Set(x1+j, i, color)
			img.Set(x2-j, i, color)
		}
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	myFont, _ := opentype.Parse(goregular.TTF)
	fontFace, _ := opentype.NewFace(myFont, &opentype.FaceOptions{
		Size: 12,
		DPI:  300.,
	})

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: fontFace,
		Dot:  point,
	}
	d.DrawString(label)
}
