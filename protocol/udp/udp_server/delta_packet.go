package udp_server

import (
	"WeNeedGameServer/game/player"
)

type DeltaPacket struct {
	pKind             uint8
	qPort             uint32
	PlayerPosition    *player.PlayerPosition
	HitInformationMap *map[uint32]int16
}

func NewDeltaPacket(qPort uint32, playerPosition *player.PlayerPosition, hitInformationMap *map[uint32]int16) *DeltaPacket {
	return &DeltaPacket{
		pKind:             'D',
		qPort:             qPort,
		PlayerPosition:    playerPosition,
		HitInformationMap: hitInformationMap,
	}
}

func (p *DeltaPacket) GetPacketKind() uint8 {
	return p.pKind
}

func (p *DeltaPacket) GetQPort() uint32 {
	return p.qPort
}
