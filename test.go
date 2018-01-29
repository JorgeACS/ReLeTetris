package main

import (
	"image"
	"os"

	_ "image/png"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)
var block_size := 16.0
var IStates [4][16]int = [4][16]int{
	{
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0,
		0, 0, 1, 0,
	},
	{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 1, 1,
		0, 0, 0, 0,
	},
	{
		0, 1, 0, 0,
		0, 1, 0, 0,
		0, 1, 0, 0,
		0, 1, 0, 0,
	},
}

var JStates [4][16]int = [4][16]int{
	{
		1, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 1, 0,
		0, 1, 0, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 0, 0, 0,
		1, 1, 1, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 0, 0,
		0, 1, 0, 0,
		1, 1, 0, 0,
		0, 0, 0, 0,
	},
}

var LStates [4][16]int = [4][16]int{
	{
		0, 0, 1, 0,
		1, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 0, 0,
		0, 1, 0, 0,
		0, 1, 1, 0,
		0, 0, 0, 0,
	},
	{
		0, 0, 0, 0,
		1, 1, 1, 0,
		1, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		1, 1, 0, 0,
		0, 1, 0, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
	},
}

var OStates [4][16]int = [4][16]int{
	{
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 1, 0,
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
}

var SStates [4][16]int = [4][16]int{
	{
		0, 1, 1, 0,
		1, 1, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 0, 0,
		0, 1, 1, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
	},
	{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		0, 0, 0, 0,
	},
	{
		1, 0, 0, 0,
		1, 1, 0, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
	},
}

var TStates [4][16]int = [4][16]int{
	{
		0, 1, 0, 0,
		1, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 0, 0,
		0, 1, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 0, 0, 0,
		1, 1, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 0, 0,
		1, 1, 0, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
	},
}

var ZStates [4][16]int = [4][16]int{
	{
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 0, 1, 0,
		0, 1, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 0,
	},
	{
		0, 0, 0, 0,
		1, 1, 0, 0,
		0, 1, 1, 0,
		0, 0, 0, 0,
	},
	{
		0, 1, 0, 0,
		1, 1, 0, 0,
		1, 0, 0, 0,
		0, 0, 0, 0,
	},
}
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

func positionBottomLeft(mat){
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.Moved(pixel.V(-400,600))
	for i:=0;i < 22-4;i++{
		mat = mat.Moved(pixel.V(0,-block_size))
	}
	return mat
}

func moveTetrominoGrid(mat,xMagnitude=1,yMagnitude=1){
	vector = pixel.V(block_size*xMagnitude,block_size*yMagnitude))
	mat = mat.Moved(vector)
	return mat
}

func drawTetromino(tetromino,state,x,y){
	mat := pixel.IM 
	mat = positionBottomLeft(mat)
	mat = moveTetrominoGrid(mat,x,y)

}
func createTetromino(letter){
	switch letter{
		case 'L':
	        return Lstates[0]
	    case 'J':
	        return Jstates[0]
	    case 'S':
	        return Sstates[0]
	    case 'Z':
	        return Zstates[0]
	    case 'T':
	        return Tstates[0]
	    case 'O':
	        return Ostates[0]
	    case 'I':
	        return Istates[0]
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