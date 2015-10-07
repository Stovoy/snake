package main

import (
	"snake/server"
	"snake/snake"
)

func main() {
	board := snake.NewBoard(8, 8)

	server.Start(board)
}
