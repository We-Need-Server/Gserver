package tserver

import (
	"encoding/binary"
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
	p.ContentLength = 1
	result := make([]byte, 5)
	binary.LittleEndian.PutUint32(result[0:4], p.ContentLength)
	result[4] = p.PKind
	//copy(result[5:], data)
	fmt.Println("O")
	fmt.Println(p.ContentLength)
	fmt.Println(result)
	return result
}
