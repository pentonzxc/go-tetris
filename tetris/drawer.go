package tetris

import (
	"fmt"
	"image"
	"image/color"
)

type drawer interface {
	UndoBlock(block)
	DrawBlock(block)
	Rotate(block) block
	MoveLeft(block) block
	MoveRight(block) block
	Init() *image.RGBA
	DrawCell(pos image.Point, color color.Color)
}

type tetrisDrawer struct {
	image  *image.RGBA
	dx, dy int
}

func (drawer *tetrisDrawer) DrawBlock(block block) {
	for y := 0; y < len(block.cells); y++ {
		for x := 0; x < len(block.cells[0]); x++ {
			if block.cells[y][x].taken {
				fmt.Printf("draw point - %v\n", image.Point{block.x + x, block.y + y})
				drawer.DrawCell(image.Point{block.x + x, block.y + y}, block.cells[y][x].color)
			}
		}
	}
}

func (drawer *tetrisDrawer) UndoBlock(block block) {
	for y := 0; y < len(block.cells); y++ {
		for x := 0; x < len(block.cells[0]); x++ {
			if block.cells[y][x].taken {
				drawer.DrawCell(image.Point{block.x + x, block.y + y}, color.White)
			}
		}
	}
}

func (drawer *tetrisDrawer) Rotate(block block) block {
	rotated := make([][]cell, len(block.cells))
	for i := range rotated {
		rotated[i] = make([]cell, len(block.cells[0]))
		copy(rotated[i], block.cells[i])
	}

	matrix := rotated

	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}

	for i := range matrix {
		for j := 0; j < i; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	return newBlock(rotated, block.x, block.y)
}

func (drawer *tetrisDrawer) MoveLeft(block block) block {
	block.x -= 1
	return block
}

func (drawer *tetrisDrawer) MoveRight(block block) block {
	block.x += 1
	return block
}

func (drawer *tetrisDrawer) DrawCell(pos image.Point, color color.Color) {
	x1 := pos.X*drawer.dx + 1
	x2 := (pos.X + 1) * drawer.dx
	y1 := pos.Y*drawer.dy + 1
	y2 := (pos.Y + 1) * drawer.dy

	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			drawer.image.Set(x, y, color)
		}
	}
}

func (drawer *tetrisDrawer) Init() *image.RGBA {
	x1 := drawer.image.Rect.Max.X
	y1 := drawer.image.Rect.Max.Y

	for y := 0; y < y1; y++ {
		for x := 0; x < x1; x++ {
			var col color.Color

			if x%drawer.dx == 0 || y%drawer.dy == 0 {
				col = color.Black
			} else {
				col = color.White
			}

			drawer.image.Set(x, y, col)
		}
	}

	return drawer.image
}

func newTetrisDrawer(image *image.RGBA) drawer {
	point := image.Rect.Max
	return &tetrisDrawer{
		image: image,
		dx:    point.X / tetrisWidth,
		dy:    point.Y / tetrisHeight,
	}
}
