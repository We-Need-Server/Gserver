package tserver

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
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
	result := make([]byte, 5+len(data))
	binary.LittleEndian.PutUint32(result[0:4], p.ContentLength)
	result[4] = p.PKind
	copy(result[5:], data)
	fmt.Println('S')
	fmt.Println(len(result))
	fmt.Println(p.ContentLength)
	fmt.Print(len(data))
	return result
}
