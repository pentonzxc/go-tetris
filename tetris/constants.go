package tetris

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
