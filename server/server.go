package server

import (
	"fmt"
	"net/http"

	"snake/snake"
	"encoding/json"
)

type Handler func(w http.ResponseWriter, r *http.Request)

func Game(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/game.html")
}

func MoveLeft(board *snake.Board) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		board.Move(snake.Left)
		WriteBoard(w, board)
	}
}

func MoveRight(board *snake.Board) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		board.Move(snake.Right)
		WriteBoard(w, board)
	}
}

func MoveForward(board *snake.Board) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		board.Move(snake.Forward)
		WriteBoard(w, board)
	}
}

func Reset(board *snake.Board) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		board.Initialize()
		WriteBoard(w, board)
	}
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

	err := http.ListenAndServe(":8000", nil)
	fmt.Println(err)
}