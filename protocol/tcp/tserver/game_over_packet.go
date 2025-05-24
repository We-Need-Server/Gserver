package tserver

type GameOverPacket struct {
	PKind uint8 `json:"packetKind"`
}

func NewGameOverPacket() *GameOverPacket {
	return &GameOverPacket{PKind: 'O'}
}

func (p *GameOverPacket) Serialize() []byte {
	return []byte{}
}
