package udp_client

import (
	"WeNeedGameServer/util"
	"fmt"
)

type EventPacket struct {
	pKind           uint8
	qPort           uint32
	SEQ             uint32
	Payload         []byte
	PayloadEndpoint int
}

func (p *EventPacket) GetPacketKind() uint8 {
	return p.pKind
}

func newEventPacket(pKind uint8, qPort uint32, SEQ uint32, Payload []byte, PayloadEndpoint int) *EventPacket {
	return &EventPacket{pKind, qPort, SEQ, Payload, PayloadEndpoint}
}

func ParseEventPacket(np []byte, endPoint int) *EventPacket {
	pKind := np[0]
	qPort := util.ConvertBinaryToUint32(np[1:5])
	SEQ := util.ConvertBinaryToUint32(np[5:9])
	Payload := np[9:endPoint]
	PayloadEndpoint := endPoint - 9
	fmt.Println(pKind, qPort, SEQ, Payload, PayloadEndpoint)
	return newEventPacket(pKind, qPort, SEQ, Payload, PayloadEndpoint)
}

func (p *EventPacket) GetQPort() uint32 {
	return p.qPort
}

func (p *EventPacket) GetSEQ() uint32 {
	return p.SEQ
}
