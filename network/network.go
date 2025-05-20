package network

import (
	"WeNeedGameServer/db/user"
	"WeNeedGameServer/network/receiver"
	"WeNeedGameServer/network/sender"
	"WeNeedGameServer/packet/udp"
	"log"
	"net"
)

type Network struct {
	tcpConnTable  map[uint32]*net.Conn
	udpConnTable  map[uint32]*net.UDPAddr
	nextSeqTable  map[uint32]uint32
	nChan         chan udp.PacketI
	udpReceiver   *receiver.UdpReceiver
	udpSender     *sender.UdpSender
	udpConn       *net.UDPConn
	tcpListener   *net.TCPListener
	tcpReceiver   *receiver.TcpReceiver
	tcpSender     *sender.TcpSender
	listenUdpAddr string
	listenTcpAddr string
}

func NewNetwork(listenUdpAddr string, listenTcpAddr string) *Network {
	return &Network{
		tcpConnTable:  make(map[uint32]*net.Conn),
		udpConnTable:  make(map[uint32]*net.UDPAddr),
		nextSeqTable:  make(map[uint32]uint32),
		nChan:         make(chan udp.PacketI),
		udpReceiver:   nil,
		udpSender:     nil,
		udpConn:       nil,
		tcpListener:   nil,
		tcpReceiver:   nil,
		tcpSender:     nil,
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
	n.udpReceiver = receiver.NewUdpReceiver(&n.udpConnTable, &n.nextSeqTable, &n.nChan, n.udpConn)
	n.udpSender = sender.NewUdpSender(&n.udpConnTable, &n.nextSeqTable, &n.nChan, n.udpConn)

	return n.udpReceiver, n.udpSender
}

func (n *Network) communicateTcpSender(message receiver.TcpReceiverMessage) {
	n.tcpSender.ProcessMessage(message)
}

func (n *Network) ReadyTcp(redTeamDb map[uint32]*user.User, blueTeamDb map[uint32]*user.User, loginFunc func(uint32, *net.TCPAddr) (uint32, error)) *receiver.TcpReceiver {
	tcpServerPoint, tcpResolveErr := net.ResolveTCPAddr("tcp", n.listenTcpAddr)
	if tcpResolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	tcpLn, tcpListenErr := net.ListenTCP("tcp", tcpServerPoint)
	if tcpListenErr != nil {
		log.Panicln("리슨 오류")
	}
	n.tcpListener = tcpLn
	n.tcpSender = sender.NewTcpSender(n.listenUdpAddr, n.tcpConnTable)
	n.tcpReceiver = receiver.NewTcpReceiver(tcpLn, n.tcpConnTable, loginFunc, n.communicateTcpSender)

	return n.tcpReceiver
}
