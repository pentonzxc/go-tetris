package main

import (
	"tetris/tetris"
)

func main() {
	game := tetris.NewTetris(600, 800)

	game.Start()
}
