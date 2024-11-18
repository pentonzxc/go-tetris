package tetris

import (
	"fmt"
	"image/color"
)

type cell struct {
	taken bool
	color color.Color
}

type block struct {
	cells [][]cell
	x, y  int
}

type tetrisState struct {
	cells      [][]cell
	score      int
	isGameOver bool

	fallingBlock *block
}

type command interface {
	name() string
}

type plainCommand string

func (s plainCommand) name() string {
	return string(s)
}

func newTetrisState() tetrisState {
	cells := make([][]cell, tetrisHeight)
	for i := 0; i < tetrisHeight; i++ {
		cells[i] = make([]cell, tetrisWidth)
	}

	return tetrisState{
		cells:      cells,
		score:      0,
		isGameOver: false,
	}
}

func newBlock(cells [][]cell, x, y int) block {
	return block{cells, x, y}
}

func (s *tetrisState) addBlock(block block) {
	for y := 0; y < len(block.cells); y++ {
		for x := 0; x < len(block.cells[0]); x++ {
			if block.cells[y][x].taken {
				s.cells[block.y+y][block.x+x] = block.cells[y][x]
			}
		}
	}
}

func (s *tetrisState) isCellsValid(block block) bool {
	for y := 0; y < len(block.cells); y++ {
		for x := 0; x < len(block.cells[0]); x++ {

			if block.cells[y][x].taken {
				cellY, cellX := block.y+y, block.x+x

				// border check
				if cellY < 0 || cellY >= tetrisHeight || cellX < 0 || cellX >= tetrisWidth {
					fmt.Printf("len - %v, y - %v\n", cellY, len(block.cells))
					return false
				}

				// freeness check
				if s.cells[cellY][cellX].taken {
					fmt.Println("here")
					return false
				}
			}

		}
	}

	return true
}
