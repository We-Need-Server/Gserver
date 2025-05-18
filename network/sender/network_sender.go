package sender

import (
	"WeNeedGameServer/packet"
	"net"
)

type Sender struct {
	connTable    *map[uint32]*net.UDPAddr
	nextSeqTable *map[uint32]uint32
	nChan        *chan packet.PacketI
	udpConn      *net.UDPConn
}

func NewSender(connTable *map[uint32]*net.UDPAddr, nextSeqTable *map[uint32]uint32, nChan *chan packet.PacketI, udpConn *net.UDPConn) *Sender {
	return &Sender{
		connTable:    connTable,
		nextSeqTable: nextSeqTable,
		nChan:        nChan,
		udpConn:      udpConn,
	}
}
