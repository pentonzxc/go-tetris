package tetris

type Tetris interface {
	Start()
	Points() int
	IsOver() bool
}
