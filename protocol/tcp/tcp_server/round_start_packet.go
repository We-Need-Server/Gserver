package tcp_server

import (
	"WeNeedGameServer/game_manager/internal/types"
)

type RoundStartPacket struct {
	PKind              uint8                   `json:"packetKind"`
	UserSpawnStatusArr []_type.UserSpawnStatus `json:"userSpawnStatusArr"`
}

func NewRoundStartPacket(userSpawnStatus []_type.UserSpawnStatus) *RoundStartPacket {
	return &RoundStartPacket{
		PKind:              'S',
		UserSpawnStatusArr: userSpawnStatus,
	}
}
