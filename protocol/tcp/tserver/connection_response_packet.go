package tserver

type ConnectionResponsePacket struct {
	PKind   uint8  `json:"packetKind"`
	QPort   uint32 `json:"qPort"`
	UdpAddr string `json:"udpAddr"`
	//UserSpawnStatusArr []internal_type.UserSpawnStatus `json:"userSpawnStatusArr"`
	//RedScore           uint16                          `json:"redScore"`
	//BlueScore          uint16                          `json:"blueScore"`
}

func NewConnectionResponsePacket(qPort uint32, udpAddr string) *ConnectionResponsePacket {
	return &ConnectionResponsePacket{
		PKind:   'I',
		QPort:   qPort,
		UdpAddr: udpAddr,
		//UserSpawnStatusArr: userSpawnStatusArr,
		//RedScore:           redScore,
		//BlueScore:          blueScore,
	}
}

func (p *ConnectionResponsePacket) Serialize() []byte {
	return []byte{}
}
