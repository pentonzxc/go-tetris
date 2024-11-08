package tetris

import (
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

type BlockPosition struct {
	x, y int
}

type Block struct {
	positions []BlockPosition
	color     color.Color
}

type Tetris interface {
	Start()
	Points() int
	IsOver() bool
}

func NewTetrisGame(w, h int) Tetris {
	// why reference here
	return &TetrisFacade{
		width:  w,
		height: h,
		state:  newTetrisState(),
		drawer: newTetrisDrawer(image.NewRGBA(image.Rect(0, 0, w, h)), func() {}),
	}
}

type TetrisState struct {
	blocks     []Block
	positions  map[BlockPosition]bool
	score      int
	isGameOver bool
}

func newTetrisState() TetrisState {
	return TetrisState{
		blocks:     make([]Block, 4),
		positions:  make(map[BlockPosition]bool),
		score:      0,
		isGameOver: false,
	}
}

func (s *TetrisState) addBlock(block Block) {
	s.blocks = append(s.blocks, block)

	for _, pos := range block.positions {
		s.positions[pos] = true
	}
}

func (s *TetrisState) isCellsValid(block Block) bool {
	for _, pos := range block.positions {
		if !(pos.y < TETRIS_HEIGHT && pos.x < TETRIS_WIDTH) {
			return false
		}
	}
	return true
}

func (s *TetrisState) isCellsFree(block Block) bool {
	for _, pos := range block.positions {
		if s.positions[pos] {
			return false
		}
	}
	return true
}

type TetrisFacade struct {
	width, height int
	state         TetrisState
	drawer        Drawer
}

func (g *TetrisFacade) Start() {
	a := app.New()
	w := a.NewWindow("Tetris")

	w.Resize(fyne.NewSize(float32(g.width), float32(g.height)))
	w.SetFixedSize(true)

	w.SetContent(canvas.NewImageFromImage(g.drawer.Init()))

	go g.run(func() { w.Canvas().Content().Refresh() })

	w.ShowAndRun()
}

func (g *TetrisFacade) run(refresh func()) {
	Pink := color.RGBA{245, 40, 145, 255}

	for {
		// here we like generate new blocks

		prevBlock := Block{
			positions: []BlockPosition{
				BlockPosition{1, 1},
				BlockPosition{2, 1},
				BlockPosition{1, 0},
			},
			color: Pink,
		}

		for {
			time.Sleep(200 * time.Millisecond)

			curBlock := Block{
				positions: make([]BlockPosition, len(prevBlock.positions)),
				color:     prevBlock.color,
			}

			copy(curBlock.positions, prevBlock.positions)
			for i := range curBlock.positions {
				curBlock.positions[i].y += 1
			}

			if !g.state.isCellsValid(curBlock) || !g.state.isCellsFree(curBlock) {
				g.state.addBlock(prevBlock)
				break
			} else {
				g.drawer.UndoBlock(prevBlock)

				g.drawer.DrawBlock(curBlock)
				refresh()

				prevBlock = curBlock
			}
		}

		// outer logic
	}
}

func whiteBlockFor(block Block) Block {
	var positions []BlockPosition = make([]BlockPosition, len(block.positions))
	copy(positions, block.positions)
	return Block{
		color:     color.White,
		positions: positions,
	}
}

func drawBlock(image *image.RGBA, block Block, w, h int) {
	dx := w / TETRIS_WIDTH
	dy := h / TETRIS_HEIGHT

	for _, pos := range block.positions {
		drawBlockInternal(image, block.color, pos.x*dx+1, (pos.x+1)*dx, pos.y*dy+1, (pos.y+1)*dy)
	}
}

func drawBlockInternal(image *image.RGBA, color color.Color, x1, x2, y1, y2 int) {
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			image.Set(x, y, color)
		}
	}
}

func (g *TetrisFacade) Points() int {
	return g.state.score
}

func (g *TetrisFacade) IsOver() bool {
	return g.state.isGameOver
}

func (g *TetrisFacade) MoveLeft() {
	// logic to move the current block left
}

func (g *TetrisFacade) MoveRight() {
	// logic to move the current block right
}

func (g *TetrisFacade) Rotate() {
	// logic to rotate the current block
}
