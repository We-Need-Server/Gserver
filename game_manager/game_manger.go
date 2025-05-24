package game_manager

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game"
	"WeNeedGameServer/game_manager/internal"
	"WeNeedGameServer/protocol/tcp"
	"WeNeedGameServer/protocol/tcp/tserver"
	"WeNeedGameServer/util"
	"fmt"
	"time"
)

type GameStatus uint16

const (
	GameReady GameStatus = iota
	RoundStart
	RoundEnd
	GameEnd
)

// 스냅샷 방식으로 살아있는 플레이어 수를 정한다.
type GameManager struct {
	GameStatus           GameStatus
	userSpawnPositionArr []int
	userDb               *db.Db
	matchScore           uint16 // 게임이 총 몇 판 몇 선제일때의 몇 판
	blueScore            uint16 // 라운드 승리 횟수
	redScore             uint16
	finalWinnerTeam      uint8
	sendTcpPacketFunc    func(message *tcp.Message)
	gameNetwork          *internal.GameNetwork
	gameTick             *internal.GameTick
	game                 *game.Game
}

func NewGameManager(playerNum int, userDb *db.Db, matchScore uint16, sendTcpPacketFunc func(message *tcp.Message), listenUdpAddr string) *GameManager {
	userSpawnPositionArr := make([]int, playerNum)
	for i := 0; i < playerNum; i++ {
		userSpawnPositionArr[i] = i + 1
	}
	return &GameManager{
		GameStatus:           GameReady,
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
	fmt.Println("게임 시작합니다.")
	gm.gameNetwork.ReadyUdp()
	go gm.gameNetwork.UdpReceiver.StartUdp()
	gm.initGame()
	gm.gameTick = internal.NewGameTick(60, gm.game, gm.gameNetwork.UdpSender, gm.userDb.CheckLogin)
	gm.sendTcpPacketFunc(tcp.NewBroadCastMessage(tserver.NewGameInitPacket(gm.gameTick.TickTime, gm.blueScore, gm.redScore, gm.game.GetPlayerSpawnStatusList())))
	gm.GameStatus = RoundStart
	go gm.gameTick.StartGameLoop()
}

func (gm *GameManager) initGame() {
	util.ShuffleIntArr(gm.userSpawnPositionArr)
	gameInstance := game.NewGame(gm.userDb.BlueTeamDb, gm.userDb.RedTeamDb, gm.userSpawnPositionArr, gm.decreasePlayer)
	gm.game = gameInstance.ReadyGame()
	if gm.GameStatus != GameReady {
		gm.sendTcpPacketFunc(tcp.NewBroadCastMessage(tserver.NewGameInitPacket(gm.gameTick.TickTime, gm.blueScore, gm.redScore, gm.game.GetPlayerSpawnStatusList())))
		gm.GameStatus = RoundStart
	}
}

func (gm *GameManager) SendGameInitPacket(userId uint32) {
	gm.sendTcpPacketFunc(tcp.NewUniCastMessage(userId, tserver.NewGameInitPacket(gm.gameTick.TickTime, gm.blueScore, gm.redScore, gm.game.GetPlayerSpawnStatusList())))
}

func (gm *GameManager) increaseTeamScore(winnerTeam db.Team) {
	if winnerTeam == db.RedTeam {
		gm.redScore += 1
	} else {
		gm.blueScore += 1
	}
}

func (gm *GameManager) readyNextRound(winnerTeam db.Team) {
	gm.matchScore -= 1
	gm.increaseTeamScore(winnerTeam)
	gm.sendTcpPacketFunc(tcp.NewBroadCastMessage(tserver.NewRoundEndPacket(winnerTeam, gm.blueScore, gm.redScore)))
	time.Sleep(5 * time.Second)
	if gm.matchScore == 0 {
		gm.GameStatus = GameEnd
		gm.sendTcpPacketFunc(tcp.NewBroadCastMessage(tserver.NewGameOverPacket()))
		// 여기서 gameTick을 종료시키는 방법을 찾아야 한다. 보니까 컨텍스트를 쓸 수 있을 거 같은데 어떻게 하는게 좋을까
		gm.game = nil
	} else {
		gm.sendTcpPacketFunc(tcp.NewBroadCastMessage(tserver.NewRoundStartPacket()))
		gm.initGame()
	}
}

func (gm *GameManager) decreasePlayer(deadPlayerTeam db.Team) {
	gm.userDb.DecreaseTeamAliveCount(deadPlayerTeam)
	if gm.userDb.GetTeamAliveCount(deadPlayerTeam) == 0 {
		gm.GameStatus = RoundEnd
		gm.readyNextRound(!deadPlayerTeam)
	}
}
