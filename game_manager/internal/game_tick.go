package internal

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/game/player"
	"WeNeedGameServer/game_manager/internal/internal_types"
	"WeNeedGameServer/game_manager/internal/sender"
	"WeNeedGameServer/protocol/udp/uclient"
	"WeNeedGameServer/protocol/udp/userver"
	"fmt"
	"log"
	"time"
)

type GameTick struct {
	TickTime          uint32
	ticker            *time.Ticker
	game              *game.Game
	udpSender         *sender.UdpSender
	ticks             [60]map[uint32]*player.PlayerPosition
	actorStatusMap    map[uint32]*ActorStatus
	stopPacket        *userver.StopPacket
	playerPositionMap map[uint32]*player.PlayerPosition
	findUserFunc      func(uint32) bool
}

type ActorStatus struct {
	Flags       uint8
	UserSEQ     uint32
	RTickNumber uint32
}

func newActorStatus() *ActorStatus {
	return &ActorStatus{}
}

func NewGameTick(tickTime int64, game *game.Game, udpSender *sender.UdpSender, findUserFunc func(uint32) bool) *GameTick {
	ticks := [60]map[uint32]*player.PlayerPosition{}
	for i := range ticks {
		ticks[i] = make(map[uint32]*player.PlayerPosition)
	}
	return &GameTick{
		TickTime:          0,
		ticker:            time.NewTicker(time.Second / time.Duration(tickTime)),
		game:              game,
		udpSender:         udpSender,
		ticks:             ticks,
		actorStatusMap:    make(map[uint32]*ActorStatus),
		stopPacket:        userver.NewStopPacket(),
		playerPositionMap: nil,
		findUserFunc:      findUserFunc,
	}
}

func (gt *GameTick) registerActorStatus(qPort uint32) {
	if _, exists := gt.actorStatusMap[qPort]; !exists {
		gt.actorStatusMap[qPort] = newActorStatus()
	}
}

func (gt *GameTick) iActorStatus(packet *uclient.TickIPacket) {
	gt.registerActorStatus(packet.GetQPort())
	gt.actorStatusMap[packet.GetQPort()].Flags = 1 << 7
}

func (gt *GameTick) rActorStatus(packet *uclient.TickRPacket) {
	gt.registerActorStatus(packet.GetQPort())
	gt.actorStatusMap[packet.GetQPort()].Flags = 1 << 6
	gt.actorStatusMap[packet.GetQPort()].RTickNumber = packet.RTickNumber
	fmt.Println("비트 테스트")
	fmt.Println(gt.actorStatusMap[packet.GetQPort()].Flags)
}

func (gt *GameTick) updateUserSEQ(seqData *internal_types.SEQData) {
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
		p := <-gt.udpSender.NChan
		switch p.GetPacketKind() {
		case 'S':
			tempMap := playerPositionMap
			playerPositionMap = make(map[uint32]*player.PlayerPosition)
			gt.playerPositionMap = tempMap
			break
		case 'I':
			if p, ok := p.(*uclient.TickIPacket); ok {
				gt.iActorStatus(p)
			}
			break
		case 'R':
			fmt.Println("왔다잉 재전송 패킷")
			if p, ok := p.(*uclient.TickRPacket); ok {
				gt.rActorStatus(p)
			}
			break
		case 'D':
			fmt.Println("delta")
			if _, exists := playerPositionMap[gt.udpSender.ConnTable[p.GetQPort()].UserId]; !exists {
				playerPositionMap[gt.udpSender.ConnTable[p.GetQPort()].UserId] = player.NewPlayerPositionD()
			}
			if p, ok := p.(*userver.DeltaPacket); ok {
				playerPositionMap[gt.udpSender.ConnTable[p.GetQPort()].UserId].CalculatePlayerPosition(p.PlayerPosition)
				fmt.Println(*playerPositionMap[gt.udpSender.ConnTable[p.GetQPort()].UserId])
				for key, val := range *p.HitInformationMap {
					if _, exists := playerPositionMap[gt.udpSender.ConnTable[p.GetQPort()].UserId]; !exists {
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
	gt.udpSender.NChan <- userver.NewStopPacket()
	for gt.playerPositionMap == nil {
		//fmt.Println("while", gt.playerPositionMap)
	}
	//fmt.Println("out", *gt.playerPositionMap)
	gt.ticks[gt.TickTime%60] = gt.playerPositionMap
	gt.game.ReflectPlayers(gt.playerPositionMap)
	//fmt.Println("out2", *gt.playerPositionMap)
	gameState := gt.game.GetGameState()

	for qPort, userConnStatus := range gt.udpSender.ConnTable {
		if gt.findUserFunc(userConnStatus.UserId) {
			gt.registerActorStatus(qPort)
			actorStatus := gt.actorStatusMap[qPort]
			var tickPacket *userver.TickPacket
			if (actorStatus.Flags & (1 << 7)) != 0 {
				tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), gt.udpSender.NextSeqTable[qPort]-1, actorStatus.Flags, gameState)
			} else if (actorStatus.Flags & (1 << 6)) != 0 {
				restoreTickCount := gt.TickTime - actorStatus.RTickNumber
				if restoreTickCount >= 60 {
					fmt.Println("좌표값 패킷 발사")
					actorStatus.Flags = (actorStatus.Flags &^ (1 << 6)) | (1 << 7)
					tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), gt.udpSender.NextSeqTable[qPort]-1, actorStatus.Flags, gameState)
				} else {
					fmt.Println("재전송 패킷 발사")
					cloneGameDeltaState := make(map[uint32]*player.PlayerPosition)
					for k, v := range gt.playerPositionMap {
						cloneGameDeltaState[k] = v
					}
					for i := actorStatus.RTickNumber; i < gt.TickTime; i++ {
						tickIdx := i % 60
						for userId, playerPosition := range gt.ticks[tickIdx] {
							if pos, exists := cloneGameDeltaState[userId]; exists {
								pos.Hp += playerPosition.Hp
								pos.PositionX += playerPosition.PositionX
								pos.PositionZ += playerPosition.PositionZ
								pos.PtAngle += playerPosition.PtAngle
								pos.YawAngle += playerPosition.YawAngle
								pos.Jp = playerPosition.Jp
								pos.IsShoot = playerPosition.IsShoot
								cloneGameDeltaState[userId] = pos
							}
						}
					}
					tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), gt.udpSender.NextSeqTable[qPort]-1, actorStatus.Flags, cloneGameDeltaState)
				}
			} else {
				fmt.Println("game_tick packet", gt.playerPositionMap)
				tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), gt.udpSender.NextSeqTable[qPort]-1, actorStatus.Flags, gt.playerPositionMap)
			}
			fmt.Println("tick")
			fmt.Println(tickPacket.Serialize())
			fmt.Println(userConnStatus.Conn)
			_, err := gt.udpSender.SendUdpPacket(tickPacket.Serialize(), userConnStatus.Conn)
			if err != nil {
				log.Println("Failed to send message:", err)
			}
			actorStatus.Flags = 0
			actorStatus.RTickNumber = 0
		} else {
			gt.game.DeletePlayer(userConnStatus.UserId)
		}
	}
	//fmt.Println("Game state sent to", len(*gt.udpSender.ConnTable), "clients")
	gt.playerPositionMap = nil
	gt.TickTime += 1
}

//func (gt *GameTick) processTick() {
//	gameDeltaState := gt.game.GetGameDeltaState()
//
//	gt.ticks[gt.TickTime%60] = gameDeltaState
//	gameState := gt.game.GetGameState()
//	for qPort, userAddr := range *gt.udpSender.ConnTable {
//		actorStatus := gt.actorStatusMap[qPort]
//		var tickPacket *userver.TickPacket
//		if (actorStatus.Flags & 1 << 7) != 0 {
//			tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameState)
//		} else if (actorStatus.Flags & 1 << 6) != 0 {
//			restoreTickCount := gt.TickTime - actorStatus.RTickNumber
//			if restoreTickCount >= 60 {
//				tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags&^(1<<6), gameState)
//			} else {
//				cloneGameDeltaState := make(map[uint32]player.PlayerPosition)
//				for k, v := range gameDeltaState {
//					cloneGameDeltaState[k] = v
//				}
//				for i := actorStatus.RTickNumber; i < gt.TickTime; i++ {
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
//				tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameState)
//			}
//		} else {
//			tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), actorStatus.UserSEQ, actorStatus.Flags, gameDeltaState)
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
//	gt.TickTime += 1
//	//fmt.Println("Game state sent to", len(gt.networkInstance.ConnTable), "clients")
//}
