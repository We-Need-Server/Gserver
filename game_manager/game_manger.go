package game_manager

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/protocol/tcp"
	"WeNeedGameServer/protocol/tcp/tserver"
)

type GameManager struct {
	userDb               *db.Db
	matchScore           uint16 // 게임이 총 몇 판 몇 선제일때의 몇 판
	blueScore            uint16 // 라운드 승리 횟수
	redScore             uint16
	blueAlivePlayerCount uint16
	redAlivePlayerCount  uint16
	finalWinnerTeam      uint8
	sendTcpPacketFunc    func(message *tcp.ReceiverMessage)
}

func NewGameManager(userDb *db.Db, matchScore uint16, sendTcpPacketFunc func(message *tcp.ReceiverMessage)) *GameManager {
	return &GameManager{
		userDb:               userDb,
		matchScore:           matchScore,
		blueScore:            0,
		redScore:             0,
		blueAlivePlayerCount: 0,
		redAlivePlayerCount:  0,
		finalWinnerTeam:      0,
		sendTcpPacketFunc:    sendTcpPacketFunc,
	}
}

func (gm *GameManager) InitAlivePlayer() {
	gm.blueAlivePlayerCount = gm.userDb.GetTeamAlivePlayerCount(db.BlueTeam)
	gm.redAlivePlayerCount = gm.userDb.GetTeamAlivePlayerCount(db.RedTeam)
}

func (gm *GameManager) IncreaseScore(winnerTeam db.Team) {
	gm.matchScore -= 1
	if winnerTeam {
		gm.redScore += 1
	} else {
		gm.blueScore += 1
	}
	if gm.matchScore == 0 {
		gm.sendTcpPacketFunc(tcp.NewBroadCastMessage('O', tserver.NewGameOverPacket()))
	} else {

	}
}

func (gm *GameManager) DecreasePlayer(deadPlayerTeam db.Team) {
	if deadPlayerTeam {
		gm.redAlivePlayerCount -= 1
		if gm.redAlivePlayerCount == 0 {
			gm.IncreaseScore(db.BlueTeam)
		}
	} else {
		gm.blueAlivePlayerCount -= 1
		if gm.blueAlivePlayerCount == 0 {
			gm.IncreaseScore(db.RedTeam)
		}
	}
}
