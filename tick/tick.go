package tick

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/game/player"
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/network/sender"
	"WeNeedGameServer/packet/client"
	"WeNeedGameServer/packet/server"
	"fmt"
	"log"
	"time"
)

type GameTick struct {
	tickTime          uint32
	ticker            *time.Ticker
	game              *game.Game
	us                *sender.Sender
	ticks             [60]map[uint32]*player.PlayerPosition
	actorStatusMap    map[uint32]*ActorStatus
	stopPacket        *server.StopPacket
	playerPositionMap *map[uint32]*player.PlayerPosition
}

type ActorStatus struct {
	Flags       uint8
	UserSEQ     uint32
	RTickNumber uint32
}

func newActorStatus() *ActorStatus {
	return &ActorStatus{}
}

func NewGameTick(tickTime int64, game *game.Game, us *sender.Sender) *GameTick {
	ticks := [60]map[uint32]*player.PlayerPosition{}
	for i := range ticks {
		ticks[i] = make(map[uint32]*player.PlayerPosition)
	}
	return &GameTick{
		tickTime:          0,
		ticker:            time.NewTicker(time.Second / time.Duration(tickTime)),
		game:              game,
		us:                us,
		ticks:             ticks,
		actorStatusMap:    make(map[uint32]*ActorStatus),
		stopPacket:        server.NewStopPacket(),
		playerPositionMap: nil,
	}
}

func (gt *GameTick) registerActorStatus(qPort uint32) {
	if _, exists := gt.actorStatusMap[qPort]; !exists {
		gt.actorStatusMap[qPort] = newActorStatus()
	}
}

func (gt *GameTick) iActorStatus(packet *client.TickIPacket) {
	gt.registerActorStatus(packet.GetQPort())
	gt.actorStatusMap[packet.GetQPort()].Flags |= 1 << 7
}

func (gt *GameTick) rActorStatus(packet *client.TickRPacket) {
	gt.registerActorStatus(packet.GetQPort())
	gt.actorStatusMap[packet.GetQPort()].Flags |= 1 << 6
	gt.actorStatusMap[packet.GetQPort()].RTickNumber = packet.RTickNumber
}

func (gt *GameTick) updateUserSEQ(seqData *internal_type.SEQData) {
	if val, exists := gt.actorStatusMap[seqData.QPort]; exists && val.UserSEQ+1 == seqData.SEQ {
		gt.actorStatusMap[seqData.QPort].UserSEQ = seqData.SEQ
	}
}

func (gt *GameTick) StartGameLoop() {
	go gt.dequeuePacket()
	// 루프 시작 틱이 될때마다 processTick함수 실행
	for {
		select {
		case <-gt.ticker.C:
			gt.processTick()
		}
	}
}

func (gt *GameTick) dequeuePacket() {
	fmt.Println()
	playerPositionMap := make(map[uint32]*player.PlayerPosition)
	for {
		p := <-*gt.us.NChan
		switch p.GetPacketKind() {
		case 'S':
			tempMap := playerPositionMap
			playerPositionMap = make(map[uint32]*player.PlayerPosition)
			gt.playerPositionMap = &tempMap
			break
		case 'I':
			if p, ok := p.(*client.TickIPacket); ok {
				gt.iActorStatus(p)
			}
			break
		case 'R':
			if p, ok := p.(*client.TickRPacket); ok {
				gt.rActorStatus(p)
			}
			break
		case 'D':
			fmt.Println("delta")
			if _, exists := playerPositionMap[p.GetQPort()]; !exists {
				playerPositionMap[p.GetQPort()] = player.NewPlayerPositionD()
			}
			if p, ok := p.(*server.DeltaPacket); ok {
				playerPositionMap[p.GetQPort()].CalculatePlayerPosition(p.PlayerPosition)
				fmt.Println(*playerPositionMap[p.GetQPort()])
				for key, val := range *p.HitInformationMap {
					if _, exists := playerPositionMap[p.GetQPort()]; !exists {
						playerPositionMap[key] = player.NewPlayerPositionD()
					}
					playerPositionMap[key].Hp += val
				}

			}
			break
		}
	}
}

func (gt *GameTick) processTick() {
	*gt.us.NChan <- server.NewStopPacket()
	for gt.playerPositionMap == nil {
		fmt.Println("while", gt.playerPositionMap)
	}
	fmt.Println("out", *gt.playerPositionMap)
	gt.ticks[gt.tickTime%60] = *gt.playerPositionMap
	gt.game.ReflectPlayers(gt.playerPositionMap)
	fmt.Println("out2", *gt.playerPositionMap)
	gameState := gt.game.GetGameState()

	for qPort, userAddr := range *gt.us.ConnTable {
		gt.registerActorStatus(qPort)
		actorStatus := gt.actorStatusMap[qPort]
		var tickPacket *server.TickPacket
		if (actorStatus.Flags & 1 << 7) != 0 {
			tickPacket = server.NewTickPacket(gt.tickTime, time.Now().Unix(), (*gt.us.NextSeqTable)[qPort]-1, actorStatus.Flags, gameState)
		} else if (actorStatus.Flags & 1 << 6) != 0 {
			restoreTickCount := gt.tickTime - actorStatus.RTickNumber
			if restoreTickCount >= 60 {
				actorStatus.Flags = (actorStatus.Flags &^ (1 << 6)) | (1 << 7)
				tickPacket = server.NewTickPacket(gt.tickTime, time.Now().Unix(), (*gt.us.NextSeqTable)[qPort]-1, actorStatus.Flags, gameState)
			} else {
				cloneGameDeltaState := make(map[uint32]*player.PlayerPosition)
				for k, v := range *gt.playerPositionMap {
					cloneGameDeltaState[k] = v
				}
				for i := actorStatus.RTickNumber; i < gt.tickTime; i++ {
					tickIdx := i % 60
					for qPort, playerPosition := range gt.ticks[tickIdx] {
						if pos, exists := cloneGameDeltaState[qPort]; exists {
							pos.Hp += playerPosition.Hp
							pos.PositionX += playerPosition.PositionX
							pos.PositionZ += playerPosition.PositionZ
							pos.PtAngle += playerPosition.PtAngle
							pos.YawAngle += playerPosition.YawAngle
							pos.Jp = playerPosition.Jp
							pos.IsShoot = playerPosition.IsShoot
							cloneGameDeltaState[qPort] = pos
						}
					}
				}
				tickPacket = server.NewTickPacket(gt.tickTime, time.Now().Unix(), (*gt.us.NextSeqTable)[qPort]-1, actorStatus.Flags, cloneGameDeltaState)
			}
		} else {
			fmt.Println("tick packet", *((*gt).playerPositionMap))
			tickPacket = server.NewTickPacket(gt.tickTime, time.Now().Unix(), (*gt.us.NextSeqTable)[qPort]-1, actorStatus.Flags, *gt.playerPositionMap)
		}

		_, err := gt.us.SendUDPPacket(tickPacket.Serialize(), userAddr)
		if err != nil {
			log.Println("Failed to send message:", err)
		}

		actorStatus.Flags = 0
		actorStatus.RTickNumber = 0
	}
	fmt.Println("Game state sent to", len(*gt.us.ConnTable), "clients")
	gt.playerPositionMap = nil
	gt.tickTime += 1
}

//func (gt *GameTick) processTick() {
//	gameDeltaState := gt.game.GetGameDeltaState()
//
//	gt.ticks[gt.tickTime%60] = gameDeltaState
//	gameState := gt.game.GetGameState()
//	for qPort, userAddr := range *gt.us.ConnTable {
//		actorStatus := gt.actorStatusMap[qPort]
//		var tickPacket *server.TickPacket
//		if (actorStatus.Flags & 1 << 7) != 0 {
//			tickPacket = server.NewTickPacket(gt.tickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameState)
//		} else if (actorStatus.Flags & 1 << 6) != 0 {
//			restoreTickCount := gt.tickTime - actorStatus.RTickNumber
//			if restoreTickCount >= 60 {
//				tickPacket = server.NewTickPacket(gt.tickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags&^(1<<6), gameState)
//			} else {
//				cloneGameDeltaState := make(map[uint32]player.PlayerPosition)
//				for k, v := range gameDeltaState {
//					cloneGameDeltaState[k] = v
//				}
//				for i := actorStatus.RTickNumber; i < gt.tickTime; i++ {
//					tickIdx := i % 60
//					for qPort, playerPosition := range gt.ticks[tickIdx] {
//						if pos, exists := cloneGameDeltaState[qPort]; exists {
//							pos.Hp += playerPosition.Hp
//							pos.PositionX += playerPosition.PositionX
//							pos.PositionZ += playerPosition.PositionZ
//							pos.PtAngle += playerPosition.PtAngle
//							pos.YawAngle += playerPosition.YawAngle
//							pos.Jp = playerPosition.Jp
//							pos.IsShoot = playerPosition.IsShoot
//							cloneGameDeltaState[qPort] = pos
//						}
//					}
//				}
//				tickPacket = server.NewTickPacket(gt.tickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameState)
//			}
//		} else {
//			tickPacket = server.NewTickPacket(gt.tickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameDeltaState)
//		}
//
//		_, err := gt.us.SendUDPPacket(tickPacket.Serialize(), userAddr)
//		if err != nil {
//			log.Println("Failed to send message:", err)
//		}
//
//		actorStatus.Flags = 0
//		actorStatus.RTickNumber = 0
//	}
//	gt.game.ResetHPDelta()
//	gt.tickTime += 1
//	//fmt.Println("Game state sent to", len(gt.networkInstance.ConnTable), "clients")
//}
