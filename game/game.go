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

func (g Game) GetGameState() map[string]player.PlayerPosition {
	//// 게임 상태를 담을 맵 생성
	gameState := make(map[string]player.PlayerPosition)

	//// 플레이어 정보를 담을 맵 생성

	//
	//// 각 플레이어 정보 저장
	for key, p := range g.Players {
		//fmt.Println(p)
		gameState[strconv.Itoa(int(key))] = player.NewPlayerPosition(p.PositionX, p.PositionY, p.PositionZ, p.YawAngle, p.PTAngle)
	}

	//fmt.Println(gameState)
	//
	//gameState["players"] = players
	//
	return gameState
}
