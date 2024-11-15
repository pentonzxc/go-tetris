package tetris

import (
	"fmt"
	"image"
	"image/color"
)

const (
	TETRIS_WIDTH  = 10
	TETRIS_HEIGHT = 20
)

type Drawer interface {
	UndoBlock(Block)
	DrawBlock(Block)
	Rotate(Block) Block
	MoveLeft(Block) Block
	MoveRight(Block) Block
	Refresh()
	Init() *image.RGBA
	DrawCell(pos image.Point, color color.Color)
}

type TetrisDrawer struct {
	image   *image.RGBA
	dx, dy  int
	refresh func()
}

func (drawer *TetrisDrawer) DrawBlock(block Block) {
	for y := 0; y < len(block.cells); y++ {
		for x := 0; x < len(block.cells[0]); x++ {
			if block.cells[y][x].NonEmpty {
				fmt.Printf("draw point - %v\n", image.Point{block.x + x, block.y + y})
				drawer.DrawCell(image.Point{block.x + x, block.y + y}, block.cells[y][x].color)
			}
		}
	}
}

func (drawer *TetrisDrawer) UndoBlock(block Block) {
	for y := 0; y < len(block.cells); y++ {
		for x := 0; x < len(block.cells[0]); x++ {
			if block.cells[y][x].NonEmpty {
				drawer.DrawCell(image.Point{block.x + x, block.y + y}, color.White)
			}
		}
	}
}

func (drawer *TetrisDrawer) Rotate(block Block) Block {
	rotated := make([][]Cell, len(block.cells))
	for i := range rotated {
		rotated[i] = make([]Cell, len(block.cells[0]))
		copy(rotated[i], block.cells[i])
	}

	matrix := rotated

	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}

	// transpose it
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < i; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	return Block{cells: rotated, x: block.x, y: block.y}
}

func (drawer *TetrisDrawer) MoveLeft(block Block) Block {
	block.x -= 1
	return block
}

func (drawer *TetrisDrawer) MoveRight(block Block) Block {
	block.x += 1
	return block
}

func (drawer *TetrisDrawer) Refresh() {
	drawer.refresh()
}

func (drawer *TetrisDrawer) DrawCell(pos image.Point, color color.Color) {
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

func (drawer *TetrisDrawer) Init() *image.RGBA {
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

func newTetrisDrawer(image *image.RGBA, refresh func()) Drawer {
	point := image.Rect.Max
	return &TetrisDrawer{
		image:   image,
		dx:      point.X / TETRIS_WIDTH,
		dy:      point.Y / TETRIS_HEIGHT,
		refresh: refresh,
	}
}
