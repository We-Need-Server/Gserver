package internal

import (
	"WeNeedGameServer/protocol/udp"
	"WeNeedGameServer/round/internal/receiver"
	"WeNeedGameServer/round/internal/sender"
	"log"
	"net"
)

type Network struct {
	udpConnTable  map[uint32]*net.UDPAddr
	nextSeqTable  map[uint32]uint32
	nChan         chan udp.PacketI
	udpReceiver   *receiver.UdpReceiver
	udpSender     *sender.UdpSender
	udpConn       *net.UDPConn
	listenUdpAddr string
}

func NewNetwork(listenUdpAddr string) *Network {
	return &Network{
		udpConnTable:  make(map[uint32]*net.UDPAddr),
		nextSeqTable:  make(map[uint32]uint32),
		nChan:         make(chan udp.PacketI),
		udpReceiver:   nil,
		udpSender:     nil,
		udpConn:       nil,
		listenUdpAddr: listenUdpAddr,
	}
}

func (n *Network) ReadyUdp() (*receiver.UdpReceiver, *sender.UdpSender) {
	udpServerPoint, udpResolveErr := net.ResolveUDPAddr("udp", n.listenUdpAddr)
	if udpResolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	udpLn, udpListenErr := net.ListenUDP("udp", udpServerPoint)
	if udpListenErr != nil {
		log.Panicln("리슨 오류")
	}
	n.udpConn = udpLn
	n.udpReceiver = receiver.NewUdpReceiver(&n.udpConnTable, &n.nextSeqTable, &n.nChan, n.udpConn)
	n.udpSender = sender.NewUdpSender(&n.udpConnTable, &n.nextSeqTable, &n.nChan, n.udpConn)

	return n.udpReceiver, n.udpSender
}
