package main

import (
	"image"
	"os"

	_ "image/png"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	pic, err := loadPicture("block.png")
	if err != nil {
		panic(err)
	}
	block_size := 16.0
	sprite := pixel.NewSprite(pic, pic.Bounds())
	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.Moved(pixel.V(-400,600))
	mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(0.25, 0.25))
	for i := 0; i < 22; i++{
		for j := 0; j < 10; j++ {
			sprite.Draw(win,mat)
			mat = mat.Moved(pixel.V(block_size,0))
		}
		mat = mat.Moved(pixel.V(-block_size*10,-block_size))
	}
	for !win.Closed() {
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}