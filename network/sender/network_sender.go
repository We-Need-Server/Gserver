package sender

import (
	"WeNeedGameServer/packet"
	"fmt"
	"log"
	"net"
)

type Sender struct {
	ConnTable    *map[uint32]*net.UDPAddr
	NextSeqTable *map[uint32]uint32
	NChan        *chan packet.PacketI
	udpConn      *net.UDPConn
}

func NewSender(connTable *map[uint32]*net.UDPAddr, nextSeqTable *map[uint32]uint32, nChan *chan packet.PacketI, udpConn *net.UDPConn) *Sender {
	return &Sender{
		ConnTable:    connTable,
		NextSeqTable: nextSeqTable,
		NChan:        nChan,
		udpConn:      udpConn,
	}
}

func (s *Sender) SendUDPPacket(b []byte, udpAddr *net.UDPAddr) (int, error) {
	fmt.Println(b)
	status, err := s.udpConn.WriteToUDP(b, udpAddr)
	if err != nil {
		log.Println("Failed to send message:", err)
		return status, err
	}
	return status, nil
}
