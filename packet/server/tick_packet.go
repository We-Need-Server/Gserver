package server

import (
	"WeNeedGameServer/game/player"
)

type TickPacket struct {
	TickNumber         int                              `json:"tickNumber"`
	Timestamp          int64                            `json:"timestamp"`
	UserSequenceNumber uint32                           `json:"userSequenceNumber"`
	Flags              uint8                            `json:"flags"`
	UserPositions      map[uint32]player.PlayerPosition `json:"userPositions"` // 내부 처리용 맵
}

func NewTickPacket(TickNumber int, Timestamp int64, UserSequenceNumber uint32, Flags uint8, UserPositions map[uint32]player.PlayerPosition) *TickPacket {
	return &TickPacket{TickNumber, Timestamp, UserSequenceNumber, Flags, UserPositions}
}

func (p *TickPacket) Serialize() []byte {
	return []byte{}
}
