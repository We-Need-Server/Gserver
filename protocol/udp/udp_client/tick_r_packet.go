package udp_client

import (
	"WeNeedGameServer/util"
)

type TickRPacket struct {
	pKind       uint8
	qPort       uint32
	SEQ         uint32
	RTickNumber uint32
}

func (p *TickRPacket) GetPacketKind() uint8 {
	return p.pKind
}

func newTickRPacket(pKind uint8, qPort uint32, seq uint32, RTickNumber uint32) *TickRPacket {
	return &TickRPacket{pKind, qPort, seq, RTickNumber}
}

func ParseTickRPacket(np []byte, endPoint int) *TickRPacket {
	pKind := np[0]
	qPort := util.ConvertBinaryToUint32(np[1:5])
	seq := util.ConvertBinaryToUint32(np[5:9])
	RTickNumber := util.ConvertBinaryToUint32(np[9:endPoint])
	return newTickRPacket(pKind, qPort, seq, RTickNumber)
}

func (p *TickRPacket) GetQPort() uint32 {
	return p.qPort
}

func (p *TickRPacket) GetSEQ() uint32 {
	return p.SEQ
}
