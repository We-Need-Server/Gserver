package network

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/packet"
	"WeNeedGameServer/packet/actor"
	"fmt"
	"net"
)

type Network struct {
	IPTable          map[uint32]string
	QPortTable       map[string]uint32
	ChanTable        map[uint32]chan *packet.Packet
	PacketActorTable map[uint32]*actor.PacketActor
	Game             *game.Game
}

func NewNetwork(Game *game.Game) *Network {
	return &Network{
		IPTable:          make(map[uint32]string),
		QPortTable:       make(map[string]uint32),
		ChanTable:        make(map[uint32]chan *packet.Packet),
		PacketActorTable: make(map[uint32]*actor.PacketActor),
		Game:             Game,
	}
}

func (n *Network) TempStart() {
	// 테스트 초기화 부분
	n.handleNewConnection(32, "127.0.0.1:4284")
	//n.handleNewConnection(42, "127.0.0.1:3030")
	UDPServerPoint, resolveErr := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if resolveErr != nil {
		fmt.Println("네트워크 리졸버 오류")
	}
	ln, listenErr := net.ListenUDP("udp", UDPServerPoint)
	if listenErr != nil {
		fmt.Println("리슨 오류")
	}
	readBuffer := make([]byte, 2048)
	for {
		readCount, addr, err := ln.ReadFromUDP(readBuffer)
		if err != nil {
			fmt.Println("잘못된 요청")
		}
		n.TempHandlePacket(readBuffer, readCount, addr.String())
	}
}

func (n *Network) TempHandlePacket(clientPacket []byte, endPoint int, userAddr string) {
	data := packet.ParsePacket(clientPacket, endPoint)
	fmt.Println(data)
	fmt.Println(n.ChanTable)
	n.TempthrowData(data)
	//n.throwData(data, userAddr)
}

// 일단 임시적으로 패킷을
// SEQ
// ACK
// QPORT
// PayLoad로 구성한다.

func (n *Network) TempthrowData(data *packet.Packet) {
	if data.QPort == 32 || data.QPort == 42 {
		n.ChanTable[data.QPort] <- data
	}
}

func (n *Network) throwData(data *packet.Packet, userAddr string) {
	if n.IPTable[data.QPort] == userAddr && n.QPortTable[userAddr] == data.QPort {
		n.ChanTable[data.QPort] <- data
	}
}

// handleConnection을 실행하도록 한다

func (n *Network) handleNewConnection(QPort uint32, userAddr string) {
	n.IPTable[QPort] = userAddr
	n.QPortTable[userAddr] = QPort
	n.ChanTable[QPort] = make(chan *packet.Packet)
	packetActor := actor.NewPacketActor(1, QPort, userAddr, n.ChanTable[QPort], n.Game.AddPlayer(QPort))
	n.PacketActorTable[QPort] = packetActor
	go n.PacketActorTable[QPort].ProcessLoopPacket()
}
