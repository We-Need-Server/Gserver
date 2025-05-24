package tserver

import "encoding/json"

type ConnectionResponsePacket struct {
	PKind   uint8  `json:"packetKind"`
	QPort   uint32 `json:"qPort"`
	UdpAddr string `json:"udpAddr"`
}

func NewConnectionResponsePacket(qPort uint32, udpAddr string) *ConnectionResponsePacket {
	return &ConnectionResponsePacket{
		PKind:   'I',
		QPort:   qPort,
		UdpAddr: udpAddr,
	}
}

func (p *ConnectionResponsePacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	return data
}
