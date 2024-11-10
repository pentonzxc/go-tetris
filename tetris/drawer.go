package tetris

import (
	"image"
	"image/color"
)

const (
	TETRIS_WIDTH  = 10
	TETRIS_HEIGHT = 20
)

type BlockPosition struct {
	x, y int
}

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
	for y := 0; y < len(block.blocks); y++ {
		for x := 0; x < len(block.blocks[0]); x++ {
			if block.blocks[y][x] {
				drawer.drawCell(BlockPosition{block.x + x, block.y + y}, block.color)
			}
		}
	}
}

func (drawer *TetrisDrawer) UndoBlock(block Block) {
	for y := 0; y < len(block.blocks); y++ {
		for x := 0; x < len(block.blocks[0]); x++ {
			if block.blocks[y][x] {
				drawer.drawCell(BlockPosition{block.x + x, block.y + y}, color.White)
			}
		}
	}
}

func (drawer *TetrisDrawer) Rotate(block Block) Block {
	rotated := make([][]bool, len(block.blocks))
	for i := range rotated {
		rotated[i] = make([]bool, len(block.blocks[0]))
		copy(rotated[i], block.blocks[i])
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

	return Block{blocks: rotated, x: block.x, y: block.y, color: block.color}
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
