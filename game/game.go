package game

import (
	"WeNeedGameServer/game/player"
)

type Game struct {
	redAlivePlayerCount  uint16
	blueAlivePlayerCount uint16
	increaseScore        func(uint8)
	players              map[uint32]*player.Player
}

func NewGame(blueAlivePlayerCount uint16, redAlivePlayerCount uint16, increaseScore func(uint8)) *Game {
	return &Game{
		blueAlivePlayerCount: blueAlivePlayerCount,
		redAlivePlayerCount:  redAlivePlayerCount,
		increaseScore:        increaseScore,
		players:              make(map[uint32]*player.Player),
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

func (g *Game) decreasePlayer(deadPlayerTeam uint8) {
	if deadPlayerTeam == 'B' {
		g.blueAlivePlayerCount -= 1
		if g.blueAlivePlayerCount == 0 {
			g.increaseScore('A')
		}
	} else {
		g.redAlivePlayerCount -= 1
	}

}
