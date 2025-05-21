package tserver

import (
	"WeNeedGameServer/game/player"
)

type GameInitPacket struct {
	PKind             uint32
	PlayerPositionMap map[uint32]*player.PlayerPosition
}

func NewGameInitPacket(playerPositionMap map[uint32]*player.PlayerPosition) *GameInitPacket {
	return &GameInitPacket{
		PKind:             'G',
		PlayerPositionMap: playerPositionMap,
	}
}
