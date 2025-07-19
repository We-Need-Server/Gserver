package tserver

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type GameOverPacket struct {
	ContentLength uint32 `json:"-"`
	PKind         uint8  `json:"-"`
}

func NewGameOverPacket() *GameOverPacket {
	return &GameOverPacket{PKind: 'O'}
}

func (p *GameOverPacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	p.ContentLength = uint32(len(data) + 1)
	result := make([]byte, 5+len(data))
	binary.LittleEndian.PutUint32(result[0:4], p.ContentLength)
	result[4] = p.PKind
	copy(result[5:], data)
	fmt.Println('O')
	fmt.Println(len(result))
	fmt.Println(p.ContentLength)
	fmt.Print(len(data))
	return result
}
