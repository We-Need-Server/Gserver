package game

import (
	"WeNeedGameServer/game/player"
)

type Game struct {
	Players map[uint32]*player.Player
}

func NewGame() *Game {
	return &Game{Players: make(map[uint32]*player.Player)}
}

func (g *Game) AddPlayer(QPort uint32) *player.Player {
	g.Players[QPort] = player.NewPlayer()
	return g.Players[QPort]
}

func (g *Game) GetGameDeltaState() map[uint32]player.PlayerPosition {
	gameDeltaState := make(map[uint32]player.PlayerPosition)
	for qPort, p := range g.Players {
		gameDeltaState[qPort] = player.NewPlayerPosition(p.XDelta, p.ZDelta, p.YawDelta, p.PTDelta, p.JP)
		p.ReflectDeltaValues()
	}
	return gameDeltaState
}

func (g *Game) GetGameState() map[uint32]player.PlayerPosition {
	gameState := make(map[uint32]player.PlayerPosition)
	for qPort, p := range g.Players {
		gameState[qPort] = player.NewPlayerPosition(p.PositionX, p.PositionZ, p.YawAngle, p.PTAngle, p.JP)
	}

	return gameState
}
