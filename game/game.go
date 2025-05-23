package game

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game/player"
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

func (g *Game) ReadyGame() {
	playerPositionIndex := 0
	// 블루팀 스폰
	for key, _ := range g.blueTeam {
		g.addPlayer(key, g.userSpawnPositionArr[playerPositionIndex])
		playerPositionIndex += 1
	}
	playerPositionIndex = 0
	// 레드팀 스폰
	for key, _ := range g.redTeam {
		g.addPlayer(key, g.userSpawnPositionArr[playerPositionIndex])
		playerPositionIndex += 1
	}
}

func (g *Game) GetGameState() map[uint32]*player.PlayerPosition {
	gameState := make(map[uint32]*player.PlayerPosition)
	for qPort, p := range g.players {
		gameState[qPort] = p.GetPlayerState()
	}
	return gameState
}

func (g *Game) addPlayer(userId uint32, respawnPosition int) {
	g.players[userId] = player.NewPlayer(respawnPosition)
}

func (g *Game) ReflectPlayers(playerPositionMap map[uint32]*player.PlayerPosition) {
	for key, val := range playerPositionMap {
		//if _, exists := g.players[key]; !exists {
		//	g.addPlayer(key)
		//}
		g.players[key].ReflectPlayerPosition(val)
	}
}
