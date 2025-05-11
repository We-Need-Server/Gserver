package client

import "WeNeedGameServer/util"

type TickIPacket struct {
	pKind uint32
	qPort uint32
	SEQ   uint32
}

func (p *TickIPacket) GetPacketKind() uint32 {
	return p.pKind
}

func newTickIPacket(pKind uint32, qPort uint32, seq uint32) *TickIPacket {

	return &TickIPacket{pKind, qPort, seq}
}

func ParseTickIPacket(np []byte, endPoint int) *TickIPacket {
	pKind := util.ConvertBinaryToUint32(np[0:4])
	qPort := util.ConvertBinaryToUint32(np[4:8])
	seq := util.ConvertBinaryToUint32(np[8:])
	return newTickIPacket(pKind, qPort, seq)
}

func (p *TickIPacket) GetQPort() uint32 {
	return p.qPort
}

func (p *TickIPacket) GetSEQ() uint32 {
	return p.SEQ
}
