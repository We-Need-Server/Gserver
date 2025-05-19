package network

import (
	"WeNeedGameServer/network/receiver"
	"WeNeedGameServer/network/sender"
	"WeNeedGameServer/packet"
	"log"
	"net"
)

type Network struct {
	connTable     map[uint32]*net.UDPAddr
	nextSeqTable  map[uint32]uint32
	nChan         chan packet.PacketI
	udpReceiver   *receiver.UdpReceiver
	udpSender     *sender.UdpSender
	udpConn       *net.UDPConn
	listenUdpAddr string
	listenTcpAddr string
}

func NewNetwork(listenUdpAddr string) *Network {
	return &Network{
		connTable:     make(map[uint32]*net.UDPAddr),
		nextSeqTable:  make(map[uint32]uint32),
		nChan:         make(chan packet.PacketI),
		udpReceiver:   nil,
		udpSender:     nil,
		udpConn:       nil,
		listenUdpAddr: listenUdpAddr,
	}
}

func (n *Network) ReadyUDP() (*receiver.UdpReceiver, *sender.UdpSender) {
	UDPServerPoint, resolveErr := net.ResolveUDPAddr("udp", n.listenUdpAddr)
	if resolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	ln, listenErr := net.ListenUDP("udp", UDPServerPoint)
	if listenErr != nil {
		log.Panicln("리슨 오류")
	}
	n.udpConn = ln
	n.udpReceiver = receiver.NewUdpReceiver(&n.connTable, &n.nextSeqTable, &n.nChan, n.udpConn)
	n.udpSender = sender.NewUdpSender(&n.connTable, &n.nextSeqTable, &n.nChan, n.udpConn)
	return n.udpReceiver, n.udpSender
}
