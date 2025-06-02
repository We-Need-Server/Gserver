package game

import (
	"WeNeedGameServer/common"
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game/entity"
)

type Game struct {
	blueTeam             map[uint32]*db.User
	redTeam              map[uint32]*db.User
	userSpawnPositionArr []int
	decreasePlayerFunc   func(team db.Team)
	players              map[uint32]*entity.Player
}

func NewGame(blueTeam map[uint32]*db.User, redTeam map[uint32]*db.User, userSpawnPositionArr []int, decreasePlayerFunc func(team db.Team)) *Game {
	return &Game{
		blueTeam:             blueTeam,
		redTeam:              redTeam,
		userSpawnPositionArr: userSpawnPositionArr,
		decreasePlayerFunc:   decreasePlayerFunc,
		players:              make(map[uint32]*entity.Player),
	}

}

func (g *Game) ReadyGame() *Game {
	playerPositionIndex := 0
	// 블루팀 스폰
	for key, _ := range g.blueTeam {
		g.addPlayer(key, -1*g.userSpawnPositionArr[playerPositionIndex], db.BlueTeam)
		playerPositionIndex += 1
	}
	playerPositionIndex = 0
	// 레드팀 스폰
	for key, _ := range g.redTeam {
		g.addPlayer(key, g.userSpawnPositionArr[playerPositionIndex], db.RedTeam)
		playerPositionIndex += 1
	}
	return g
}

func (g *Game) GetGameState() map[uint32]*entity.PlayerState {
	gameState := make(map[uint32]*entity.PlayerState)
	for userId, p := range g.players {
		gameState[userId] = p.GetPlayerState()
		if !gameState[userId].IsAlive {
			g.decreasePlayerFunc(gameState[userId].Team)
		}
	}
	return gameState
}

func (g *Game) GetPlayerSpawnStatusList() []*common.UserSpawnStatus {
	var userSpawnStatusArr []*common.UserSpawnStatus
	for key, val := range g.players {
		userSpawnStatusArr = append(userSpawnStatusArr, common.NewUserSpawnStatus(key, int16(val.RespawnPoint)))
	}
	return userSpawnStatusArr
}

func (g *Game) addPlayer(userId uint32, respawnPosition int, team db.Team) {
	g.players[userId] = entity.NewPlayer(respawnPosition, team)
}

func (g *Game) DeletePlayer(userId uint32) {
	delete(g.players, userId)
}

func (g *Game) ReflectPlayers(playerPositionMap map[uint32]*entity.PlayerState) {
	for key, val := range playerPositionMap {
		g.players[key].ReflectPlayer(val)
	}
}
