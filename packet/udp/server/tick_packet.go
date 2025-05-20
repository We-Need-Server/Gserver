package server

import (
	"WeNeedGameServer/game/player"
	"bytes"
	"encoding/binary"
	"fmt"
)

type TickPacket struct {
	TickNumber         uint32                            `json:"tickNumber"`
	Timestamp          int64                             `json:"timestamp"`
	UserSequenceNumber uint32                            `json:"userSequenceNumber"`
	Flags              uint8                             `json:"flags"`
	UserPositions      map[uint32]*player.PlayerPosition `json:"userPositions"` // 내부 처리용 맵
}

func NewTickPacket(TickNumber uint32, Timestamp int64, UserSequenceNumber uint32, Flags uint8, UserPositions map[uint32]*player.PlayerPosition) *TickPacket {
	return &TickPacket{TickNumber, Timestamp, UserSequenceNumber, Flags, UserPositions}
}

func (p *TickPacket) Serialize() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p.TickNumber)
	binary.Write(buf, binary.LittleEndian, p.Timestamp)
	binary.Write(buf, binary.LittleEndian, p.UserSequenceNumber)
	binary.Write(buf, binary.LittleEndian, p.Flags)
	for qPort, playerPosition := range p.UserPositions {
		fmt.Println("packet")
		fmt.Println(qPort, playerPosition.PositionX, playerPosition.PositionZ, playerPosition.Hp, playerPosition.IsShoot)

		buf.WriteByte('I')
		buf.WriteByte('D')
		binary.Write(buf, binary.LittleEndian, qPort)
		if (*playerPosition).Hp != 0 {
			buf.WriteByte('H')
			buf.WriteByte('T')
			binary.Write(buf, binary.LittleEndian, (*playerPosition).Hp)
		}
		if (*playerPosition).PositionZ != 0 {
			buf.WriteByte('F')
			buf.WriteByte('B')
			binary.Write(buf, binary.LittleEndian, (*playerPosition).PositionZ)
		}
		if (*playerPosition).PositionX != 0 {
			buf.WriteByte('L')
			buf.WriteByte('R')
			binary.Write(buf, binary.LittleEndian, (*playerPosition).PositionX)
		}
		if (*playerPosition).PtAngle != 0 {
			buf.WriteByte('P')
			buf.WriteByte('T')
			binary.Write(buf, binary.LittleEndian, (*playerPosition).PtAngle)
		}
		if (*playerPosition).YawAngle != 0 {
			buf.WriteByte('Y')
			buf.WriteByte('W')
			binary.Write(buf, binary.LittleEndian, (*playerPosition).YawAngle)
		}

		if (*playerPosition).Jp {
			buf.WriteByte('J')
			buf.WriteByte('P')
			fmt.Println("JP-Tick")
		}

		if (*playerPosition).IsShoot {
			buf.WriteByte('S')
			buf.WriteByte('H')
			fmt.Println("SH-Tick")
		}
		if (*playerPosition).IsReload {
			buf.WriteByte('R')
			buf.WriteByte('E')
			fmt.Println("RE-Tick")
		}
	}
	fmt.Println(buf.Bytes())
	return buf.Bytes()
}
