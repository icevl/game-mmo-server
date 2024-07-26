package main

import (
	"server/gameserver"
	"server/http"
)

const (
	WorldFilePath = "/Users/ice/MMO/Assets/Editor/level.txt"
)

func main() {
	ch := make(chan int)

	go gameserver.StartGameServer()
	go http.Start()

	ch <- 1
}
