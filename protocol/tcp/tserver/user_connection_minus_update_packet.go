package tserver

type UserConnectionMUpdatePacket struct {
	PKind    uint8    `json:"packetKind"`
	UserList []uint32 `json:"userList"`
}

func NewUserConnectionMUpdatePacket(userList []uint32) *UserConnectionMUpdatePacket {
	return &UserConnectionMUpdatePacket{
		PKind:    'M',
		UserList: userList,
	}
}

func (p *UserConnectionMUpdatePacket) Serialize() []byte {
	return []byte{}
}
