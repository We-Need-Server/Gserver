package tickrpacket

import "WeNeedGameServer/util"

type TickRPacket struct {
	QPort       uint32
	RTickNumber uint32
}

func (*TickRPacket) GetPacketKind() uint32 {
	return 50
}

func newTickRPacket(QPort uint32, RTickNumber uint32) *TickRPacket {
	return &TickRPacket{QPort, RTickNumber}
}

func ParseTickRPacket(np []byte, endPoint int) *TickRPacket {
	QPort := util.ConvertBinaryToUint32(np[4:8])
	RTickNumber := util.ConvertBinaryToUint32(np[8:12])
	return newTickRPacket(QPort, RTickNumber)
}

func (p *TickRPacket) GetQPort() uint32 {
	return p.QPort
}
