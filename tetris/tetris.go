package tetris

import (
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

const (
	TETRIS_WIDTH  = 10
	TETRIS_HEIGHT = 20
)

type Tetris interface {
	Start()
	Points() int
	IsOver() bool
}

func NewGame(w, h int) Tetris {
	// why reference here
	return &GameFacade{
		width:  w,
		height: h,
		state:  TetrisState{},
	}
}

type TetrisState struct {
	blocks     []Block
	score      int
	isGameOver bool
}

type Block struct {
	positions []BlockPosition
	color     color.Color
}

type BlockPosition struct {
	x, y int
}

// it's like an impl of Tetris
type GameFacade struct {
	width, height int
	state         TetrisState
}

func (g *GameFacade) Start() {
	a := app.New()
	w := a.NewWindow("Tetris")

	dx := g.width / TETRIS_WIDTH
	dy := g.height / TETRIS_HEIGHT

	w.Resize(fyne.NewSize(float32(g.width), float32(g.height)))
	w.SetFixedSize(true)

	rawImage := image.NewRGBA(image.Rect(0, 0, g.width, g.height))

	// make a grid
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			var col color.Color

			if x%dx == 0 || y%dy == 0 {
				col = color.Black
			} else {
				col = color.White
			}

			rawImage.Set(x, y, col)
		}
	}

	canvasImage := canvas.NewImageFromImage(rawImage)
	w.SetContent(canvasImage)

	go g.run(rawImage, func() { canvasImage.Refresh() })

	w.ShowAndRun()
}

func (g *GameFacade) run(image *image.RGBA, refresh func()) {
	time.Sleep(3 * time.Second)
	Pink := color.RGBA{245, 40, 145, 255}

	dx := g.width / TETRIS_WIDTH
	dy := g.height / TETRIS_HEIGHT

	drawBlock(image, Pink, dx*1, dx*2, dy*1, dy*2)
	refresh()
}

func drawBlock(image *image.RGBA, color color.Color, x1, x2, y1, y2 int) {
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			image.Set(x, y, color)
		}
	}
}

func (g *GameFacade) Points() int {
	return g.state.score
}

func (g *GameFacade) IsOver() bool {
	return g.state.isGameOver
}

func (g *GameFacade) MoveLeft() {
	// logic to move the current block left
}

func (g *GameFacade) MoveRight() {
	// logic to move the current block right
}

func (g *GameFacade) Rotate() {
	// logic to rotate the current block
}
