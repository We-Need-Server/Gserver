package internal

import (
	"WeNeedGameServer/game_manager/internal/receiver"
	"WeNeedGameServer/game_manager/internal/sender"
	"WeNeedGameServer/game_manager/internal/types"
	"WeNeedGameServer/protocol/udp"
	"log"
	"net"
)

type GameNetwork struct {
	udpConnTable  map[uint32]*types.UdpUserConnStatus
	nextSeqTable  map[uint32]uint32
	nChan         chan udp.PacketI
	UdpReceiver   *receiver.UdpReceiver
	UdpSender     *sender.UdpSender
	udpConn       *net.UDPConn
	listenUdpAddr string
	findUserFunc  func(uint32) uint32
}

func NewGameNetwork(listenUdpAddr string, findUserFunc func(uint32) uint32) *GameNetwork {
	return &GameNetwork{
		udpConnTable:  make(map[uint32]*types.UdpUserConnStatus),
		nextSeqTable:  make(map[uint32]uint32),
		nChan:         make(chan udp.PacketI),
		UdpReceiver:   nil,
		UdpSender:     nil,
		udpConn:       nil,
		listenUdpAddr: listenUdpAddr,
		findUserFunc:  findUserFunc,
	}
}

func (n *GameNetwork) ReadyUdp() {
	udpServerPoint, udpResolveErr := net.ResolveUDPAddr("udp", n.listenUdpAddr)
	if udpResolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	udpLn, udpListenErr := net.ListenUDP("udp", udpServerPoint)
	if udpListenErr != nil {
		log.Panicln("리슨 오류")
	}
	n.udpConn = udpLn
	n.UdpReceiver = receiver.NewUdpReceiver(n.udpConnTable, n.nextSeqTable, n.nChan, n.udpConn, n.findUserFunc)
	n.UdpSender = sender.NewUdpSender(n.udpConnTable, n.nextSeqTable, n.nChan, n.udpConn)

}
