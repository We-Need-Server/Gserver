package tick

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/game/player"
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/mediator"
	"WeNeedGameServer/network"
	"WeNeedGameServer/packet/client"
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
	Ticks           [60]map[string]player.PlayerPosition
	ActorStatusMap  map[uint32]*ActorStatus
}

type ActorStatus struct {
	Flags       uint8
	UserSEQ     uint32
	RTickNumber uint32
}

func newActorStatus() *ActorStatus {
	return &ActorStatus{}
}

func NewGameTick(tickTime int, game *game.Game, networkInstance *network.Network) *GameTick {
	ticks := [60]map[string]player.PlayerPosition{}
	for i := range ticks {
		ticks[i] = make(map[string]player.PlayerPosition)
	}
	return &GameTick{
		TickTime:        tickTime,
		Ticker:          time.NewTicker(time.Second / time.Duration(tickTime)),
		Game:            game,
		networkInstance: networkInstance,
		stopChan:        make(chan struct{}),
		Mediator:        nil,
		Ticks:           ticks,
		ActorStatusMap:  make(map[uint32]*ActorStatus),
	}
}

func (gt *GameTick) Register(m *mediator.Mediator) {
	gt.Mediator = m
}

func (gt *GameTick) Send(receiverName string, message interface{}) {
	gt.Mediator.Notify("network", receiverName, message)
}

func (gt *GameTick) Receive(senderName string, message interface{}) {
	if senderName == "network" {
		switch pkt := message.(type) {
		case *client.TickIPacket:
			gt.iActorStatus(pkt)
		case *client.TickRPacket:
			gt.rActorStatus(pkt)
		case *internal_type.SEQData:
			gt.updateUserSEQ(pkt)
		}
	} else if senderName == "actor" {
		if val, ok := message.(uint32); ok {
			gt.registerActorStatus(val)
		}
	}
}

func (gt *GameTick) registerActorStatus(qPort uint32) {
	if _, exists := gt.ActorStatusMap[qPort]; !exists {
		gt.ActorStatusMap[qPort] = newActorStatus()
	}
}

func (gt *GameTick) iActorStatus(packet *client.TickIPacket) {
	if _, exists := gt.ActorStatusMap[packet.GetQPort()]; exists {
		gt.ActorStatusMap[packet.GetQPort()].Flags |= 1 << 7
	}
}

func (gt *GameTick) rActorStatus(packet *client.TickRPacket) {
	if _, exists := gt.ActorStatusMap[packet.GetQPort()]; exists {
		gt.ActorStatusMap[packet.GetQPort()].Flags |= 1 << 6
		gt.ActorStatusMap[packet.GetQPort()].RTickNumber = packet.RTickNumber
	}
}

func (gt *GameTick) updateUserSEQ(seqData *internal_type.SEQData) {
	if val, exists := gt.ActorStatusMap[seqData.QPort]; exists && val.UserSEQ+1 == seqData.SEQ {
		gt.ActorStatusMap[seqData.QPort].UserSEQ = seqData.SEQ
	}
}

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

	// 최종적으로 게임 객체의 Delta를 업데이트 하고.
	// 각 유저의 상태에 따라 틱 패킷을 다르게 만들어 보내는 로직을 구현하면 됨.
	// 추가적으로 60틱 저장 로직도 구현해야 함
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
