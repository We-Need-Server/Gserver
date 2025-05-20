package client

import "WeNeedGameServer/util"

type TickIPacket struct {
	pKind uint8
	qPort uint32
	SEQ   uint32
}

func (p *TickIPacket) GetPacketKind() uint8 {
	return p.pKind
}

func newTickIPacket(pKind uint8, qPort uint32, seq uint32) *TickIPacket {

	return &TickIPacket{pKind, qPort, seq}
}

func ParseTickIPacket(np []byte, endPoint int) *TickIPacket {
	pKind := np[0]
	qPort := util.ConvertBinaryToUint32(np[1:5])
	seq := util.ConvertBinaryToUint32(np[5:])
	return newTickIPacket(pKind, qPort, seq)
}

func (p *TickIPacket) GetQPort() uint32 {
	return p.qPort
}

func (p *TickIPacket) GetSEQ() uint32 {
	return p.SEQ
}
