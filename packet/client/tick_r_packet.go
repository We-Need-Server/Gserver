package client

import "WeNeedGameServer/util"

type TickRPacket struct {
	pKind       uint32
	qPort       uint32
	RTickNumber uint32
}

func (p *TickRPacket) GetPacketKind() uint32 {
	return p.pKind
}

func newTickRPacket(pKind uint32, qPort uint32, RTickNumber uint32) *TickRPacket {
	return &TickRPacket{pKind, qPort, RTickNumber}
}

func ParseTickRPacket(np []byte, endPoint int) *TickRPacket {
	pKind := util.ConvertBinaryToUint32(np[0:4])
	qPort := util.ConvertBinaryToUint32(np[4:8])
	RTickNumber := util.ConvertBinaryToUint32(np[8:endPoint])
	return newTickRPacket(pKind, qPort, RTickNumber)
}

func (p *TickRPacket) GetQPort() uint32 {
	return p.qPort
}
