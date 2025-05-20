package tcp_server

type GameOverPacket struct {
	PKind uint8 `json:"packetKind"`
}

func NewGameOverPacket() *GameOverPacket {
	return &GameOverPacket{PKind: 'O'}
}
