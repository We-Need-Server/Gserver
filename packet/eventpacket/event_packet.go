package eventpacket

import (
	"WeNeedGameServer/util"
	"fmt"
)

type EventPacket struct {
	QPort           uint32
	SEQ             uint32
	ACK             uint32
	Payload         []byte
	PayloadEndpoint int
}

func (*EventPacket) GetPacketKind() uint32 {
	return 46
}

func newEventPacket(QPort uint32, SEQ uint32, ACK uint32, Payload []byte, PayloadEndpoint int) *EventPacket {
	return &EventPacket{QPort, SEQ, ACK, Payload, PayloadEndpoint}
}

func ParseEventPacket(np []byte, endPoint int) *EventPacket {
	QPort := util.ConvertBinaryToUint32(np[4:8])
	SEQ := util.ConvertBinaryToUint32(np[8:12])
	ACK := util.ConvertBinaryToUint32(np[12:16])
	Payload := np[16:endPoint]
	PayloadEndpoint := endPoint - 16
	fmt.Println(QPort, SEQ, ACK, Payload, PayloadEndpoint)
	return newEventPacket(QPort, SEQ, ACK, Payload, PayloadEndpoint)
}

func (p *EventPacket) GetQPort() uint32 {
	return p.QPort
}
