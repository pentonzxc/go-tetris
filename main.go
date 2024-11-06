package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

const TETRIS_WIDTH = 10
const TETRIS_HEIGHT = 20

func MakeTetrisGrid() *canvas.Raster {
	raster := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		dx := w / TETRIS_WIDTH
		dy := h / TETRIS_HEIGHT

		if x%dx == 0 || y%dy == 0 {
			return color.Black
		} else {
			return color.White
		}
	})

	return raster
}

type BlockPosition struct {
	x, y int
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

func DrawBlock(grid *canvas.Raster, pos BlockPosition, color color.Color) *canvas.Raster {
	return canvas.NewRaster(func(w, h int) image.Image {
		dx := w / TETRIS_WIDTH
		dy := h / TETRIS_HEIGHT

		sourceImg := grid.Generator(w, h)
		newImg := image.NewRGBA(image.Rect(0, 0, w, h))

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				newImg.Set(x, y, sourceImg.At(x, y))
			}
		}

		x1 := dx * pos.x
		y1 := dy * pos.y
		fmt.Println(x1)
		fmt.Println(y1)

		for y := y1; y < y1+dy; y++ {
			for x := x1; x < x1+dx; x++ {
				fmt.Printf("x - %d, y - %d", x, y)
				newImg.Set(x, y, color)
			}
		}

		return newImg
	})
}

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	w.Resize(fyne.NewSize(600, 800))
	w.SetFixedSize(true)

	grid := MakeTetrisGrid()

	grid2 := DrawBlock(grid, BlockPosition{1, 1}, color.RGBA{245, 40, 145, 1})

	w.SetContent(grid2)

	// hello := widget.NewLabel("Hello Fyne!")

	// w.SetContent(container.NewVBox(
	// 	hello,
	// 	widget.NewButton("Hi!", func() {
	// 		hello.SetText("Welcome :)")
	// 	}),
	// ))

	w.ShowAndRun()
}
