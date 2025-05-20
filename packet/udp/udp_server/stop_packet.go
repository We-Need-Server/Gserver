package udp_server

type StopPacket struct {
	pKind uint8
	qPort uint32
}

func NewStopPacket() *StopPacket {
	return &StopPacket{
		pKind: 'S',
		qPort: uint32(0),
	}
}

func (p *StopPacket) GetPacketKind() uint8 {
	return p.pKind
}

func (p *StopPacket) GetQPort() uint32 {
	return p.qPort
}
