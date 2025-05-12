package client

import (
	"WeNeedGameServer/util"
	"fmt"
)

type EventPacket struct {
	pKind           uint32
	qPort           uint32
	SEQ             uint32
	Payload         []byte
	PayloadEndpoint int
}

func (p *EventPacket) GetPacketKind() uint32 {
	return p.pKind
}

func newEventPacket(pKind uint32, qPort uint32, SEQ uint32, Payload []byte, PayloadEndpoint int) *EventPacket {
	return &EventPacket{pKind, qPort, SEQ, Payload, PayloadEndpoint}
}

func ParseEventPacket(np []byte, endPoint int) *EventPacket {
	pKind := util.ConvertBinaryToUint32(np[0:4])
	qPort := util.ConvertBinaryToUint32(np[4:8])
	SEQ := util.ConvertBinaryToUint32(np[8:12])
	Payload := np[12:endPoint]
	PayloadEndpoint := endPoint - 12
	fmt.Println(pKind, qPort, SEQ, Payload, PayloadEndpoint)
	return newEventPacket(pKind, qPort, SEQ, Payload, PayloadEndpoint)
}

func (p *EventPacket) GetQPort() uint32 {
	return p.qPort
}

func (p *EventPacket) GetSEQ() uint32 {
	return p.SEQ
}
