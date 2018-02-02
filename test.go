package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"time"
)

const BLOCK_SIZE = 16.0
const GRID_SIZE = 4
const PLAYFIELD_WIDTH = 10
const PLAYFIELD_HEIGHT = 22
const PLAYFIELD_MAX_HEIGHT = 40
const INITIAL_STATE = 0
const MIN_GRAVITY = 0.125
const MAX_GRAVITY = 20
const NUM_PIECES = 7
const NEXT_LENGTH = 5

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

//----------------Game------------------
//Structure for the game
//contains the playfield plus other neccesary values
type Game struct {
	playfield    Playfield
	gravity      float32  //Rate at which pieces fall (X cells per frame)
	score        int      //Player score
	currentPiece string   // Piece on the playfield at the moment
	hold         string   // Piece currently held
	currentBag   []string // Currrent piece bag
	nextBag      []string // Next piece bag
	bagIndex     int
	next         []string // Piece preview
	canSwap      bool     // Determines if the current piece can be swapped
	pieceHeld    bool     // Determines if there is a held piece
	timer        int64    // Time since game started in milliseconds
	timePast     int64    //Time past in milliseconds
}

func createGame() Game {
	bag1 := generate7PieceBag()
	bag2 := generate7PieceBag()
	return Game{createPlayfield(), MIN_GRAVITY, 0, bag1[0], "", bag1, bag2, 0, bag1[1:6], true, false, 0, 0}
}

func updateGame(game Game, milliseconds int64) {
	game.timePast += milliseconds
	if float32(game.timePast/1000) > game.gravity {
		//pushPieceDown(Tetromino)
	}

}

//--------------Playfield---------------
//Structure for the playfield
//Contains the stack, which is composed of one or more Tetrominos
type Playfield struct {
	stack []Tetromino
}

//--------------Tetromino---------------
//Structure for a tetromino piece
//Contains its value, state, and it's position using (x,y) coordinates
type Tetromino struct {
	letter string
	state  int
	x      int
	y      int
}

func pushPieceDown(piece Tetromino) {
	piece.y -= 1
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

	MINO_PICTURE, err := loadPicture("block.png")
	MINO_SPRITE := pixel.NewSprite(MINO_PICTURE, MINO_PICTURE.Bounds())
	if err != nil {
		panic(err)
	}
	game := createGame()
	last := time.Now()
	for !win.Closed() {
		//Draw
		mat := pixel.IM
		mat = mat.Moved(win.Bounds().Center())
		mat = mat.Moved(pixel.V(-400, 600))
		mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(0.25, 0.25))
		for i := 0; i < PLAYFIELD_HEIGHT; i++ {
			for j := 0; j < PLAYFIELD_WIDTH; j++ {
				MINO_SPRITE.Draw(win, mat)
				mat = mat.Moved(pixel.V(BLOCK_SIZE, 0))
			}
			mat = mat.Moved(pixel.V(-BLOCK_SIZE*PLAYFIELD_WIDTH, -BLOCK_SIZE))
		}
		//Update
		updateGame(game, time.Since(last).Nanoseconds()*1000000)
		win.Update()
	}

}

func generate7PieceBag() []string {
	bag := []string{"I", "O", "T", "L", "J", "S", "Z"}
	shuffleBag(bag, 100)
	return bag

}
func shuffleBag(bag []string, shuffles int) {
	for aux := 0; aux < shuffles; aux++ {
		i := rand.Intn(NUM_PIECES)
		j := rand.Intn(NUM_PIECES)
		bag[i], bag[j] = bag[j], bag[i]
	}

}
func createPlayfield() Playfield {
	return Playfield{stack: []Tetromino{}}
}
func pushTetromino(playfield Playfield, piece Tetromino) {
	playfield.stack = append(playfield.stack, piece)
}

//positions the tettromino grid on the bottom left corner
func positionBottomLeft(win pixelgl.Window, mat pixel.Matrix) pixel.Matrix {
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.Moved(pixel.V(-400, 600))
	for i := 0; i < PLAYFIELD_HEIGHT-GRID_SIZE; i++ {
		mat = mat.Moved(pixel.V(0, -BLOCK_SIZE))
	}
	return mat
}

func moveTetrominoGrid(mat pixel.Matrix, xMagnitude int, yMagnitude int) pixel.Matrix {
	vector := pixel.V(float64(BLOCK_SIZE*xMagnitude), float64(BLOCK_SIZE*yMagnitude*1.0))
	mat = mat.Moved(vector)
	return mat
}

/*
func drawTetromino(win, piece Tetromino, state int, x int, y int) {
	mat := pixel.IM
	mat = positionBottomLeft(win, mat)
	mat = moveTetrominoGrid(mat, piece.x, piece.y)

	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			if Tetromino.piece == 'S' {

			}
		}
	}

}*/
func createTetromino(letter string) Tetromino {
	return Tetromino{letter, 0, PLAYFIELD_WIDTH / 2, PLAYFIELD_HEIGHT}
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

/*
	Ejemplo 4.2
	Ejercicio 4.5 (paginas 65 y 66)
	Ejercicio 4.9 (pagina 69)
*/
