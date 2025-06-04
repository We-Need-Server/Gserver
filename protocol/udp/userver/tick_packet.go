package userver

import (
	"WeNeedGameServer/game_type"
	"bytes"
	"encoding/binary"
	"fmt"
)

type TickPacket struct {
	TickNumber         uint32
	Timestamp          int64
	UserSequenceNumber uint32
	Flags              uint8
	PlayerStateMap     map[uint32]*game_type.PlayerState
}

func NewTickPacket(TickNumber uint32, Timestamp int64, UserSequenceNumber uint32, Flags uint8, UserPositions map[uint32]*game_type.PlayerState) *TickPacket {
	return &TickPacket{TickNumber, Timestamp, UserSequenceNumber, Flags, UserPositions}
}

func (p *TickPacket) Serialize() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p.TickNumber)
	binary.Write(buf, binary.LittleEndian, p.Timestamp)
	binary.Write(buf, binary.LittleEndian, p.UserSequenceNumber)
	binary.Write(buf, binary.LittleEndian, p.Flags)
	for qPort, playerState := range p.PlayerStateMap {
		buf.WriteByte('I')
		buf.WriteByte('D')
		binary.Write(buf, binary.LittleEndian, qPort)
		if playerState.Damage != 0 {
			buf.WriteByte('H')
			buf.WriteByte('T')
			binary.Write(buf, binary.LittleEndian, playerState.Damage)
		}
		if playerState.PositionZ != 0 {
			buf.WriteByte('F')
			buf.WriteByte('B')
			binary.Write(buf, binary.LittleEndian, playerState.PositionZ)
		}
		if playerState.PositionX != 0 {
			buf.WriteByte('L')
			buf.WriteByte('R')
			binary.Write(buf, binary.LittleEndian, playerState.PositionX)
		}
		if playerState.PtAngle != 0 {
			buf.WriteByte('P')
			buf.WriteByte('T')
			binary.Write(buf, binary.LittleEndian, playerState.PtAngle)
		}
		if playerState.YawAngle != 0 {
			buf.WriteByte('Y')
			buf.WriteByte('W')
			binary.Write(buf, binary.LittleEndian, playerState.YawAngle)
		}

		if playerState.Jp {
			buf.WriteByte('J')
			buf.WriteByte('P')
			fmt.Println("JP-Tick")
		}

		if playerState.IsShoot {
			buf.WriteByte('S')
			buf.WriteByte('H')
			fmt.Println("SH-Tick")
		}

		if playerState.IsReload {
			buf.WriteByte('R')
			buf.WriteByte('E')
			fmt.Println("RE-Tick")
		}

		buf.WriteByte('R')
		buf.WriteByte('P')
		binary.Write(buf, binary.LittleEndian, playerState.RespawnPoint)

	}
	return buf.Bytes()
}
