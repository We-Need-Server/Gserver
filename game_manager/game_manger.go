package game_manager

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game"
	"WeNeedGameServer/game_manager/internal"
	"WeNeedGameServer/protocol/tcp"
	"WeNeedGameServer/protocol/tcp/tserver"
	"WeNeedGameServer/util"
	"time"
)

// 스냅샷 방식으로 살아있는 플레이어 수를 정한다.
type GameManager struct {
	userSpawnPositionArr []int
	userDb               *db.Db
	matchScore           uint16 // 게임이 총 몇 판 몇 선제일때의 몇 판
	blueScore            uint16 // 라운드 승리 횟수
	redScore             uint16
	finalWinnerTeam      uint8
	sendTcpPacketFunc    func(message *tcp.ReceiverMessage)
	gameNetwork          *internal.GameNetwork
	gameTick             *internal.GameTick
	game                 *game.Game
}

func NewGameManager(playerNum int, userDb *db.Db, matchScore uint16, sendTcpPacketFunc func(message *tcp.ReceiverMessage), listenUdpAddr string) *GameManager {
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
		sendTcpPacketFunc:    sendTcpPacketFunc,
		gameNetwork:          internal.NewGameNetwork(listenUdpAddr, userDb.FindUserByQPort),
		gameTick:             nil,
		game:                 nil,
	}
}

//func (gm *GameManager) InitAlivePlayer() {
//	gm.blueAlivePlayerCount = uint16(gm.userDb.GetTeamAliveCount(db.BlueTeam))
//	gm.redAlivePlayerCount = gm.userDb.GetTeamAlivePlayerCount(db.RedTeam)
//}

func (gm *GameManager) StartGameManager() {
	gm.gameNetwork.ReadyUdp()
	go gm.gameNetwork.UdpReceiver.StartUdp()
	gm.InitGame()
	gm.gameTick = internal.NewGameTick(60, gm.game, gm.gameNetwork.UdpSender)
	go gm.gameTick.StartGameLoop()
}

func (gm *GameManager) InitGame() {
	util.ShuffleIntArr(gm.userSpawnPositionArr)
	gameInstance := game.NewGame(gm.userDb.BlueTeamDb, gm.userDb.RedTeamDb, gm.userSpawnPositionArr, gm.decreasePlayer)
	gm.game = gameInstance.ReadyGame()
}

func (gm *GameManager) IncreaseTeamScore(winnerTeam db.Team) {
	if winnerTeam == db.RedTeam {
		gm.redScore += 1
	} else {
		gm.blueScore += 1
	}
}

func (gm *GameManager) ReadyNextRound(winnerTeam db.Team) {
	gm.matchScore -= 1
	gm.IncreaseTeamScore(winnerTeam)
	gm.sendTcpPacketFunc(tcp.NewBroadCastMessage('E', tserver.NewRoundEndPacket(winnerTeam, gm.blueScore, gm.redScore)))
	time.Sleep(5 * time.Second)
	if gm.matchScore == 0 {
		gm.sendTcpPacketFunc(tcp.NewBroadCastMessage('O', tserver.NewGameOverPacket()))
		gm.game = nil
	} else {
		gm.InitGame()
	}
}

func (gm *GameManager) decreasePlayer(deadPlayerTeam db.Team) {
	gm.userDb.DecreaseTeamAliveCount(deadPlayerTeam)
	if gm.userDb.GetTeamAliveCount(deadPlayerTeam) == 0 {
		gm.ReadyNextRound(!deadPlayerTeam)
	}
}
