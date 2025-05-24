package tserver

import "encoding/json"

type GameOverPacket struct {
	PKind uint8 `json:"-"`
}

func NewGameOverPacket() *GameOverPacket {
	return &GameOverPacket{PKind: 'O'}
}

func (p *GameOverPacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	result := make([]byte, 1+len(data))
	result[0] = p.PKind
	copy(result[1:], data)

	return result
}
