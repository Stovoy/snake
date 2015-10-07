package main

import (
	"snake/server"
	"snake/snake"
)

func main() {
	board := snake.NewBoard(50, 50)

	server.Start(board)
}
