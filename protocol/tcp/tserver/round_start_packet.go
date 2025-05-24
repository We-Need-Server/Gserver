package tserver

import "encoding/json"

type RoundStartPacket struct {
	PKind uint8 `json:"packetKind"`
}

func NewRoundStartPacket() *RoundStartPacket {
	return &RoundStartPacket{
		PKind: 'S',
	}
}

func (p *RoundStartPacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	return data
}
