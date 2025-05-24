package tserver

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
	return []byte{}
}
