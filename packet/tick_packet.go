package packet

import (
	"WeNeedGameServer/game/player"
	"time"
)

type TickPacket struct {
	TickNumber         uint32                           `json:"tickNumber"`
	Timestamp          time.Time                        `json:"timestamp"`
	UserSequenceNumber uint32                           `json:"userSequenceNumber"`
	Flag               rune                             `json:"flag"`
	UserPositions      map[string]player.PlayerPosition `json:"userPositions"` // 내부 처리용 맵
}

func NewTickPacket(TickNumber uint32, Timestamp time.Time, UserSequenceNumber uint32, Flag rune, UserPositions map[string]player.PlayerPosition) *TickPacket {
	return &TickPacket{TickNumber, Timestamp, UserSequenceNumber, Flag, UserPositions}
}
