package tserver

type RoundStartPacket struct {
	PKind uint8 `json:"packetKind"`
}

func NewRoundStartPacket() *RoundStartPacket {
	return &RoundStartPacket{
		PKind: 'S',
	}
}

func (p *RoundStartPacket) Serialize() []byte {
	return []byte{}
}
