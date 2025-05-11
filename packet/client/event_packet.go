package client

import (
	"WeNeedGameServer/util"
	"fmt"
)

type EventPacket struct {
	pKind           uint32
	qPort           uint32
	SEQ             uint32
	ACK             uint32
	Payload         []byte
	PayloadEndpoint int
}

func (p *EventPacket) GetPacketKind() uint32 {
	return p.pKind
}

func newEventPacket(pKind uint32, qPort uint32, SEQ uint32, ACK uint32, Payload []byte, PayloadEndpoint int) *EventPacket {
	return &EventPacket{pKind, qPort, SEQ, ACK, Payload, PayloadEndpoint}
}

func ParseEventPacket(np []byte, endPoint int) *EventPacket {
	pKind := util.ConvertBinaryToUint32(np[0:4])
	qPort := util.ConvertBinaryToUint32(np[4:8])
	SEQ := util.ConvertBinaryToUint32(np[8:12])
	ACK := util.ConvertBinaryToUint32(np[12:16])
	Payload := np[16:endPoint]
	PayloadEndpoint := endPoint - 16
	fmt.Println(pKind, qPort, SEQ, ACK, Payload, PayloadEndpoint)
	return newEventPacket(pKind, qPort, SEQ, ACK, Payload, PayloadEndpoint)
}

func (p *EventPacket) GetQPort() uint32 {
	return p.qPort
}
