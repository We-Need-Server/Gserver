package main

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/lobby"
)

//var PacketRegisterInstance = make(register.PacketRegister)

//func main() {
//	mediatorInstance := mediator.NewMediator()
//	gameInstance := legacy3.NewGame()
//	networkInstance := legacy.NewNetwork(gameInstance)
//	if _, err := mediatorInstance.Register("network", networkInstance); err != nil {
//		log.Panicln("메디에이터 등록 실패")
//	}
//	go networkInstance.Start()
//	gameTickInstance := legacy2.NewGameTick(0, gameInstance, networkInstance)
//	if _, err := mediatorInstance.Register("game_tick", gameTickInstance); err != nil {
//		log.Panicln("메디에이터 등록 실패")
//	}
//	gameTickInstance.StartGameLoop()
//}

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

//
//1. 게임 객체가 만들어짐
//2. 네트워크 객체가 만들어짐
//3. 게임틱 객체가 만들어짐
//
//네트워크 객체가 연결을 받아서 이제 새로운 연결이면 액터를 만든다.
//	이 액터의 역할은 네트워크 객체가 해석한 패킷(이벤트 패킷인지 틱패킷인지 해석)을 받아서 처리하는 거란 말이야.
//	N개의 액터를 만들면, 락킹을 안해도 된다.
//
//	틱 인스턴스는 게임 객체의 정보를 추출해서 그걸 1초에 60번 보내줌
//
//
//네트워크 -> 액터 -> 게임
//		게임 <- 틱
