package game

import (
	"WeNeedGameServer/common"
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game/player"
	"fmt"
)

type Game struct {
	blueTeam             map[uint32]*db.User
	redTeam              map[uint32]*db.User
	userSpawnPositionArr []int
	decreasePlayerFunc   func(team db.Team)
	players              map[uint32]*player.Player
}

// 이제 게임 부분만 구축하면 끝!
func NewGame(blueTeam map[uint32]*db.User, redTeam map[uint32]*db.User, userSpawnPositionArr []int, decreasePlayerFunc func(team db.Team)) *Game {
	return &Game{
		blueTeam:             blueTeam,
		redTeam:              redTeam,
		userSpawnPositionArr: userSpawnPositionArr,
		decreasePlayerFunc:   decreasePlayerFunc,
		players:              make(map[uint32]*player.Player),
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

func (g *Game) GetGameState() map[uint32]*player.PlayerPosition {
	gameState := make(map[uint32]*player.PlayerPosition)
	for userId, p := range g.players {
		gameState[userId] = p.GetPlayerState()
		if !gameState[userId].IsAlive {
			fmt.Println("유저가 죽었습니다", userId, gameState[userId].Hp, gameState[userId].IsAlive)
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
	g.players[userId] = player.NewPlayer(respawnPosition, team)
}

func (g *Game) DeletePlayer(userId uint32) {
	delete(g.players, userId)
}
func (g *Game) ReflectPlayers(playerPositionMap map[uint32]*player.PlayerPosition) {
	for key, val := range playerPositionMap {
		if _, exists := g.players[key]; !exists {
			g.addPlayer(key, val.RespawnPoint, val.Team)
		}
		g.players[key].ReflectPlayerPosition(val)
	}
}
