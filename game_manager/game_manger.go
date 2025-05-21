package game_manager

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game"
	"WeNeedGameServer/protocol/tcp"
	"WeNeedGameServer/protocol/tcp/tserver"
	"WeNeedGameServer/util"
)

// 스냅샷 방식으로 살아있는 플레이어 수를 정한다.
type GameManager struct {
	userSpawnPositionArr []int
	userDb               *db.Db
	matchScore           uint16 // 게임이 총 몇 판 몇 선제일때의 몇 판
	blueScore            uint16 // 라운드 승리 횟수
	redScore             uint16
	blueAlivePlayerCount uint16
	redAlivePlayerCount  uint16
	finalWinnerTeam      uint8
	sendTcpPacketFunc    func(message *tcp.ReceiverMessage)
	game                 *game.Game
}

func NewGameManager(playerNum int, userDb *db.Db, matchScore uint16, sendTcpPacketFunc func(message *tcp.ReceiverMessage)) *GameManager {
	userSpawnPositionArr := make([]int, playerNum)
	for i := 0; i < playerNum; i++ {
		userSpawnPositionArr[i] = i + 1
	}
	return &GameManager{
		userSpawnPositionArr: userSpawnPositionArr,
		userDb:               userDb,
		matchScore:           matchScore,
		blueScore:            0,
		redScore:             0,
		blueAlivePlayerCount: 0,
		redAlivePlayerCount:  0,
		finalWinnerTeam:      0,
		sendTcpPacketFunc:    sendTcpPacketFunc,
		game:                 nil,
	}
}

func (gm *GameManager) InitAlivePlayer() {
	gm.blueAlivePlayerCount = gm.userDb.GetTeamAlivePlayerCount(db.BlueTeam)
	gm.redAlivePlayerCount = gm.userDb.GetTeamAlivePlayerCount(db.RedTeam)
}

func (gm *GameManager) InitGame() {
	gm.InitAlivePlayer()
	util.ShuffleIntArr(gm.userSpawnPositionArr)
	gm.game = game.NewGame(gm.userDb.BlueTeamDb, gm.userDb.RedTeamDb, gm.decreasePlayer)
}

func (gm *GameManager) IncreaseScore(winnerTeam db.Team) {
	gm.matchScore -= 1
	if winnerTeam {
		gm.redScore += 1
		gm.sendTcpPacketFunc(tcp.NewBroadCastMessage('E', tserver.NewRoundEndPacket('R', gm.blueScore, gm.redScore)))
	} else {
		gm.blueScore += 1
		gm.sendTcpPacketFunc(tcp.NewBroadCastMessage('E', tserver.NewRoundEndPacket('B', gm.blueScore, gm.redScore)))
	}

	// matchScore가 0이면 게임 종료 패킷을 던진다.
	// matchScore가 0이 아닐 때는 게임을 초기화 후 라운드 시작 패킷과 라운드 이니셜라이저 패킷을 던진다.
	if gm.matchScore == 0 {
		gm.sendTcpPacketFunc(tcp.NewBroadCastMessage('O', tserver.NewGameOverPacket()))
	} else {

	}
}

func (gm *GameManager) decreasePlayer(deadPlayerTeam db.Team) {
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
