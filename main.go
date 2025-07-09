package main

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/lobby"
)

func main() {
	userDbInstance := db.NewUserDb()
	userDbInstance.Init()
	lobbyInstance := lobby.NewLobby(userDbInstance, ":20001", ":20000", 10)
	tcpReceiver, _ := lobbyInstance.ReadyTcp()
	tcpReceiver.StartTcp()
	//userDbInstance := db.NewUserDb()
	//userDbInstance.Init()
	//networkInstance := internal.NewNetwork(":20000", ":20001")
	//udpReceiver, udpSender := networkInstance.ReadyUdp()
	//go udpReceiver.StartUdp()
	//gameInstance := game.NewGame()
	//tickInstance := tick.NewGameTick(60, gameInstance, udpSender)
	//tickInstance.StartGameLoop()
}
