package tserver

import (
	"encoding/binary"
	"encoding/json"
)

type RoundStartPacket struct {
	ContentLength uint32 `json:"-"`
	PKind         uint8  `json:"-"`
}

func NewRoundStartPacket() *RoundStartPacket {
	return &RoundStartPacket{
		PKind: 'S',
	}
}

func (p *RoundStartPacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	p.ContentLength = uint32(len(data) + 1)
	result := make([]byte, 9+len(data))
	binary.LittleEndian.PutUint32(result[0:4], p.ContentLength)
	result[4] = p.PKind
	copy(result[5:len(result)-4], data)
	copy(result[len(result)-4:], "\r\n\r\n")
	return result
}
