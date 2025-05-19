package network

import (
	"WeNeedGameServer/network/receiver"
	"WeNeedGameServer/network/sender"
	"WeNeedGameServer/packet"
	"log"
	"net"
)

type Network struct {
	connTable    map[uint32]*net.UDPAddr
	nextSeqTable map[uint32]uint32
	nChan        chan packet.PacketI
	udpReceiver  *receiver.Receiver
	udpSender    *sender.Sender
	udpConn      *net.UDPConn
	listenAddr   string
}

func NewNetwork(listenAddr string) *Network {
	return &Network{
		connTable:    make(map[uint32]*net.UDPAddr),
		nextSeqTable: make(map[uint32]uint32),
		nChan:        make(chan packet.PacketI),
		udpReceiver:  nil,
		udpSender:    nil,
		udpConn:      nil,
		listenAddr:   listenAddr,
	}
}

func (n *Network) ReadyUDP() (*receiver.Receiver, *sender.Sender) {
	UDPServerPoint, resolveErr := net.ResolveUDPAddr("udp", n.listenAddr)
	if resolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	ln, listenErr := net.ListenUDP("udp", UDPServerPoint)
	if listenErr != nil {
		log.Panicln("리슨 오류")
	}
	n.udpConn = ln
	n.udpReceiver = receiver.NewReceiver(&n.connTable, &n.nextSeqTable, &n.nChan, n.udpConn)
	n.udpSender = sender.NewSender(&n.connTable, &n.nextSeqTable, &n.nChan, n.udpConn)
	return n.udpReceiver, n.udpSender
}
