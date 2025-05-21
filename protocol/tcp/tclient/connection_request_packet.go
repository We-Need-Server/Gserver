package tclient

type ConnectionRequestPacket struct {
	UserId uint32
}

func NewConnectionRequestPacket(userId uint32) *ConnectionRequestPacket {
	return &ConnectionRequestPacket{UserId: userId}
}
