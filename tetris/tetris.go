package tetris

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

type Command string

const (
	Rotate    Command = "rotate"
	MoveLeft  Command = "moveLeft"
	MoveRight Command = "moveRight"
	MoveDown  Command = "moveDown"
	Place     Command = "place"
	Generate  Command = "generate"
)

type Block struct {
	blocks [][]bool
	x, y   int
	color  color.Color
}

type Tetris interface {
	Start()
	Points() int
	IsOver() bool
}

type TetrisState struct {
	blocks     [][]bool
	score      int
	isGameOver bool

	// can be nil
	lastBlock *Block
}

type TetrisFacade struct {
	width, height int
	state         TetrisState
	drawer        Drawer
	commandQueue  chan Command
}

func NewTetrisGame(w, h int) Tetris {
	// why reference here
	return &TetrisFacade{
		width:        w,
		height:       h,
		state:        newTetrisState(),
		drawer:       newTetrisDrawer(image.NewRGBA(image.Rect(0, 0, w, h)), func() {}),
		commandQueue: make(chan Command, 10),
	}
}

func newTetrisState() TetrisState {
	arr := make([][]bool, TETRIS_HEIGHT)
	for i := 0; i < TETRIS_HEIGHT; i++ {
		arr[i] = make([]bool, TETRIS_WIDTH)
	}

	return TetrisState{
		blocks:     arr,
		score:      0,
		isGameOver: false,
	}
}

func (s *TetrisState) addBlock(block Block) {
	for y := 0; y < len(block.blocks); y++ {
		for x := 0; x < len(block.blocks[0]); x++ {
			if block.blocks[y][x] {
				s.blocks[block.y+y][block.x+x] = true
			}
		}
	}
}

func (s *TetrisState) isCellsValid(block Block) bool {
	if !(block.y < TETRIS_HEIGHT && block.x < TETRIS_WIDTH) {
		return false
	}

	for y := 0; y < len(block.blocks); y++ {
		for x := 0; x < len(block.blocks[0]); x++ {
			if block.blocks[y][x] {
				cellY, cellX := block.y+y, block.x+x
				if cellY < 0 || cellY >= TETRIS_HEIGHT || cellX < 0 || cellX >= TETRIS_WIDTH {
					return false
				}
			}
		}
	}

	return true
}

func (s *TetrisState) isCellsFree(block Block) bool {
	for y := 0; y < len(block.blocks); y++ {
		for x := 0; x < len(block.blocks[0]); x++ {
			if block.blocks[y][x] {
				cellY, cellX := block.y+y, block.x+x
				if cellY < 0 || cellY >= TETRIS_HEIGHT || cellX < 0 || cellX >= TETRIS_WIDTH {
					return false // Out of bounds
				}
				if s.blocks[cellY][cellX] {
					fmt.Println(cellY, cellX)
					return false // Cell is occupied
				}
			}
		}
	}

	return true
}

func (g *TetrisFacade) Start() {
	a := app.New()
	w := a.NewWindow("Tetris")

	w.Resize(fyne.NewSize(float32(g.width), float32(g.height)))
	w.SetFixedSize(true)

	w.SetContent(canvas.NewImageFromImage(g.drawer.Init()))

	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		switch event.Name {
		case "Up":
			g.commandQueue <- Rotate
		case "Left":
			g.commandQueue <- MoveLeft
		case "Right":
			g.commandQueue <- MoveRight
		}

	})

	go g.processCommands(func() { w.Canvas().Content().Refresh() })
	g.commandQueue <- Generate

	w.ShowAndRun()
}

func (g *TetrisFacade) processCommands(refresh func()) {
	for command := range g.commandQueue {
		if g.state.lastBlock != nil || command == Generate {
			var prev Block
			last := g.state.lastBlock
			if last != nil {
				prev = *last
				prev.y -= 1
			}

			switch command {
			case Rotate:
				rotated := g.drawer.Rotate(*last)
				if g.state.isCellsValid(rotated) && g.state.isCellsFree(rotated) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(rotated)

					g.state.lastBlock = &rotated
				}
			case MoveLeft:
				moved := g.drawer.MoveLeft(*last)
				if g.state.isCellsValid(moved) && g.state.isCellsFree(moved) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(moved)

					g.state.lastBlock = &moved
				}
			case MoveRight:
				moved := g.drawer.MoveRight(*last)
				if g.state.isCellsValid(moved) && g.state.isCellsFree(moved) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(moved)

					g.state.lastBlock = &moved
				}
			case Place:
				g.drawer.UndoBlock(prev)
				g.drawer.DrawBlock(*last)
			case Generate:
				new := generateBlock()
				g.state.lastBlock = &new

				g.commandQueue <- MoveDown
			case MoveDown:
				g.state.lastBlock.y += 1

				if g.state.isCellsValid(*last) && g.state.isCellsFree(*last) {
					g.commandQueue <- Place
					g.commandQueue <- MoveDown
				} else {
					g.state.addBlock(prev)

					g.state.lastBlock = nil
					g.commandQueue <- Generate
				}

				time.Sleep(200 * time.Millisecond)
			}

			refresh()
		}
	}

}

func generateBlock() Block {
	Pink := color.RGBA{245, 40, 145, 255}
	return Block{
		[][]bool{
			{true, false},
			{true, true},
		},
		TETRIS_WIDTH / 2,
		-1,
		Pink,
	}
}

func (g *TetrisFacade) Points() int {
	return g.state.score
}

func (g *TetrisFacade) IsOver() bool {
	return g.state.isGameOver
}
