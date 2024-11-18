package tetris

import (
	"image/color"
	"math/rand"
)

type shape block

func newShape(blocks [][]bool, color color.Color, x int) shape {
	var res block

	res.cells = make([][]cell, len(blocks))
	for y := range blocks {
		res.cells[y] = make([]cell, len(blocks[y]))

		for x := range blocks[y] {
			res.cells[y][x] = cell{blocks[y][x], color}
		}
	}

	res.x = x
	// start position
	res.y = -1
	return shape(res)
}

var (
	oShape = newShape(
		[][]bool{{true, true}, {true, true}},
		color.RGBA{255, 255, 0, 255},
		tetrisWidth/2-1,
	)

	iShape = newShape(
		[][]bool{
			{false, true, false, false},
			{false, true, false, false},
			{false, true, false, false},
			{false, true, false, false},
		},
		color.RGBA{107, 202, 226, 255},
		tetrisWidth/2,
	)

	sShape = newShape(
		[][]bool{
			{false, false, false},
			{false, true, true},
			{true, true, false},
		},
		color.RGBA{255, 0, 0, 255},
		tetrisWidth/2-1,
	)

	zShape = newShape(
		[][]bool{
			{false, false, false},
			{true, true, false},
			{false, true, true},
		},
		color.RGBA{0, 128, 0, 255},
		tetrisWidth/2-1,
	)

	lShape = newShape(
		[][]bool{
			{true, false, false},
			{true, false, false},
			{true, true, false},
		},
		color.RGBA{255, 165, 0, 255},
		tetrisWidth/2-1,
	)

	jShape = newShape(
		[][]bool{
			{false, false, true},
			{false, false, true},
			{false, true, true},
		},
		color.RGBA{255, 105, 180, 255},
		tetrisWidth/2-1,
	)

	tShape = newShape(
		[][]bool{
			{false, false, false},
			{true, true, true},
			{false, true, false},
		},
		color.RGBA{128, 0, 128, 255},
		tetrisWidth/2-1,
	)

	shapes = map[int]shape{
		0: tShape,
		1: oShape,
		2: iShape,
		3: sShape,
		4: zShape,
		5: lShape,
		6: jShape,
	}
)

func generateBlock() block {
	var res block
	shape := shapes[rand.Int()%7]

	res = block(shape)
	res.cells = make([][]cell, len(shape.cells))

	for i := range shape.cells {
		res.cells[i] = make([]cell, len(shape.cells[0]))

		for j := range res.cells[0] {
			res.cells[j] = shape.cells[j]
		}

	}

	return res
}
