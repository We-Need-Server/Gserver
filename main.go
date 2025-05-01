package main

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/network"
	"WeNeedGameServer/tick"
)

func main() {
	gameInstance := game.NewGame()
	networkInstance := network.NewNetwork(gameInstance)
	go networkInstance.Start()
	gameTickInstance := tick.NewGameTick(1, gameInstance, networkInstance)
	gameTickInstance.StartGameLoop()
}
