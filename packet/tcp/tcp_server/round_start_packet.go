package tcp_server

import "WeNeedGameServer/internal_type"

type RoundStartPacket struct {
	PKind              uint8                           `json:"packetKind"`
	UserSpawnStatusArr []internal_type.UserSpawnStatus `json:"userSpawnStatusArr"`
}

func NewRoundStartPacket(userSpawnStatus []internal_type.UserSpawnStatus) *RoundStartPacket {
	return &RoundStartPacket{
		PKind:              'S',
		UserSpawnStatusArr: userSpawnStatus,
	}
}
