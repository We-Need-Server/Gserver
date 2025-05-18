package sender

import (
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/packet"
	"net"
)

type Sender struct {
	connTable *map[uint32]*net.UDPAddr
	nQueue    *internal_type.Queue[*packet.PacketI]
	udpConn   *net.UDPConn
}

func NewSender(connTable *map[uint32]*net.UDPAddr, nQueue *internal_type.Queue[*packet.PacketI], udpConn *net.UDPConn) *Sender {
	return &Sender{connTable, nQueue, udpConn}
}
