package tick

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/game/player"
	udp_client2 "WeNeedGameServer/protocol/udp/udp_client"
	udp_server2 "WeNeedGameServer/protocol/udp/udp_server"
	"WeNeedGameServer/round/internal/sender"
	"WeNeedGameServer/round/internal/types"
	"fmt"
	"log"
	"time"
)

type GameTick struct {
	tickTime          uint32
	ticker            *time.Ticker
	game              *game.Game
	udpSender         *sender.UdpSender
	ticks             [60]map[uint32]*player.PlayerPosition
	actorStatusMap    map[uint32]*ActorStatus
	stopPacket        *udp_server2.StopPacket
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

func NewGameTick(tickTime int64, game *game.Game, udpSender *sender.UdpSender) *GameTick {
	ticks := [60]map[uint32]*player.PlayerPosition{}
	for i := range ticks {
		ticks[i] = make(map[uint32]*player.PlayerPosition)
	}
	return &GameTick{
		tickTime:          0,
		ticker:            time.NewTicker(time.Second / time.Duration(tickTime)),
		game:              game,
		udpSender:         udpSender,
		ticks:             ticks,
		actorStatusMap:    make(map[uint32]*ActorStatus),
		stopPacket:        udp_server2.NewStopPacket(),
		playerPositionMap: nil,
	}
}

func (gt *GameTick) registerActorStatus(qPort uint32) {
	if _, exists := gt.actorStatusMap[qPort]; !exists {
		gt.actorStatusMap[qPort] = newActorStatus()
	}
}

func (gt *GameTick) iActorStatus(packet *udp_client2.TickIPacket) {
	gt.registerActorStatus(packet.GetQPort())
	gt.actorStatusMap[packet.GetQPort()].Flags = 1 << 7
}

func (gt *GameTick) rActorStatus(packet *udp_client2.TickRPacket) {
	gt.registerActorStatus(packet.GetQPort())
	gt.actorStatusMap[packet.GetQPort()].Flags = 1 << 6
	gt.actorStatusMap[packet.GetQPort()].RTickNumber = packet.RTickNumber
	fmt.Println("비트 테스트")
	fmt.Println(gt.actorStatusMap[packet.GetQPort()].Flags)
}

func (gt *GameTick) updateUserSEQ(seqData *types.SEQData) {
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
		p := <-*gt.udpSender.NChan
		switch p.GetPacketKind() {
		case 'S':
			tempMap := playerPositionMap
			playerPositionMap = make(map[uint32]*player.PlayerPosition)
			gt.playerPositionMap = &tempMap
			break
		case 'I':
			if p, ok := p.(*udp_client2.TickIPacket); ok {
				gt.iActorStatus(p)
			}
			break
		case 'R':
			fmt.Println("왔다잉 재전송 패킷")
			if p, ok := p.(*udp_client2.TickRPacket); ok {
				gt.rActorStatus(p)
			}
			break
		case 'D':
			fmt.Println("delta")
			if _, exists := playerPositionMap[p.GetQPort()]; !exists {
				playerPositionMap[p.GetQPort()] = player.NewPlayerPositionD()
			}
			if p, ok := p.(*udp_server2.DeltaPacket); ok {
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
	*gt.udpSender.NChan <- udp_server2.NewStopPacket()
	for gt.playerPositionMap == nil {
		//fmt.Println("while", gt.playerPositionMap)
	}
	//fmt.Println("out", *gt.playerPositionMap)
	gt.ticks[gt.tickTime%60] = *gt.playerPositionMap
	gt.game.ReflectPlayers(gt.playerPositionMap)
	//fmt.Println("out2", *gt.playerPositionMap)
	gameState := gt.game.GetGameState()

	for qPort, userAddr := range *gt.udpSender.ConnTable {
		gt.registerActorStatus(qPort)
		actorStatus := gt.actorStatusMap[qPort]
		var tickPacket *udp_server2.TickPacket
		if (actorStatus.Flags & (1 << 7)) != 0 {
			tickPacket = udp_server2.NewTickPacket(gt.tickTime, time.Now().Unix(), (*gt.udpSender.NextSeqTable)[qPort]-1, actorStatus.Flags, gameState)
		} else if (actorStatus.Flags & (1 << 6)) != 0 {
			restoreTickCount := gt.tickTime - actorStatus.RTickNumber
			if restoreTickCount >= 60 {
				fmt.Println("좌표값 패킷 발사")
				actorStatus.Flags = (actorStatus.Flags &^ (1 << 6)) | (1 << 7)
				tickPacket = udp_server2.NewTickPacket(gt.tickTime, time.Now().Unix(), (*gt.udpSender.NextSeqTable)[qPort]-1, actorStatus.Flags, gameState)
			} else {
				fmt.Println("재전송 패킷 발사")
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
				tickPacket = udp_server2.NewTickPacket(gt.tickTime, time.Now().Unix(), (*gt.udpSender.NextSeqTable)[qPort]-1, actorStatus.Flags, cloneGameDeltaState)
			}
		} else {
			fmt.Println("tick packet", *((*gt).playerPositionMap))
			tickPacket = udp_server2.NewTickPacket(gt.tickTime, time.Now().Unix(), (*gt.udpSender.NextSeqTable)[qPort]-1, actorStatus.Flags, *gt.playerPositionMap)
		}

		_, err := gt.udpSender.SendUdpPacket(tickPacket.Serialize(), userAddr)
		if err != nil {
			log.Println("Failed to send message:", err)
		}

		actorStatus.Flags = 0
		actorStatus.RTickNumber = 0
	}
	//fmt.Println("Game state sent to", len(*gt.udpSender.ConnTable), "clients")
	gt.playerPositionMap = nil
	gt.tickTime += 1
}

//func (gt *GameTick) processTick() {
//	gameDeltaState := gt.game.GetGameDeltaState()
//
//	gt.ticks[gt.tickTime%60] = gameDeltaState
//	gameState := gt.game.GetGameState()
//	for qPort, userAddr := range *gt.udpSender.ConnTable {
//		actorStatus := gt.actorStatusMap[qPort]
//		var tickPacket *udp_server.TickPacket
//		if (actorStatus.Flags & 1 << 7) != 0 {
//			tickPacket = udp_server.NewTickPacket(gt.tickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameState)
//		} else if (actorStatus.Flags & 1 << 6) != 0 {
//			restoreTickCount := gt.tickTime - actorStatus.RTickNumber
//			if restoreTickCount >= 60 {
//				tickPacket = udp_server.NewTickPacket(gt.tickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags&^(1<<6), gameState)
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
//				tickPacket = udp_server.NewTickPacket(gt.tickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameState)
//			}
//		} else {
//			tickPacket = udp_server.NewTickPacket(gt.tickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameDeltaState)
//		}
//
//		_, err := gt.udpSender.SendUdpPacket(tickPacket.Serialize(), userAddr)
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
