package tserver

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type UserConnectionMUpdatePacket struct {
	ContentLength uint32   `json:"-"`
	PKind         uint8    `json:"-"`
	UserList      []uint32 `json:"userList"`
}

func NewUserConnectionMUpdatePacket(userList []uint32) *UserConnectionMUpdatePacket {
	return &UserConnectionMUpdatePacket{
		PKind:    'M',
		UserList: userList,
	}
}

func (p *UserConnectionMUpdatePacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	p.ContentLength = uint32(len(data) + 1)
	result := make([]byte, 5+len(data))
	binary.LittleEndian.PutUint32(result[0:4], p.ContentLength)
	result[4] = p.PKind
	copy(result[5:], data)
	fmt.Println(len(result))
	fmt.Println(p.ContentLength)
	fmt.Print(len(data))
	return result
}
