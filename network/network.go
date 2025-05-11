package network

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/game/actor"
	"WeNeedGameServer/mediator"
	"WeNeedGameServer/packet"
	"log"
	"net"
)

type Network struct {
	ChanTable        map[uint32]chan packet.PacketI
	ConnTable        map[uint32]*net.UDPAddr
	NextSEQTable     map[uint32]uint32
	PacketActorTable map[uint32]*actor.PacketActor
	Ln               *net.UDPConn
	Game             *game.Game
	Mediator         *mediator.Mediator
}

func NewNetwork(Game *game.Game) *Network {
	return &Network{
		ChanTable:        make(map[uint32]chan packet.PacketI),
		ConnTable:        make(map[uint32]*net.UDPAddr),
		PacketActorTable: make(map[uint32]*actor.PacketActor),
		Game:             Game,
		Mediator:         nil,
	}
}

func (n *Network) Start() {
	UDPServerPoint, resolveErr := net.ResolveUDPAddr("udp", ":20000")
	if resolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	ln, listenErr := net.ListenUDP("udp", UDPServerPoint)
	if listenErr != nil {
		log.Panicln("리슨 오류")
	}
	n.Ln = ln
	readBuffer := make([]byte, 2048)
	for {
		readCount, addr, err := ln.ReadFromUDP(readBuffer)
		if err != nil {
			log.Panicln("잘못된 요청")
		}
		n.handlePacket(readBuffer, readCount, addr)
	}
}

func (n *Network) Register(m *mediator.Mediator) {
	n.Mediator = m
}

func (n *Network) Send(receiverName string, message interface{}) {
	n.Mediator.Notify("network", receiverName, message)
}

func (n *Network) Receive(senderName string, message interface{}) {
}

func (n *Network) handlePacket(clientPacket []byte, endPoint int, userAddr *net.UDPAddr) {
	data, err := packet.ParsePacketByKind(clientPacket, endPoint)
	if err != nil {
		log.Panicln("잘못된 요청")
	}

	if QPort := n.ConnTable[data.GetQPort()]; QPort == nil {
		n.tempHandleNewConnection(data.GetQPort(), userAddr)
	}
	n.throwData(data, userAddr)
}

//func (n *Network) handleNewConnection(QPort uint32, userAddr string) bool {
//	checkUser := api.CheckUserValidation(userAddr)
//	if checkUser {
//		n.IPTable[QPort] = userAddr
//		n.QPortTable[userAddr] = QPort
//		n.ChanTable[QPort] = make(chan *packet.Packet)
//		packetActor := actor.NewPacketActor(1, QPort, userAddr, n.ChanTable[QPort], n.Game.AddPlayer(QPort))
//		n.PacketActorTable[QPort] = packetActor
//		go n.PacketActorTable[QPort].ProcessLoopPacket()
//	}
//	return checkUser
//}

func (n *Network) throwData(data packet.PacketI, userAddr *net.UDPAddr) {
	if n.ConnTable[data.GetQPort()] != nil || n.NextSEQTable[data.GetQPort()] == data.GetSEQ() {
		n.NextSEQTable[data.GetQPort()] += 1
		n.ChanTable[data.GetQPort()] <- data
	}
}

// handleConnection을 실행하도록 한다

func (n *Network) tempHandleNewConnection(QPort uint32, userAddr *net.UDPAddr) {
	n.ChanTable[QPort] = make(chan packet.PacketI)
	n.ConnTable[QPort] = userAddr
	n.NextSEQTable[QPort] = 1
	packetActor := actor.NewPacketActor(1, QPort, userAddr, n.ChanTable[QPort], n.Game.AddPlayer(QPort))
	n.PacketActorTable[QPort] = packetActor
	go n.PacketActorTable[QPort].ProcessLoopPacket()
}
