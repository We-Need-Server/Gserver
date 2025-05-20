package tcp_server

import "WeNeedGameServer/internal_type"

type ConnectionResponsePacket struct {
	PKind              uint8                           `json:"packetKind"`
	QPort              uint32                          `json:"qPort"`
	UdpAddr            string                          `json:"udpAddr"`
	UserSpawnStatusArr []internal_type.UserSpawnStatus `json:"userSpawnStatusArr"`
	RedScore           uint16                          `json:"redScore"`
	BlueScore          uint16                          `json:"blueScore"`
}

func NewConnectionResponsePacket(qPort uint32, udpAddr string, userSpawnStatusArr []internal_type.UserSpawnStatus, redScore uint16, blueScore uint16) *ConnectionResponsePacket {
	return &ConnectionResponsePacket{
		PKind:              'I',
		QPort:              qPort,
		UdpAddr:            udpAddr,
		UserSpawnStatusArr: userSpawnStatusArr,
		RedScore:           redScore,
		BlueScore:          blueScore,
	}
}
