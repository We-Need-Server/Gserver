package server

import (
	"WeNeedGameServer/game/player"
	"bytes"
	"encoding/binary"
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
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p.TickNumber)
	binary.Write(buf, binary.LittleEndian, p.Timestamp)
	binary.Write(buf, binary.LittleEndian, p.UserSequenceNumber)
	binary.Write(buf, binary.LittleEndian, p.Flags)
	for qPort, playerPosition := range p.UserPositions {
		binary.Write(buf, binary.LittleEndian, "ID")
		binary.Write(buf, binary.LittleEndian, qPort)
		binary.Write(buf, binary.LittleEndian, playerPosition)
	}
	return buf.Bytes()
}
