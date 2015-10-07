package snake

import (
	"math/rand"
)

type Rotation byte

const (
	Left Rotation = iota
	Right
	Forward
)

type Direction byte

const (
	North Direction = iota
	West
	South
	East
)

func (heading *Direction) rotate(rotation Rotation) {
	switch rotation {
	case Left:
		switch *heading {
		case North:
			*heading = West
		case West:
			*heading = South
		case South:
			*heading = East
		case East:
			*heading = North
		}
	case Right:
		switch *heading {
		case North:
			*heading = East
		case West:
			*heading = North
		case South:
			*heading = West
		case East:
			*heading = South
		}
	}
}

type Point struct {
	X int
	Y int
}

type Food struct {
	Point
}

// An element that can move.
type Motive struct {
	Point
	Heading Direction
}

// Information on a specific rotation.
type RotationPoint struct {
	Rotation
	Point
}

func (motive *Motive) move() {
	switch motive.Heading {
	case North:
		motive.Y--
	case West:
		motive.X--
	case South:
		motive.Y++
	case East:
		motive.X++
	}
}

type GameState byte

const (
	Playing GameState = iota
	GameOver
	GameWon
)

type Board struct {
	Width  int
	Height int
	// Position of the single food pellet. There will always be a food pellet.
	Food Food
	// Stores all empty points.
	Empty []Point
	// The two most important parts of the snake.
	// The middle components can be inferred through Empty and Food.
	SnakeHead Motive
	SnakeTail Motive
	// Stores a list of rotations the snake made, in order.
	// This allows the tail to follow.
	Rotations []RotationPoint

	State GameState
	// Stores whether food was eaten on the last movement.
	AteFood bool

	History []Board
}

func NewBoard(width int, height int) *Board {
	board := Board{Width: width, Height: height}

	board.Initialize()

	return &board
}

func (board *Board) Initialize() {
	// Initialize the empty list.
	board.Empty = make([]Point, (board.Width * board.Height))
	i := 0
	for x := 0; x < board.Width; x++ {
		for y := 0; y < board.Height; y++ {
			board.Empty[i] = Point{x, y}
			i++
		}
	}

	board.placeSnake()
	board.placeFood()

	board.Rotations = make([]RotationPoint, 0)

	board.State = Playing

	board.History = make([]Board, 0)
}

func (board *Board) Move(rotation Rotation) {
	if board.State != Playing {
		return
	}

	// Store this move in the history.
	board.saveState()

	if rotation != Forward {
		board.SnakeHead.Heading.rotate(rotation)
		board.enqueueRotation(rotation)
	}

	board.SnakeHead.move()
	if !board.isValid() {
		board.State = GameOver
		board.restoreState()
		return
	}
	if board.SnakeHead.Point == board.Food.Point {
		if len(board.Empty) == 0 {
			// This must be the last position - victory!
			board.State = GameWon
			return
		}
		board.AteFood = true
	}

	// Where the head is is no longer empty.
	board.clearEmpty(board.SnakeHead.Point)

	if !board.AteFood {
		// Where the tail was is now empty.
		board.addEmpty(board.SnakeTail.Point)

		if len(board.Rotations) > 0 {
			// Check if the snake tail should rotate.
			nextRotation := board.Rotations[0]
			if board.SnakeTail.Point == nextRotation.Point {
				board.SnakeTail.Heading.rotate(nextRotation.Rotation)
				board.popRotation()
			}
		}

		board.SnakeTail.move()
	} else {
		board.placeFood()
	}

	board.AteFood = false
}

func (board *Board) placeSnake() {
	i := rand.Intn(len(board.Empty))
	position := board.Empty[i]
	board.SnakeHead = Motive{position, Direction(rand.Intn(4))}
	board.SnakeTail = board.SnakeHead

	// Remove the position from Empty.
	board.clearEmptyIndex(i)
}

// Make sure this is only called if board.Empty has at least one spot.
func (board *Board) placeFood() {
	i := rand.Intn(len(board.Empty))
	position := board.Empty[i]
	board.Food = Food{position}

	// Remove the position from Empty.
	board.clearEmptyIndex(i)
}

func (board *Board) enqueueRotation(rotation Rotation) {
	board.Rotations = append(board.Rotations, RotationPoint{rotation, board.SnakeHead.Point})
}

func (board *Board) popRotation() {
	board.Rotations = append(board.Rotations[:0], board.Rotations[1:]...)
}

func (board *Board) clearEmptyIndex(i int) {
	board.Empty = append(board.Empty[:i], board.Empty[i+1:]...)
}

func (board *Board) clearEmpty(point Point) {
	for i := 0; i < len(board.Empty); i++ {
		if point == board.Empty[i] {
			board.clearEmptyIndex(i)
			return
		}
	}
}

func (board *Board) addEmpty(point Point) {
	board.Empty = append(board.Empty, point)
}

func (board *Board) isValid() bool {
	position := board.SnakeHead.Point
	if position.X < 0 || position.X >= board.Width {
		return false
	}
	if position.Y < 0 || position.Y >= board.Height {
		return false
	}
	if position == board.Food.Point {
		return true
	}
	for i := 0; i < len(board.Empty); i++ {
		if position == board.Empty[i] {
			return true
		}
	}
	return false
}

func (board *Board) saveState() {
	history := board.History
	board.History = nil
	history = append(history, Board(*board))
	board.History = history
}

func (board *Board) restoreState() {
	history := board.History
	previousBoard, history := history[len(history)-1], history[:len(history)-1]
	*board = previousBoard
	board.History = history
}
