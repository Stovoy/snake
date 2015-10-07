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
	GameWon
)

type Board struct {
	Width  int
	Height int
	// Position of the single food pellet. There will always be a food pellet.
	Food Food
	// The head of the snake.
	SnakeHead Motive
	// The snake's body, in order from front to back.
	SnakeBody []Point

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
	board.placeSnake()
	board.placeFood()

	board.SnakeBody = make([]Point, 0)

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
	}

	oldHeadPosition := board.SnakeHead.Point
	board.SnakeHead.move()
	if !board.isValid() {
		board.restoreState()
		return
	}
	if board.SnakeHead.Point == board.Food.Point {
		// The snake body doesn't include the head.
		// The snake body will be one longer due to this food.
		if len(board.SnakeBody) == (board.Width*board.Height)-2 {
			// This must be the last position - victory!
			board.State = GameWon
		}
		board.AteFood = true
	}

	if !board.AteFood {
		if len(board.SnakeBody) > 0 {
			board.SnakeBody = board.SnakeBody[1:]
			board.SnakeBody = append(board.SnakeBody, oldHeadPosition)
		}
	} else {
		board.SnakeBody = append(board.SnakeBody, oldHeadPosition)
		board.placeFood()
	}

	board.AteFood = false
}

func (board *Board) Rewind() {
	board.restoreState()
}

func (board *Board) placeSnake() {
	board.SnakeHead = Motive{
		Point{rand.Intn(board.Width), rand.Intn(board.Height)},
		Direction(rand.Intn(4))}
}

// Make sure this is only called if board.Empty has at least one spot.
func (board *Board) placeFood() {
	empty := [][]bool{}
	emptyCount := 0
	for x := 0; x < board.Width; x++ {
		empty = append(empty, make([]bool, board.Height))
		for y := 0; y < board.Height; y++ {
			empty[x][y] = true
			emptyCount++
		}
	}
	for i := 0; i < len(board.SnakeBody); i++ {
		bodyPosition := board.SnakeBody[i]
		empty[bodyPosition.X][bodyPosition.Y] = false
		emptyCount--
	}
	empty[board.SnakeHead.X][board.SnakeHead.Y] = false
	emptyCount--

	emptyTarget := rand.Intn(emptyCount)
	i := -1

	var position Point
	for x := 0; x < board.Width; x++ {
		for y := 0; y < board.Height; y++ {
			if empty[x][y] {
				i++
				if i == emptyTarget {
					position = Point{x, y}
					break
				}
			}
		}
		if i == emptyTarget {
			break
		}
	}
	board.Food = Food{position}
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
	for i := 0; i < len(board.SnakeBody); i++ {
		if position == board.SnakeBody[i] {
			return false
		}
	}
	return true
}

func (board *Board) saveState() {
	history := board.History
	board.History = nil
	savedBoard := Board(*board)

	// Copy SnakeBody, so it is unique.
	savedBoard.SnakeBody = append([]Point{}, board.SnakeBody...)
	history = append(history, savedBoard)
	board.History = history
}

func (board *Board) restoreState() {
	history := board.History
	board.History = nil
	previousBoard, history := history[len(history)-1], history[:len(history)-1]
	*board = previousBoard
	board.History = history
}
