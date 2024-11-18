package tetris

import (
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

type tetrisFacade struct {
	width, height int
	state         tetrisState
	drawer        drawer
	commandQueue  chan command
}

func NewTetris(w, h int) Tetris {
	return &tetrisFacade{
		width:        w,
		height:       h,
		state:        newTetrisState(),
		drawer:       newTetrisDrawer(image.NewRGBA(image.Rect(0, 0, w, h))),
		commandQueue: make(chan command, 10),
	}
}

func (g *tetrisFacade) Start() {
	a := app.New()
	w := a.NewWindow("Tetris")

	w.Resize(fyne.NewSize(float32(g.width), float32(g.height)))
	w.SetFixedSize(true)

	w.SetContent(canvas.NewImageFromImage(g.drawer.Init()))

	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		switch event.Name {
		case "Up":
			g.commandQueue <- rotate
		case "Left":
			g.commandQueue <- moveLeft
		case "Right":
			g.commandQueue <- moveRight
		}

	})

	go g.processCommands(func() { w.Canvas().Content().Refresh() })
	g.commandQueue <- generate

	w.ShowAndRun()
}

func (g *tetrisFacade) processCommands(refresh func()) {
	for command := range g.commandQueue {
		if g.state.fallingBlock != nil || generate == command || rowCollapse == command {
			var prev block
			last := g.state.fallingBlock

			// previous element
			if last != nil {
				prev = *last
				prev.y -= 1
			}

			switch command {
			case rotate:
				rotated := g.drawer.Rotate(*last)
				if g.state.isCellsValid(rotated) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(rotated)

					g.state.fallingBlock = &rotated
				}
			case moveLeft:
				moved := g.drawer.MoveLeft(*last)
				if g.state.isCellsValid(moved) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(moved)

					g.state.fallingBlock = &moved
				}
			case moveRight:
				moved := g.drawer.MoveRight(*last)
				if g.state.isCellsValid(moved) {
					g.drawer.UndoBlock(*last)
					g.drawer.UndoBlock(prev)
					g.drawer.DrawBlock(moved)

					g.state.fallingBlock = &moved
				}
			case rowCollapse:
				for {
					var removed []int

					for y := range g.state.cells {
						full := true

						for x := range g.state.cells[y] {
							full = full && g.state.cells[y][x].taken
						}

						if full {
							removed = append(removed, y)

							for x := range g.state.cells[y] {
								g.state.cells[y][x] = cell{}
								g.drawer.DrawCell(image.Point{x, y}, color.White)
							}

						}
					}

					log.Printf("rowCollapse: removed arr - %v\n", removed)

					// exit case
					if len(removed) == 0 {
						break
					}

					visited := make(map[image.Point]bool, 0)
					lagged := make([]block, 0)

					for y := removed[0] - 1; y >= 0; y-- {
						for x := 0; x < tetrisWidth; x++ {
							if g.state.cells[y][x].taken {
								lagged = append(lagged, findLaggedBlock(x, y, visited, g.state))
							}
						}
					}

					var wg sync.WaitGroup

					wg.Add(len(lagged))

					for i := range lagged {
						log.Printf("rowCollapse: lagged - %v\n", lagged[i])

						go func() {
							block := lagged[i]
							for {
								var prev = block
								block.y += 1

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
			case place:
				g.drawer.UndoBlock(prev)
				g.drawer.DrawBlock(*last)
			case generate:
				new := generateBlock()
				g.state.fallingBlock = &new

				g.commandQueue <- moveDown
			case moveDown:
				prev.y = g.state.fallingBlock.y
				g.state.fallingBlock.y += 1

				if g.state.isCellsValid(*last) {
					g.commandQueue <- place
					go time.AfterFunc(250*time.Millisecond, func() {
						g.commandQueue <- moveDown
					})
				} else {
					g.state.addBlock(prev)
					g.state.fallingBlock = nil

					g.commandQueue <- rowCollapse
					g.commandQueue <- generate
				}
			}

			refresh()
		}
	}

}

func findLaggedBlock(x, y int, visited map[image.Point]bool, state tetrisState) block {
	var (
		res   block
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

	tmpCells := make([][]cell, tetrisHeight)

	for y := range tmpCells {
		tmpCells[y] = make([]cell, tetrisWidth)
	}

	tmpCells[y][x] = state.cells[y][x]

	queue = append(queue, image.Point{x, y})
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		if !visited[p] {
			for _, way := range ways {
				newP := p.Add(way)
				if newP.X >= 0 && newP.X < tetrisWidth &&
					newP.Y >= 0 && newP.Y < tetrisHeight &&
					state.cells[newP.Y][newP.X].taken {

					minX = int(math.Min(float64(minX), float64(newP.X)))
					minY = int(math.Min(float64(minY), float64(newP.Y)))
					maxX = int(math.Max(float64(maxX), float64(newP.X)))
					maxY = int(math.Max(float64(maxY), float64(newP.Y)))

					tmpCells[newP.Y][newP.X] = state.cells[newP.Y][newP.X]
					state.cells[newP.Y][newP.X] = cell{}

					queue = append(queue, newP)
				}
			}

			visited[p] = true
		}
	}

	// log.Printf("max x - %v\n", maxX)
	// log.Printf("min x - %v\n", minX)
	// log.Printf("max y - %v\n", maxY)
	// log.Printf("min y - %v\n", minY)

	size := int(math.Max(float64(maxX-minX+1), float64(maxY-minY+1)))

	res.cells = make([][]cell, size)
	for dy := 0; dy < size; dy++ {
		res.cells[dy] = make([]cell, size)
		for dx := 0; dx < size; dx++ {
			cellX := minX + dx
			cellY := minY + dy
			if cellX >= minX && cellX <= maxX && cellY >= minY && cellY <= maxY {
				res.cells[dy][dx] = tmpCells[cellY][cellX]
			} else {
				// Fill empty cells with default (non-empty = false)
				res.cells[dy][dx] = cell{}
			}
		}
	}

	// Set the block's position
	res.x = minX
	res.y = minY

	return res
}

func generateBlock() block {
	var block, copy block
	randomNum := rand.Int() % 7

	switch randomNum {
	case 0:
		block = tShape
	case 1:
		block = oShape
	case 2:
		block = iShape
	case 3:
		block = sShape
	case 4:
		block = zShape
	case 5:
		block = lShape
	case 6:
		block = jShape
	}

	copy = block
	copy.cells = make([][]cell, len(block.cells))

	for i := range block.cells {
		copy.cells[i] = make([]cell, len(block.cells[0]))

		for j := range copy.cells[0] {
			copy.cells[j] = block.cells[j]
		}

	}

	return copy
}

func (g *tetrisFacade) Points() int {
	return g.state.score
}

func (g *tetrisFacade) IsOver() bool {
	return g.state.isGameOver
}

// func (g *tetrisFacade) testRowCollapse() {
// 	for y := len(g.state.cells) - 1; y >= len(g.state.cells)-2; y-- {
// 		for x := 0; x < len(g.state.cells[0]); x++ {
// 			g.state.cells[y][x] = cell{true, color.RGBA{0, 128, 0, 255}}
// 			g.drawer.DrawCell(image.Point{x, y}, g.state.cells[y][x].color)
// 		}
// 	}
// }
