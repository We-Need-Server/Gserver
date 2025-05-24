package tserver

import "encoding/json"

type GameOverPacket struct {
	PKind uint8 `json:"packetKind"`
}

func NewGameOverPacket() *GameOverPacket {
	return &GameOverPacket{PKind: 'O'}
}

func (p *GameOverPacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	return data
}
