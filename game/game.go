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
		for key, val := range p.ShootHitInformation {
			g.Players[key].HPDelta += val
		}
		p.ReflectHitInformation()
		gameDeltaState[qPort] = player.NewPlayerPosition(&p.HPDelta, p.XDelta, p.ZDelta, p.YawDelta, p.PTDelta, p.JP, p.IsShoot)
		p.ReflectDeltaValues()
	}
	return gameDeltaState
}

func (g *Game) GetGameState() map[uint32]player.PlayerPosition {
	gameState := make(map[uint32]player.PlayerPosition)
	for qPort, p := range g.Players {
		p.ReflectDamageHP()
		gameState[qPort] = player.NewPlayerPosition(&p.HP, p.PositionX, p.PositionZ, p.YawAngle, p.PTAngle, p.JP, p.IsShoot)
	}
	return gameState
}
