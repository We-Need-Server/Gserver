package tserver

import "encoding/json"

type UserConnectionPUpdatePacket struct {
	PKind    uint8    `json:"packetKind"`
	UserList []uint32 `json:"userList"`
}

func NewUserConnectionPUpdatePacket(userList []uint32) *UserConnectionPUpdatePacket {
	return &UserConnectionPUpdatePacket{
		PKind:    'P',
		UserList: userList,
	}
}

func (p *UserConnectionPUpdatePacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	return data
}
