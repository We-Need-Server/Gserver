package game

import (
	"WeNeedGameServer/game/player"
	"strconv"
)

type Game struct {
	Players map[uint32]*player.Player
}

func NewGame() *Game {
	return &Game{Players: make(map[uint32]*player.Player)}
}

func (g *Game) GetPlayer(QPort uint32) player.Player {
	return g.Players[QPort].GetPlayerInfo()
}

func (g *Game) AddPlayer(QPort uint32) *player.Player {
	g.Players[QPort] = player.NewPlayer()
	return g.Players[QPort]
}

func (g Game) GetGameState() string {
	//// 게임 상태를 담을 맵 생성
	gameState := make(map[string]interface{})
	gameState["playerCount"] = len(g.Players)

	//// 플레이어 정보를 담을 맵 생성

	//
	//// 각 플레이어 정보 저장
	for key, p := range g.Players {
		gameState[strconv.Itoa(int(key))] = p
	}
	//
	//gameState["players"] = players
	//
	//// JSON으로 마샬링하여 바이트 배열로 변환
	return "test"
	//return g.GetPlayer(32)
}
