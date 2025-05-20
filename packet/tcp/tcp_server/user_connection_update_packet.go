package tcp_server

type UserConnectionUpdatePacket struct {
	PacketKind uint8    `json:"packetKind"`
	UserList   []uint32 `json:"userList"`
}

func NewUserConnectionUpdatePacket(userList []uint32) *UserConnectionUpdatePacket {
	return &UserConnectionUpdatePacket{
		PacketKind: 'U',
		UserList:   userList,
	}
}
