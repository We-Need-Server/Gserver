package packet

import (
	"encoding/binary"
	"fmt"
)

type Packet struct {
	SEQ             uint32
	ACK             uint32
	QPort           uint32
	PKind           uint32
	Payload         []byte
	PayloadEndpoint int
}

func NewPacket(SEQ uint32, ACK uint32, QPort uint32, PKind uint32, Payload []byte, PayloadEndpoint int) *Packet {
	return &Packet{SEQ, ACK, QPort, PKind, Payload, PayloadEndpoint}
}

func ParsePacket(np []byte, endPoint int) *Packet {
	PKind := binary.BigEndian.Uint32(np[0:4])
	SEQ := binary.BigEndian.Uint32(np[4:8])
	ACK := binary.BigEndian.Uint32(np[8:12])
	QPort := binary.BigEndian.Uint32(np[12:16])
	Payload := np[16:endPoint]
	PayloadEndpoint := endPoint - 16
	fmt.Println(PKind, SEQ, ACK, QPort, Payload, PayloadEndpoint)
	return NewPacket(SEQ, ACK, QPort, PKind, Payload, PayloadEndpoint)
}
