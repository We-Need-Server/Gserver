package sender

import (
	"WeNeedGameServer/protocol/udp"
	"fmt"
	"log"
	"net"
)

type UdpSender struct {
	ConnTable    *map[uint32]*net.UDPAddr
	NextSeqTable *map[uint32]uint32
	NChan        *chan udp.PacketI
	udpConn      *net.UDPConn
}

func NewUdpSender(connTable *map[uint32]*net.UDPAddr, nextSeqTable *map[uint32]uint32, nChan *chan udp.PacketI, udpConn *net.UDPConn) *UdpSender {
	return &UdpSender{
		ConnTable:    connTable,
		NextSeqTable: nextSeqTable,
		NChan:        nChan,
		udpConn:      udpConn,
	}
}

func (s *UdpSender) SendUdpPacket(b []byte, udpAddr *net.UDPAddr) (int, error) {
	fmt.Println(b)
	status, err := s.udpConn.WriteToUDP(b, udpAddr)
	if err != nil {
		log.Println("Failed to send message:", err)
		return status, err
	}
	return status, nil
}
