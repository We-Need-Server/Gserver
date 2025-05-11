package network

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/game/actor"
	"WeNeedGameServer/packet"
	"WeNeedGameServer/packet/client"
	"log"
	"math"
	"net"
)

type Network struct {
	IPTable map[uint32]*net.UDPAddr
	//QPortTable       map[*net.UDPAddr]uint32
	ChanTable        map[uint32]chan *packet.PacketI
	ConnTable        map[uint32]*net.UDPAddr
	PacketActorTable map[uint32]*actor.PacketActor
	Ln               *net.UDPConn
	Game             *game.Game
}

func NewNetwork(Game *game.Game) *Network {
	return &Network{
		IPTable: make(map[uint32]*net.UDPAddr),
		//QPortTable:       make(map[*net.UDPAddr]uint32),
		ChanTable:        make(map[uint32]chan *packet.PacketI),
		ConnTable:        make(map[uint32]*net.UDPAddr),
		PacketActorTable: make(map[uint32]*actor.PacketActor),
		Game:             Game,
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

func (n *Network) handlePacket(clientPacket []byte, endPoint int, userAddr *net.UDPAddr) {
	pKind, data, err := packet.ParsePacketByKind(clientPacket, endPoint)
	if err != nil || pKind != uint32(math.MaxUint32) {
		log.Panicln("잘못된 요청")
	}
	switch pKind {
	case 41:
		data = data.(*client.EventPacket)
	case 46:
		data = data.(*client.TickIPacket)
	case 50:
		data = data.(*client.TickRPacket)
	}
	if QPort := n.IPTable[data.GetQPort()]; QPort == nil {
		n.tempHandleNewConnection(data.GetQPort(), userAddr)
	}
	//if userAddr := n.IPTable[data.QPort]; userAddr == "" {
	//	if checkUser := n.handleNewConnection(data.QPort, userAddr); checkUser {
	//		log.Println(data.QPort, userAddr, "was denied because of not certificated user")
	//		return
	//	}
	//}
	n.throwData(pKind, data, userAddr)
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

//func (n *Network) TempStart() {
//	// 테스트 초기화 부분
//	n.tempHandleNewConnection(32, "127.0.0.1:4284")
//	UDPServerPoint, resolveErr := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
//	if resolveErr != nil {
//		fmt.Println("네트워크 리졸버 오류")
//	}
//	ln, listenErr := net.ListenUDP("udp", UDPServerPoint)
//	if listenErr != nil {
//		fmt.Println("리슨 오류")
//	}
//	readBuffer := make([]byte, 2048)
//	for {
//		readCount, addr, err := ln.ReadFromUDP(readBuffer)
//		if err != nil {
//			fmt.Println("잘못된 요청")
//		}
//		n.tempHandlePacket(readBuffer, readCount, addr.String())
//	}
//}

//func (n *Network) tempHandlePacket(clientPacket []byte, endPoint int, userAddr string) {
//	data := packet.ParseEventPacket(clientPacket, endPoint)
//	fmt.Println(data)
//	fmt.Println(n.ChanTable)
//	n.TempthrowData(data)
//	//n.throwData(data, userAddr)
//}

// 일단 임시적으로 패킷을
// SEQ
// ACK
// QPORT
// PayLoad로 구성한다.

//func (n *Network) TempthrowData(data *packet.EventPacket) {
//	if data.QPort == 32 || data.QPort == 42 {
//		n.ChanTable[data.QPort] <- data
//	}
//}

func (n *Network) throwData(data packet.PacketI, userAddr *net.UDPAddr) {
	if n.IPTable[data.GetQPort()] != nil {
		n.ChanTable[data.GetQPort()] <- &data
	}
}

// handleConnection을 실행하도록 한다

func (n *Network) tempHandleNewConnection(QPort uint32, userAddr *net.UDPAddr) {
	n.IPTable[QPort] = userAddr
	//n.QPortTable[userAddr] = QPort
	n.ChanTable[QPort] = make(chan *packet.PacketI)
	//clientAddr, err := net.ResolveUDPAddr("udp", userAddr)
	//if err != nil {
	//	log.Println("ClientAddr Error for", userAddr, ":", err)
	//}
	//clientConn, clientConnErr := net.DialUDP("udp", nil, clientAddr)
	//if clientConnErr != nil {
	//	log.Println("Client Connection Error for", userAddr, ":", clientConnErr)
	//} else {
	n.ConnTable[QPort] = userAddr
	//}
	packetActor := actor.NewPacketActor(1, QPort, userAddr, n.ChanTable[QPort], n.Game.AddPlayer(QPort))
	n.PacketActorTable[QPort] = packetActor
	go n.PacketActorTable[QPort].ProcessLoopPacket()
}
