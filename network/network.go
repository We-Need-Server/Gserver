package network

import (
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/network/receiver"
	"WeNeedGameServer/network/sender"
	"WeNeedGameServer/packet"
	"log"
	"net"
)

type Network struct {
	connTable  map[uint32]*net.UDPAddr
	nQueue     *internal_type.Queue[packet.PacketI]
	ur         *receiver.Receiver
	us         *sender.Sender
	udpConn    *net.UDPConn
	listenAddr string
}

func NewNetwork(listenAddr string) *Network {
	return &Network{
		connTable:  make(map[uint32]*net.UDPAddr),
		nQueue:     internal_type.NewQueue[packet.PacketI](),
		ur:         nil,
		us:         nil,
		udpConn:    nil,
		listenAddr: listenAddr,
	}
}

func (n *Network) StartUDP() (*receiver.Receiver, *sender.Sender) {
	UDPServerPoint, resolveErr := net.ResolveUDPAddr("udp", n.listenAddr)
	if resolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	ln, listenErr := net.ListenUDP("udp", UDPServerPoint)
	if listenErr != nil {
		log.Panicln("리슨 오류")
	}
	n.udpConn = ln
	n.ur = receiver.NewReceiver(&n.connTable, n.nQueue, n.udpConn)
	n.us = sender.NewSender(&n.connTable, n.nQueue, n.udpConn)
	return n.ur, n.us
}
