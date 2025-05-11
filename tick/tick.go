package tick

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/mediator"
	"WeNeedGameServer/network"
	"WeNeedGameServer/packet/server"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type GameTick struct {
	TickTime        int
	Ticker          *time.Ticker
	Game            *game.Game
	networkInstance *network.Network
	stopChan        chan struct{}
	Mediator        *mediator.Mediator
}

func NewGameTick(tickTime int, game *game.Game, networkInstance *network.Network) *GameTick {
	return &GameTick{
		TickTime:        tickTime,
		Ticker:          time.NewTicker(time.Second / time.Duration(tickTime)),
		Game:            game,
		networkInstance: networkInstance,
		stopChan:        make(chan struct{}),
		Mediator:        nil,
	}
}

func (gt *GameTick) Register(m *mediator.Mediator) {
	gt.Mediator = m
}

func (gt *GameTick) Send(receiverName string, message interface{}) {
	gt.Mediator.Notify("network", receiverName, message)
}

func (gt *GameTick) Receive(senderName string, message interface{}) {
}

//func (gt *GameTick) RegisterActorChan(qPort uint32, actorStatusChan chan *ActorStatus) {
//	gt.actorChanMap[qPort] = actorStatusChan
//	//go gt.startActorChan(qPort)
//}

//func (gt *GameTick) startActorChan(qPort uint32) {
//	for {
//		select {
//		//case <-
//		}
//	}
//}

func (gt *GameTick) StartGameLoop() {
	// 루프 시작 틱이 될때마다 processTick함수 실행
	for {
		select {
		case <-gt.Ticker.C:
			gt.processTick()
		case <-gt.stopChan:
			gt.Ticker.Stop()
			return
		}
	}
}

func (gt *GameTick) StopGameLoop() {
	// StartGameLoop가 실행되는 중 어떤 문제가 생길 경우 멈추는 역할
	gt.stopChan <- struct{}{}
}

func (gt *GameTick) processTick() {

	gameState := gt.Game.GetGameState()
	tickPacket := server.TickPacket{
		TickNumber:         gt.TickTime,
		Timestamp:          time.Now(),
		UserSequenceNumber: 0,
		UserPositions:      gameState,
	}

	jsonData, err := json.Marshal(tickPacket)
	if err != nil {
		log.Println("Failed to marshal tickPacket to JSON:", err)
		return
	}

	for _, userAddr := range gt.networkInstance.ConnTable {
		_, err := gt.networkInstance.Ln.WriteToUDP(jsonData, userAddr)
		if err != nil {
			log.Println("Failed to send message:", err)
		}
	}

	gt.TickTime += 1
	fmt.Println("Game state sent to", len(gt.networkInstance.ConnTable), "clients")
}
