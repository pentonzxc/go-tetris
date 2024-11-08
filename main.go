package main

import (
	"errors"
	"image"
	"image/color"
	"tetris/tetris"
	"time"

	"fyne.io/fyne/v2/canvas"
)

const TETRIS_WIDTH = 10
const TETRIS_HEIGHT = 20

func MakeTetrisGrid() *canvas.Raster {
	return canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		dx := w / TETRIS_WIDTH
		dy := h / TETRIS_HEIGHT

		if x%dx == 0 || y%dy == 0 {
			return color.Black
		} else {
			return color.White
		}
	})
}

type BlockPosition struct {
	x, y int
}

type Block struct {
	positions []BlockPosition
	color     color.Color
}

func MakeGridPosition(x, y int) (BlockPosition, error) {
	var pos BlockPosition

	if x < TETRIS_WIDTH && y < TETRIS_HEIGHT {
		return BlockPosition{
			x: x,
			y: y,
		}, nil
	} else {
		return pos, errors.New("invalid grid sizes")
	}
}

func DrawBlock(oldRaster *canvas.Raster, pos BlockPosition, color color.Color) {
	oldRasterGenerator := oldRaster.Generator

	newRaster := canvas.NewRaster(func(w, h int) image.Image {
		dx := w / TETRIS_WIDTH
		dy := h / TETRIS_HEIGHT

		sourceImg := oldRasterGenerator(w, h)
		newImg := image.NewRGBA(image.Rect(0, 0, w, h))

		// copy old image to new one
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				newImg.Set(x, y, sourceImg.At(x, y))
			}
		}

		x1 := dx * pos.x
		y1 := dy * pos.y
		// fmt.Println(x1)
		// fmt.Println(y1)

		// draw on new image a block
		for y := y1; y < y1+dy; y++ {
			for x := x1; x < x1+dx; x++ {
				// fmt.Printf("x - %d, y - %d", x, y)
				newImg.Set(x, y, color)
			}
		}

		return newImg
	})

	oldRaster.Generator = newRaster.Generator
	oldRaster.Refresh()
}

// func DrawBlockInternal(image *image.RGBA, pos BlockPosition, color color.Color, dx, dy int) {
// 	x1 := dx * pos.x
// 	y1 := dy * pos.y

// 	// draw on new image a block
// 	for y := y1; y < y1+dy; y++ {
// 		for x := x1; x < x1+dx; x++ {
// 			// fmt.Printf("x - %d, y - %d", x, y)
// 			image.Set(x, y, color)
// 		}
// 	}
// }

func DrawBlockInternal(image *image.RGBA, color color.Color, x1, x2, y1, y2 int) {
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			image.Set(x, y, color)
		}
	}
}

func FallingShape(raster *canvas.Raster, block Block) {
	// block needs to be sorted
	// block needs to be coupled

	generator := raster.Generator

	newRaster := canvas.NewRaster(func(w, h int) image.Image {
		dx := w / TETRIS_WIDTH
		dy := h / TETRIS_HEIGHT

		sourceImg := generator(w, h)
		newImg := image.NewRGBA(image.Rect(0, 0, w, h))

		// copy old image to new one
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				newImg.Set(x, y, sourceImg.At(x, y))
			}
		}

		for i := 0; i < len(block.positions); i++ {
			pos := block.positions[i]

			// WHY: maybe it's not necessary
			// if block.positions[i].y > 0 {
			// oldPos := BlockPosition{pos.x, pos.y - 1}
			// DrawBlockInternal2(newImg, color.White, oldPos.x*dx+1, oldPos.x*(dx+1)-1, oldPos.y*dy+1, oldPos.y*(dy+1)-1)
			// }

			DrawBlockInternal(newImg, block.color, pos.x*dx, (pos.x+1)*dx, pos.y*dy, (pos.y+1)*dy)
			block.positions[i].y += 1

		}

		return newImg
	})

	raster.Generator = newRaster.Generator

	for {
		if block.positions[0].y < TETRIS_HEIGHT {
			raster.Refresh()
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
}

func main() {
	// Pink := color.RGBA{245, 40, 145, 255}

	// a := app.New()
	// w := a.NewWindow("Hello")

	// w.Resize(fyne.NewSize(600, 800))
	// w.SetFixedSize(true)

	// grid := MakeTetrisGrid()
	// w.SetContent(grid)

	// block := Block{
	// 	positions: []BlockPosition{BlockPosition{1, 1}, BlockPosition{2, 1}, BlockPosition{1, 0}},
	// 	color:     Pink,
	// }

	// go FallingShape(grid, block)

	// // DrawBlock(grid, BlockPosition{1, 1}, Pink)
	// // DrawBlock(grid, BlockPosition{1, 2}, Pink)

	// w.ShowAndRun()

	game := tetris.NewGame(600, 800)

	game.Start()
}
