package tetris

import "image/color"

const (
	tetrisWidth  = 10
	tetrisHeight = 20
)

var (
	rotate    command = plainCommand("rotate")
	moveLeft  command = plainCommand("moveLeft")
	moveRight command = plainCommand("moveRight")
	moveDown  command = plainCommand("moveDown")
	place     command = plainCommand("place")
	generate  command = plainCommand("generate")

	rowCollapse command = plainCommand("rowCollapse")
)

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
)
