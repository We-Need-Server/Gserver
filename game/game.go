package game

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game/entity"
	"WeNeedGameServer/game_type"
)

type Game struct {
	blueTeam             map[uint32]*db.User
	redTeam              map[uint32]*db.User
	userSpawnPositionArr []int
	decreasePlayerFunc   func(team game_type.Team)
	players              map[uint32]*entity.Player
}

func NewGame(blueTeam map[uint32]*db.User, redTeam map[uint32]*db.User, userSpawnPositionArr []int, decreasePlayerFunc func(team game_type.Team)) *Game {
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
		g.addPlayer(key, -1*g.userSpawnPositionArr[playerPositionIndex], game_type.BlueTeam)
		playerPositionIndex += 1
	}
	playerPositionIndex = 0
	// 레드팀 스폰
	for key, _ := range g.redTeam {
		g.addPlayer(key, g.userSpawnPositionArr[playerPositionIndex], game_type.RedTeam)
		playerPositionIndex += 1
	}
	return g
}

func (g *Game) GetGameState() map[uint32]*game_type.PlayerState {
	gameState := make(map[uint32]*game_type.PlayerState)
	for userId, p := range g.players {
		gameState[userId] = p.GetPlayerState()
		if !gameState[userId].IsAlive {
			g.decreasePlayerFunc(gameState[userId].Team)
		}
	}
	return gameState
}

func (g *Game) GetPlayerSpawnStatusList() []*game_type.UserSpawnStatus {
	var userSpawnStatusArr []*game_type.UserSpawnStatus
	for key, val := range g.players {
		userSpawnStatusArr = append(userSpawnStatusArr, game_type.NewUserSpawnStatus(key, int16(val.RespawnPoint)))
	}
	return userSpawnStatusArr
}

func (g *Game) addPlayer(userId uint32, respawnPosition int, team game_type.Team) {
	g.players[userId] = entity.NewPlayer(respawnPosition, team)
}

func (g *Game) DeletePlayer(userId uint32) {
	delete(g.players, userId)
}

func (g *Game) ReflectPlayers(playerPositionMap map[uint32]*game_type.PlayerState) {
	for key, val := range playerPositionMap {
		g.players[key].ReflectPlayer(val)
	}
}
