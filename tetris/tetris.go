package tetris

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sync"
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

	res.cells = make([][]Cell, len(blocks))
	for y := range len(blocks) {
		res.cells[y] = make([]Cell, len(blocks[0]))

		for x := range len(blocks[0]) {
			res.cells[y][x] = Cell{blocks[y][x], color}
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
	NonEmpty bool
	color    color.Color
}

type Block struct {
	cells [][]Cell
	x, y  int
}

type Tetris interface {
	Start()
	Points() int
	IsOver() bool
}

type TetrisState struct {
	cells      [][]Cell
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
	cells := make([][]Cell, TETRIS_HEIGHT)
	for i := 0; i < TETRIS_HEIGHT; i++ {
		cells[i] = make([]Cell, TETRIS_WIDTH)
	}

	return TetrisState{
		cells:      cells,
		score:      0,
		isGameOver: false,
	}
}

func (s *TetrisState) addBlock(block Block) {
	for y := 0; y < len(block.cells); y++ {
		for x := 0; x < len(block.cells[0]); x++ {

			if block.cells[y][x].NonEmpty {
				s.cells[block.y+y][block.x+x] = block.cells[y][x]
			}

		}
	}
}

// check borders and free of cells
func (s *TetrisState) isCellsValid(block Block) bool {
	for y := 0; y < len(block.cells); y++ {
		for x := 0; x < len(block.cells[0]); x++ {

			if block.cells[y][x].NonEmpty {
				cellY, cellX := block.y+y, block.x+x

				// border check
				if cellY < 0 || cellY >= TETRIS_HEIGHT || cellX < 0 || cellX >= TETRIS_WIDTH {
					fmt.Printf("len - %v, y - %v\n", cellY, len(block.cells))
					return false
				}

				// freeness check
				if s.cells[cellY][cellX].NonEmpty {
					fmt.Println("here")
					return false
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

	g.testRowCollapse()

	go g.processCommands(func() { w.Canvas().Content().Refresh() })
	g.commandQueue <- Generate

	w.ShowAndRun()
}

func (g *TetrisFacade) processCommands(refresh func()) {
	for command := range g.commandQueue {
		if g.state.lastBlock != nil || command == Generate || command == RowCollapse {
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
					log.Println("called RowCollapse")
					removed := make([]int, 0)

					for y := range g.state.cells {
						full := true

						for x := range g.state.cells[y] {
							full = full && g.state.cells[y][x].NonEmpty
						}

						if full {
							removed = append(removed, y)

							for x := range g.state.cells[y] {
								g.state.cells[y][x] = Cell{}
								g.drawer.DrawCell(image.Point{x, y}, color.White)
							}

						}
					}

					log.Printf("removed arr - %v\n", removed)

					// exit case
					if len(removed) == 0 {
						break
					}

					visited := make(map[image.Point]bool, 0)
					lagged := make([]Block, 0)

					for y := removed[0] - 1; y >= 0; y-- {
						for x := 0; x < TETRIS_WIDTH; x++ {
							if g.state.cells[y][x].NonEmpty {
								lagged = append(lagged, findLaggedBlock(x, y, visited, g.state))
							}
						}
					}

					var wg sync.WaitGroup

					wg.Add(len(lagged))

					for i := range lagged {
						fmt.Printf("lagged - %v\n", lagged[i])

						go func() {
							block := lagged[i]
							for {
								var prev = block
								block.y += 1
								fmt.Println(block.y)

								if g.state.isCellsValid(block) {
									time.Sleep(10 * time.Millisecond)
									g.drawer.UndoBlock(prev)
									g.drawer.DrawBlock(block)
									refresh()
								} else {
									g.state.addBlock(prev)
									wg.Done()
									break
								}
							}
						}()
					}

					wg.Wait()
				}
			case Place:
				// fmt.Println("Place a block", *last)
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

					fmt.Println("go to next block")
					g.commandQueue <- RowCollapse
					g.commandQueue <- Generate
				}
			}

			refresh()
		}
	}

}

func findLaggedBlock(x, y int, visited map[image.Point]bool, state TetrisState) Block {
	var (
		res   Block
		queue []image.Point  = make([]image.Point, 0)
		ways  [4]image.Point = [4]image.Point{
			image.Point{1, 0},
			image.Point{0, 1},
			image.Point{-1, 0},
			image.Point{0, -1},
		}
		minX, minY int = x, y
		maxX, maxY int = x, y
	)

	tmpCells := make([][]Cell, TETRIS_HEIGHT)

	for y := range tmpCells {
		tmpCells[y] = make([]Cell, TETRIS_WIDTH)
	}

	tmpCells[y][x] = state.cells[y][x]

	queue = append(queue, image.Point{x, y})
	for len(queue) > 0 {
		p := queue[0]
		fmt.Println(p)
		queue = queue[1:]

		if !visited[p] {
			for _, way := range ways {
				newP := p.Add(way)
				if newP.X >= 0 && newP.X < TETRIS_WIDTH &&
					newP.Y >= 0 && newP.Y < TETRIS_HEIGHT &&
					state.cells[newP.Y][newP.X].NonEmpty {

					minX = int(math.Min(float64(minX), float64(newP.X)))
					minY = int(math.Min(float64(minY), float64(newP.Y)))
					maxX = int(math.Max(float64(maxX), float64(newP.X)))
					maxY = int(math.Max(float64(maxY), float64(newP.Y)))

					tmpCells[newP.Y][newP.X] = state.cells[newP.Y][newP.X]
					state.cells[newP.Y][newP.X] = Cell{}

					queue = append(queue, newP)
				}
			}

			visited[p] = true
		}
	}

	// fmt.Printf("max x - %v\n", maxX)
	// fmt.Printf("min x - %v\n", minX)
	// fmt.Printf("max y - %v\n", maxY)
	// fmt.Printf("min y - %v\n", minY)

	size := int(math.Max(float64(maxX-minX+1), float64(maxY-minY+1)))

	res.cells = make([][]Cell, size)
	for dy := 0; dy < size; dy++ {
		res.cells[dy] = make([]Cell, size)
		for dx := 0; dx < size; dx++ {
			cellX := minX + dx
			cellY := minY + dy
			if cellX >= minX && cellX <= maxX && cellY >= minY && cellY <= maxY {
				res.cells[dy][dx] = tmpCells[cellY][cellX]
			} else {
				// Fill empty cells with default (non-empty = false)
				res.cells[dy][dx] = Cell{}
			}
		}
	}

	// Set the block's position
	res.x = minX
	res.y = minY

	return res
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
	copy.cells = make([][]Cell, len(block.cells))

	for i := range block.cells {
		copy.cells[i] = make([]Cell, len(block.cells[0]))

		for j := range copy.cells[0] {
			copy.cells[j] = block.cells[j]
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

func (g *TetrisFacade) testRowCollapse() {
	for y := len(g.state.cells) - 1; y >= len(g.state.cells)-2; y-- {
		for x := 0; x < len(g.state.cells[0]); x++ {
			g.state.cells[y][x] = Cell{true, color.RGBA{0, 128, 0, 255}}
			g.drawer.DrawCell(image.Point{x, y}, g.state.cells[y][x].color)
		}
	}
}
