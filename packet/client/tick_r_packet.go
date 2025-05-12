package client

import "WeNeedGameServer/util"

type TickRPacket struct {
	pKind       uint32
	qPort       uint32
	SEQ         uint32
	RTickNumber uint32
}

func (p *TickRPacket) GetPacketKind() uint32 {
	return p.pKind
}

func newTickRPacket(pKind uint32, qPort uint32, seq uint32, RTickNumber uint32) *TickRPacket {
	return &TickRPacket{pKind, qPort, seq, RTickNumber}
}

func ParseTickRPacket(np []byte, endPoint int) *TickRPacket {
	pKind := util.ConvertBinaryToUint32(np[0:4])
	qPort := util.ConvertBinaryToUint32(np[4:8])
	seq := util.ConvertBinaryToUint32(np[8:12])
	RTickNumber := util.ConvertBinaryToUint32(np[12:endPoint])
	return newTickRPacket(pKind, qPort, seq, RTickNumber)
}

func (p *TickRPacket) GetQPort() uint32 {
	return p.qPort
}

func (p *TickRPacket) GetSEQ() uint32 {
	return p.SEQ
}
