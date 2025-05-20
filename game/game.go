package game

import (
	"WeNeedGameServer/game/player"
	"WeNeedGameServer/game/round"
)

type Game struct {
	round *round.Round

	players map[uint32]*player.Player
}

func NewGame(round *round.Round) *Game {
	return &Game{
		round: round,

		players: make(map[uint32]*player.Player),
	}
}

func (g *Game) GetGameState() map[uint32]*player.PlayerPosition {
	gameState := make(map[uint32]*player.PlayerPosition)
	for qPort, p := range g.players {
		gameState[qPort] = p.GetPlayerState()
	}
	return gameState
}

func (g *Game) addPlayer(qPort uint32) {
	g.players[qPort] = player.NewPlayer()
}

func (g *Game) ReflectPlayers(playerPositionMap *map[uint32]*player.PlayerPosition) {
	for key, val := range *playerPositionMap {
		if _, exists := g.players[key]; !exists {
			g.addPlayer(key)
		}
		g.players[key].ReflectPlayerPosition(val)
	}
}
