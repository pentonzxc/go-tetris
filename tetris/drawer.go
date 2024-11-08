package tetris

import (
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
}

type TetrisDrawer struct {
	image   *image.RGBA
	dx, dy  int
	refresh func()
}

func (drawer *TetrisDrawer) DrawBlock(block Block) {
	for _, pos := range block.positions {
		drawer.drawCell(pos, block.color)
	}
}

func (drawer *TetrisDrawer) UndoBlock(block Block) {
	for _, pos := range block.positions {
		drawer.drawCell(pos, color.White)
	}
}

func (drawer *TetrisDrawer) Rotate(block Block) Block {
	panic("not implemented")
}

func (drawer *TetrisDrawer) MoveLeft(block Block) Block {
	panic("not implemented")
}

func (drawer *TetrisDrawer) MoveRight(block Block) Block {
	panic("not implemented")
}

func (drawer *TetrisDrawer) Refresh() {
	drawer.refresh()
}

func (drawer *TetrisDrawer) drawCell(pos BlockPosition, color color.Color) {
	x1 := pos.x*drawer.dx + 1
	x2 := (pos.x + 1) * drawer.dx
	y1 := pos.y*drawer.dy + 1
	y2 := (pos.y + 1) * drawer.dy

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
