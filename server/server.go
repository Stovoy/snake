package server

import (
	"fmt"
	"net/http"

	"encoding/json"
	"snake/snake"
)

type Handler func(w http.ResponseWriter, r *http.Request)
type Modifier func()

func Game(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/game.html")
}

func ModifyAndWrite(board *snake.Board, modifier Modifier) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		modifier()
		WriteBoard(w, board)
	}
}

func MoveLeft(board *snake.Board) Handler {
	return ModifyAndWrite(board, func() {
		board.Move(snake.Left)
	})
}

func MoveRight(board *snake.Board) Handler {
	return ModifyAndWrite(board, func() {
		board.Move(snake.Right)
	})
}

func MoveForward(board *snake.Board) Handler {
	return ModifyAndWrite(board, func() {
		board.Move(snake.Forward)
	})
}

func Reset(board *snake.Board) Handler {
	return ModifyAndWrite(board, func() {
		board.Initialize()
	})
}

func Rewind(board *snake.Board) Handler {
	return ModifyAndWrite(board, func() {
		board.Rewind()
	})
}

func WriteBoard(w http.ResponseWriter, board *snake.Board) {
	data, err := json.Marshal(board)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, string(data))
	}
}

func Start(board *snake.Board) {
	// Serve static files.
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", Game)
	http.HandleFunc("/snake/move/left", MoveLeft(board))
	http.HandleFunc("/snake/move/right", MoveRight(board))
	http.HandleFunc("/snake/move/forward", MoveForward(board))
	http.HandleFunc("/snake/reset", Reset(board))
	http.HandleFunc("/snake/rewind", Rewind(board))

	err := http.ListenAndServe(":8000", nil)
	fmt.Println(err)
}
