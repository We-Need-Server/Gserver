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
	tcpListener   *net.TCPListener
	tcpReceiver   *receiver.TcpReceiver
	listenUdpAddr string
	listenTcpAddr string
}

func NewNetwork(listenUdpAddr string, listenTcpAddr string) *Network {
	return &Network{
		connTable:     make(map[uint32]*net.UDPAddr),
		nextSeqTable:  make(map[uint32]uint32),
		nChan:         make(chan packet.PacketI),
		udpReceiver:   nil,
		udpSender:     nil,
		udpConn:       nil,
		tcpListener:   nil,
		tcpReceiver:   nil,
		listenUdpAddr: listenUdpAddr,
		listenTcpAddr: listenTcpAddr,
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
	n.udpReceiver = receiver.NewUdpReceiver(&n.connTable, &n.nextSeqTable, &n.nChan, n.udpConn)
	n.udpSender = sender.NewUdpSender(&n.connTable, &n.nextSeqTable, &n.nChan, n.udpConn)

	return n.udpReceiver, n.udpSender
}

func (n *Network) ReadyTcp() *receiver.TcpReceiver {
	tcpServerPoint, tcpResolveErr := net.ResolveTCPAddr("tcp", n.listenTcpAddr)
	if tcpResolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	tcpLn, tcpListenErr := net.ListenTCP("tcp", tcpServerPoint)
	if tcpListenErr != nil {
		log.Panicln("리슨 오류")
	}
	n.tcpListener = tcpLn
	n.tcpReceiver = receiver.NewTcpReceiver(tcpLn)

	return n.tcpReceiver
}
