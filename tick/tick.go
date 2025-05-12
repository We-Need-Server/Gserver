package tick

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/game/player"
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/mediator"
	"WeNeedGameServer/network"
	"WeNeedGameServer/packet/client"
	"WeNeedGameServer/packet/server"
	"fmt"
	"log"
	"time"
)

type GameTick struct {
	TickTime        uint32
	Ticker          *time.Ticker
	Game            *game.Game
	networkInstance *network.Network
	stopChan        chan struct{}
	Mediator        *mediator.Mediator
	Ticks           [60]map[uint32]player.PlayerPosition
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

func NewGameTick(tickTime uint32, game *game.Game, networkInstance *network.Network) *GameTick {
	ticks := [60]map[uint32]player.PlayerPosition{}
	for i := range ticks {
		ticks[i] = make(map[uint32]player.PlayerPosition)
	}
	return &GameTick{
		TickTime:        tickTime,
		Ticker:          time.NewTicker(time.Second / time.Duration(10)),
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
	gameDeltaState := gt.Game.GetGameDeltaState()
	gt.Ticks[gt.TickTime%60] = gameDeltaState
	gameState := gt.Game.GetGameState()
	for qPort, userAddr := range gt.networkInstance.ConnTable {
		actorStatus := gt.ActorStatusMap[qPort]
		var tickPacket *server.TickPacket
		if (actorStatus.Flags & 1 << 7) != 0 {
			tickPacket = server.NewTickPacket(gt.TickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameState)
		} else if (actorStatus.Flags & 1 << 6) != 0 {
			restoreTickCount := gt.TickTime - actorStatus.RTickNumber
			if restoreTickCount >= 60 {
				tickPacket = server.NewTickPacket(gt.TickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags&^(1<<6), gameState)
			} else {
				cloneGameDeltaState := make(map[uint32]player.PlayerPosition)
				for k, v := range gameDeltaState {
					cloneGameDeltaState[k] = v
				}
				for i := actorStatus.RTickNumber; i < gt.TickTime; i++ {
					tickIdx := i % 60
					for qPort, playerPosition := range gt.Ticks[tickIdx] {
						if pos, exists := cloneGameDeltaState[qPort]; exists {
							pos.PositionX += playerPosition.PositionX
							pos.PositionZ += playerPosition.PositionZ
							pos.PTAngle += playerPosition.PTAngle
							pos.YawAngle += playerPosition.YawAngle
							pos.JP = playerPosition.JP
							cloneGameDeltaState[qPort] = pos
						}
					}
				}
				tickPacket = server.NewTickPacket(gt.TickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameState)
			}
		} else {
			tickPacket = server.NewTickPacket(gt.TickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameDeltaState)
		}

		_, err := gt.networkInstance.Ln.WriteToUDP(tickPacket.Serialize(), userAddr)
		if err != nil {
			log.Println("Failed to send message:", err)
		}

		actorStatus.Flags = 0
		actorStatus.RTickNumber = 0
	}

	gt.TickTime += 1
	fmt.Println("Game state sent to", len(gt.networkInstance.ConnTable), "clients")
}
