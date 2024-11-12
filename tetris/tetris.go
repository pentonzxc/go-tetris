package tetris

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

type Command interface {
	name() string
}

type PlainCommand string

type BlockCommand struct {
	command string
	block   Block
}

func (bc *BlockCommand) name() string {
	return bc.command
}

func (s PlainCommand) name() string {
	return string(s)
}

var (
	Rotate    Command = PlainCommand("rotate")
	MoveLeft  Command = PlainCommand("moveLeft")
	MoveRight Command = PlainCommand("moveRight")
	MoveDown  Command = PlainCommand("moveDown")
	Place     Command = PlainCommand("place")
	Generate  Command = PlainCommand("generate")

	RowCollapse         Command = PlainCommand("rowCollapse")
	MoveDownLaggedBlock Command = PlainCommand("laggedBlock")
)

func Shape(blocks [][]bool, x int, y int, color color.Color) Block {
	var res Block

	res.blocks = make([][]Cell, len(blocks))
	for y := range len(blocks) {
		res.blocks[y] = make([]Cell, len(blocks[0]))

		for x := range len(blocks[0]) {
			res.blocks[y][x] = Cell{blocks[y][x], color}
		}
	}

	res.x = x
	res.y = y
	return res
}

const StartY = -1

var (
	OShape = Shape(
		[][]bool{{true, true}, {true, true}},
		TETRIS_WIDTH/2-1,
		StartY,
		color.RGBA{255, 255, 0, 255},
	)

	IShape = Shape(
		[][]bool{
			{false, true, false, false},
			{false, true, false, false},
			{false, true, false, false},
			{false, true, false, false},
		},
		TETRIS_WIDTH/2,
		StartY,
		color.RGBA{107, 202, 226, 255},
	)

	SShape = Shape(
		[][]bool{
			{false, false, false},
			{false, true, true},
			{true, true, false},
		},
		TETRIS_WIDTH/2-1,
		StartY,
		color.RGBA{255, 0, 0, 255},
	)

	ZShape = Shape(
		[][]bool{
			{false, false, false},
			{true, true, false},
			{false, true, true},
		},
		TETRIS_WIDTH/2-1,
		StartY,
		color.RGBA{0, 128, 0, 255},
	)

	LShape = Shape(
		[][]bool{
			{true, false, false},
			{true, false, false},
			{true, true, false},
		},
		TETRIS_WIDTH/2-1,
		StartY,
		color.RGBA{255, 165, 0, 255},
	)

	JShape = Shape(
		[][]bool{
			{false, false, true},
			{false, false, true},
			{false, true, true},
		},
		TETRIS_WIDTH/2-1,
		StartY,
		color.RGBA{255, 105, 180, 255},
	)

	TShape = Shape(
		[][]bool{
			{false, false, false},
			{true, true, true},
			{false, true, false},
		},
		TETRIS_WIDTH/2-1,
		StartY,
		color.RGBA{128, 0, 128, 255},
	)
)

type Cell struct {
	free  bool
	color color.Color
}

type Block struct {
	blocks [][]Cell
	x, y   int
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
			if block.blocks[y][x].free {
				s.blocks[block.y+y][block.x+x] = true
			}
		}
	}
}

// check borders and free of cells
func (s *TetrisState) isCellsValid(block Block) bool {
	for y := 0; y < len(block.blocks); y++ {
		for x := 0; x < len(block.blocks[0]); x++ {
			if block.blocks[y][x].free {
				cellY, cellX := block.y+y, block.x+x
				if cellY < 0 || cellY >= TETRIS_HEIGHT || cellX < 0 || cellX >= TETRIS_WIDTH {
					return false // Out of bounds
				}

				if s.blocks[cellY][cellX] {
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

			// previous element
			if last != nil {
				prev = *last
				prev.y -= 1
			}

			switch command {
			case Rotate:
				rotated := g.drawer.Rotate(*last)
				if g.state.isCellsValid(rotated) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(rotated)

					g.state.lastBlock = &rotated
				}
			case MoveLeft:
				moved := g.drawer.MoveLeft(*last)
				if g.state.isCellsValid(moved) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(moved)

					g.state.lastBlock = &moved
				}
			case MoveRight:
				moved := g.drawer.MoveRight(*last)
				if g.state.isCellsValid(moved) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(moved)

					g.state.lastBlock = &moved
				}
			case RowCollapse:
				for {
					removed := make([]int, 5)

					for y := range g.state.blocks {
						full := true

						for x := range g.state.blocks[y] {
							full = full && g.state.blocks[y][x]
						}

						if full {
							_ = append(removed, y)
						}
					}

					for {
						// maybe not init with size,
						visited := make(map[image.Point]bool, 10)
						lagged := make([]Block, 10)
						for y := removed[len(removed)-1] + 1; y >= 0; y++ {
							for x := 0; x < TETRIS_WIDTH; x++ {
								if g.state.blocks[y][x] {
									_ = append(lagged, findLaggedBlock(x, y, visited))
								}
							}
						}
					}

					if len(removed) == 0 {
						// return and update state
					}
				}
			case Place:
				fmt.Println("Place a block", *last)
				g.drawer.UndoBlock(prev)
				g.drawer.DrawBlock(*last)
			case Generate:
				new := generateBlock()
				g.state.lastBlock = &new

				g.commandQueue <- MoveDown
			case MoveDown:
				prev.y = g.state.lastBlock.y
				g.state.lastBlock.y += 1

				if g.state.isCellsValid(*last) {
					g.commandQueue <- Place
					go time.AfterFunc(250*time.Millisecond, func() {
						g.commandQueue <- MoveDown
					})
				} else {
					g.state.addBlock(prev)
					g.state.lastBlock = nil

					// g.commandQueue <- RowCollapse
					g.commandQueue <- Generate
				}
			}

			refresh()
		}
	}

}

func findLaggedBlock(x, y int, visited map[image.Point]bool) Block {

	return Block{}

}

func generateBlock() Block {
	var (
		block Block
		copy  Block
	)
	randomNum := rand.Int() % 7

	switch randomNum {
	case 0:
		block = TShape
	case 1:
		block = OShape
	case 2:
		block = IShape
	case 3:
		block = SShape
	case 4:
		block = ZShape
	case 5:
		block = LShape
	case 6:
		block = JShape
	}

	copy = block
	copy.blocks = make([][]Cell, len(block.blocks))

	for i := range block.blocks {
		copy.blocks[i] = make([]Cell, len(block.blocks[0]))
		for j := range copy.blocks[0] {
			copy.blocks[j] = block.blocks[j]
		}
	}

	return copy
}

func (g *TetrisFacade) Points() int {
	return g.state.score
}

func (g *TetrisFacade) IsOver() bool {
	return g.state.isGameOver
}
