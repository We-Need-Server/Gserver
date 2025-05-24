package tserver

import "encoding/json"

type ConnectionResponsePacket struct {
	PKind      uint8  `json:"packetKind"`
	QPort      uint32 `json:"qPort"`
	UdpAddr    string `json:"udpAddr"`
	MatchScore uint16 `json:"matchScore"`
}

func NewConnectionResponsePacket(qPort uint32, udpAddr string, matchScore uint16) *ConnectionResponsePacket {
	return &ConnectionResponsePacket{
		PKind:      'I',
		QPort:      qPort,
		UdpAddr:    udpAddr,
		MatchScore: matchScore,
	}
}

func (p *ConnectionResponsePacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	return data
}
