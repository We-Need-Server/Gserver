package packet

import "WeNeedGameServer/util"

type TickIPacket struct {
	QPort uint32
}

func (*TickIPacket) GetPacketKind() uint32 {
	return 41
}

func newTickIPacket(QPort uint32) *TickIPacket {
	return &TickIPacket{QPort}
}

func ParseTickIPacket(np []byte, endPoint int) *TickIPacket {
	QPort := util.ConvertBinaryToUint32(np[4:8])
	return newTickIPacket(QPort)
}

func (p *TickIPacket) GetQPort() uint32 {
	return p.QPort
}
