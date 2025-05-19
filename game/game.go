package game

import "WeNeedGameServer/game/player"

type Game struct {
	players map[uint32]*player.Player
}

func NewGame() *Game {
	return &Game{players: make(map[uint32]*player.Player)}
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
