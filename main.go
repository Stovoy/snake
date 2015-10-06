package main

import (
	"fmt"
	"snake/server"
	"snake/snake"
)

func main() {
	board := snake.NewBoard(5, 5)
	fmt.Println(board)

	server.Start(board)
}
