package tetris

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

type BlockPosition struct {
	x, y int
}

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
	blocks     [][]bool
	score      int
	isGameOver bool
	lastBlock  Block
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
			if block.blocks[x][y] {
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
			if block.blocks[x][y] {
				if !(block.y+y < TETRIS_HEIGHT && block.x+x < TETRIS_WIDTH) {
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
			if block.blocks[x][y] {
				if s.blocks[block.y+y][block.x+x] {
					return false
				}
			}
		}
	}

	return true
}

type TetrisFacade struct {
	width, height int
	state         TetrisState
	drawer        Drawer
	mutex         sync.Mutex
}

func (g *TetrisFacade) Start() {
	a := app.New()
	w := a.NewWindow("Tetris")

	w.Resize(fyne.NewSize(float32(g.width), float32(g.height)))
	w.SetFixedSize(true)

	w.SetContent(canvas.NewImageFromImage(g.drawer.Init()))

	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		g.mutex.Lock()

		switch event.Name {
		case "Up":
			fmt.Println(g.state.lastBlock)
			g.drawer.UndoBlock(g.state.lastBlock)
			fmt.Println("Do rotate")
			g.drawer.DrawBlock(g.drawer.Rotate(g.state.lastBlock))
			w.Canvas().Content().Refresh()
		case "Left":
			// g.drawer.MoveLeft(g.state.failingBlock)
		case "Right":
			// g.drawer.MoveRight(g.state.failingBlock)
		}
		g.mutex.Unlock()
	})

	go g.run(func() { w.Canvas().Content().Refresh() })

	w.ShowAndRun()
}

func (g *TetrisFacade) run(refresh func()) {
	Pink := color.RGBA{245, 40, 145, 255}

	for {
		// here we like generate new blocks
		prevBlock := Block{
			[][]bool{
				{true, false},
				{true, true},
			},
			TETRIS_WIDTH / 2,
			0,
			Pink,
		}

		// left := g.drawer.Rotate(prevBlock)
		for {
			curBlock := Block{
				blocks: prevBlock.blocks,
				x:      prevBlock.x,
				y:      prevBlock.y,
				color:  prevBlock.color,
			}

			// copy(curBlock.positions, prevBlock.positions)

			g.state.lastBlock = curBlock

			g.mutex.Lock()

			curBlock.y += 1

			if !g.state.isCellsValid(curBlock) || !g.state.isCellsFree(curBlock) {
				g.state.addBlock(prevBlock)

				g.mutex.Unlock()

				break
			} else {
				// g.drawer.UndoBlock(prevBlock)
				// g.drawer.DrawBlock(curBlock)
				refresh()

				// g.state.failingBlock = curBlock
				prevBlock = curBlock
			}

			g.mutex.Unlock()

			time.Sleep(200 * time.Millisecond)
		}

		// outer logic
	}
}

func (g *TetrisFacade) Points() int {
	return g.state.score
}

func (g *TetrisFacade) IsOver() bool {
	return g.state.isGameOver
}
