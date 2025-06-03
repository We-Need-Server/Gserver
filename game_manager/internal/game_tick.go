package internal

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/game_manager/internal/internal_types"
	"WeNeedGameServer/game_manager/internal/sender"
	"WeNeedGameServer/game_type"
	"WeNeedGameServer/protocol/udp/uclient"
	"WeNeedGameServer/protocol/udp/userver"
	"fmt"
	"log"
	"time"
)

type GameTick struct {
	TickTime       uint32
	ticker         *time.Ticker
	game           *game.Game
	udpSender      *sender.UdpSender
	ticks          [60]map[uint32]*game_type.PlayerState
	actorStatusMap map[uint32]*ActorStatus
	stopPacket     *userver.StopPacket
	playerStateMap map[uint32]*game_type.PlayerState
	findUserFunc   func(uint32) bool
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
	ticks := [60]map[uint32]*game_type.PlayerState{}
	for i := range ticks {
		ticks[i] = make(map[uint32]*game_type.PlayerState)
	}
	return &GameTick{
		TickTime:       0,
		ticker:         time.NewTicker(time.Second / time.Duration(tickTime)),
		game:           game,
		udpSender:      udpSender,
		ticks:          ticks,
		actorStatusMap: make(map[uint32]*ActorStatus),
		stopPacket:     userver.NewStopPacket(),
		playerStateMap: nil,
		findUserFunc:   findUserFunc,
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
	playerStateMap := make(map[uint32]*game_type.PlayerState)
	for {
		p := <-gt.udpSender.NChan
		switch p.GetPacketKind() {
		case 'S':
			tempMap := playerStateMap
			playerStateMap = make(map[uint32]*game_type.PlayerState)
			gt.playerStateMap = tempMap
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
			if _, exists := playerStateMap[p.GetQPort()]; !exists {
				playerStateMap[p.GetQPort()] = game_type.NewPlayerStateDefault()
			}
			if p, ok := p.(*userver.DeltaPacket); ok {
				playerStateMap[p.GetQPort()].CalculatePlayerState(p.PlayerPosition)
				for key, val := range *p.HitInformationMap {
					if _, exists := playerStateMap[p.GetQPort()]; !exists {
						playerStateMap[key] = game_type.NewPlayerStateDefault()
					}
					playerStateMap[key].Damage += val
				}
			}
			break
		}
	}
}

func (gt *GameTick) processTick() {
	gt.udpSender.NChan <- userver.NewStopPacket()
	for gt.playerStateMap == nil {
		//fmt.Println("while", gt.playerStateMap)
	}
	//fmt.Println("out", *gt.playerStateMap)
	gt.ticks[gt.TickTime%60] = gt.playerStateMap
	gt.game.ReflectPlayers(gt.playerStateMap)
	//fmt.Println("out2", *gt.playerStateMap)
	// 여기서 종료까지 같은 스레드에서 해버리는게 문제
	// 이거 때문에 이전 라운드에 대한 정보를 다음 라운드가 가져가버림
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
					cloneGameDeltaState := make(map[uint32]*game_type.PlayerState)
					for k, v := range gt.playerStateMap {
						cloneGameDeltaState[k] = v
					}
					for i := actorStatus.RTickNumber; i < gt.TickTime; i++ {
						tickIdx := i % 60
						for userId, playerState := range gt.ticks[tickIdx] {
							if pos, exists := cloneGameDeltaState[userId]; exists {
								pos.Damage += playerState.Damage
								pos.PositionX += playerState.PositionX
								pos.PositionZ += playerState.PositionZ
								pos.PtAngle += playerState.PtAngle
								pos.YawAngle += playerState.YawAngle
								pos.Jp = playerState.Jp
								pos.IsShoot = playerState.IsShoot
								cloneGameDeltaState[userId] = pos
							}
						}
					}
					tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), gt.udpSender.NextSeqTable[qPort]-1, actorStatus.Flags, cloneGameDeltaState)
				}
			} else {
				fmt.Println("game_tick packet", gt.playerStateMap)
				tickPacket = userver.NewTickPacket(gt.TickTime, time.Now().Unix(), gt.udpSender.NextSeqTable[qPort]-1, actorStatus.Flags, gt.playerStateMap)
			}
			_, err := gt.udpSender.SendUdpPacket(tickPacket.Serialize(), userConnStatus.Conn)
			if err != nil {
				log.Println("Failed to send message:", err)
			}
			actorStatus.Flags = 0
			actorStatus.RTickNumber = 0
		} else {
			// 여기서는 게임 인스턴스를 죽이는 것보다는 connTable에서 제거하는 게 낫나?
			gt.game.DeletePlayer(userConnStatus.UserId)
		}
	}
	//fmt.Println("Game state sent to", len(*gt.udpSender.ConnTable), "clients")
	gt.playerStateMap = nil
	gt.TickTime += 1
}
