package main

import (
	"tetris/tetris"
)

func main() {
	game := tetris.NewTetrisGame(600, 800)

	game.Start()
}
