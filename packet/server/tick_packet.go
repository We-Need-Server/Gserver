package server

import (
	"WeNeedGameServer/game/player"
	"time"
)

type TickPacket struct {
	TickNumber         int                              `json:"tickNumber"`
	Timestamp          time.Time                        `json:"timestamp"`
	UserSequenceNumber uint32                           `json:"userSequenceNumber"`
	Flags              uint8                            `json:"flags"`
	UserPositions      map[string]player.PlayerPosition `json:"userPositions"` // 내부 처리용 맵
}

func NewTickPacket(TickNumber int, Timestamp time.Time, UserSequenceNumber uint32, Flags uint8, UserPositions map[string]player.PlayerPosition) *TickPacket {
	return &TickPacket{TickNumber, Timestamp, UserSequenceNumber, Flags, UserPositions}
}
