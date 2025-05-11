package client

import "WeNeedGameServer/util"

type TickIPacket struct {
	pKind uint32
	qPort uint32
}

func (p *TickIPacket) GetPacketKind() uint32 {
	return p.pKind
}

func newTickIPacket(pKind uint32, qPort uint32) *TickIPacket {
	return &TickIPacket{pKind, qPort}
}

func ParseTickIPacket(np []byte, endPoint int) *TickIPacket {
	pKind := util.ConvertBinaryToUint32(np[0:4])
	qPort := util.ConvertBinaryToUint32(np[4:8])
	return newTickIPacket(pKind, qPort)
}

func (p *TickIPacket) GetQPort() uint32 {
	return p.qPort
}
