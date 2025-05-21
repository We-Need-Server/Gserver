package internal

import (
	"WeNeedGameServer/protocol/udp"
	"WeNeedGameServer/round/internal/receiver"
	"WeNeedGameServer/round/internal/sender"
	"log"
	"net"
)

type Network struct {
	udpConnTable map[uint32]*net.UDPAddr
	nextSeqTable map[uint32]uint32
	nChan        chan udp.PacketI
	udpReceiver  *receiver.UdpReceiver
	udpSender    *sender.UdpSender
	udpConn      *net.UDPConn
	tcpListener  *net.TCPListener
	//tcpReceiver   *internal.TcpReceiver
	//tcpSender     *internal.TcpSender
	listenUdpAddr string
	listenTcpAddr string
}

func NewNetwork(listenUdpAddr string, listenTcpAddr string) *Network {
	return &Network{
		udpConnTable: make(map[uint32]*net.UDPAddr),
		nextSeqTable: make(map[uint32]uint32),
		nChan:        make(chan udp.PacketI),
		udpReceiver:  nil,
		udpSender:    nil,
		udpConn:      nil,
		tcpListener:  nil,
		//tcpReceiver:   nil,
		//tcpSender:     nil,
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

//
//func (n *Network) communicateTcpSender(message internal.TcpReceiverMessage) {
//	n.tcpSender.ProcessMessage(message)
//}

//func (n *Network) ReadyTcp(blueTeamDb map[uint32]*user.User, redTeamDb map[uint32]*user.User, loginFunc func(uint32, *net.TCPAddr) (uint32, error)) (*internal.TcpReceiver, *internal.TcpSender) {
//	tcpServerPoint, tcpResolveErr := net.ResolveTCPAddr("tcp", n.listenTcpAddr)
//	if tcpResolveErr != nil {
//		log.Panicln("네트워크 리졸버 오류")
//	}
//	tcpLn, tcpListenErr := net.ListenTCP("tcp", tcpServerPoint)
//	if tcpListenErr != nil {
//		log.Panicln("리슨 오류")
//	}
//	n.tcpListener = tcpLn
//	n.tcpReceiver = internal.NewTcpReceiver(tcpLn, loginFunc, n.communicateTcpSender)
//	n.tcpSender = internal.NewTcpSender(n.listenUdpAddr, blueTeamDb, redTeamDb)
//	return n.tcpReceiver, n.tcpSender
//}
