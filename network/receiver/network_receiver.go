package receiver

import (
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/packet"
	"net"
)

type Receiver struct {
	connTable *map[uint32]*net.UDPAddr
	nQueue    *internal_type.Queue[packet.PacketI]
	udpConn   *net.UDPConn
}

func NewReceiver(connTable *map[uint32]*net.UDPAddr, nQueue *internal_type.Queue[packet.PacketI], udpConn *net.UDPConn) *Receiver {
	return &Receiver{connTable, nQueue, udpConn}
}
