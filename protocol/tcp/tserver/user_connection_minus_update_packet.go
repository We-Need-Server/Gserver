package tserver

import "encoding/json"

type UserConnectionMUpdatePacket struct {
	PKind    uint8    `json:"-"`
	UserList []uint32 `json:"userList"`
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
	result := make([]byte, 1+len(data))
	result[0] = p.PKind
	copy(result[1:], data)

	return result
}
