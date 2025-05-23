package legacy

//import (
//	"WeNeedGameServer/game/legacy"
//	"WeNeedGameServer/game/legacy/actor"
//	"WeNeedGameServer/legacy/internal_type"
//	"WeNeedGameServer/legacy/mediator"
//	"WeNeedGameServer/protocol/udp"
//	"log"
//	"net"
//)
//
//internal_types Network struct {
//	chanTable        map[uint32]chan udp.PacketI
//	ConnTable        map[uint32]*net.UDPAddr
//	nextSEQTable     map[uint32]uint32
//	packetActorTable map[uint32]*actor.PacketActor
//	ln               *net.UDPConn
//	game             *legacy.Game
//	Mediator         *mediator.Mediator
//}
//
//func NewNetwork(game *legacy.Game) *Network {
//	return &Network{
//		chanTable:        make(map[uint32]chan udp.PacketI),
//		ConnTable:        make(map[uint32]*net.UDPAddr),
//		nextSEQTable:     make(map[uint32]uint32),
//		packetActorTable: make(map[uint32]*actor.PacketActor),
//		game:             game,
//		Mediator:         nil,
//	}
//}
//
//func (n *Network) Start() {
//	UDPServerPoint, resolveErr := net.ResolveUDPAddr("udp", ":20000")
//	if resolveErr != nil {
//		log.Panicln("네트워크 리졸버 오류")
//	}
//	ln, listenErr := net.ListenUDP("udp", UDPServerPoint)
//	if listenErr != nil {
//		log.Panicln("리슨 오류")
//	}
//	n.ln = ln
//	readBuffer := make([]byte, 2048)
//	for {
//		readCount, addr, err := ln.ReadFromUDP(readBuffer)
//		if err != nil {
//			log.Panicln("잘못된 요청")
//		}
//		n.handlePacket(readBuffer, readCount, addr)
//	}
//}
//
//func (n *Network) Register(m *mediator.Mediator) {
//	n.Mediator = m
//}
//
//func (n *Network) Send(receiverName string, message interface{}) {
//	n.Mediator.Notify("network", receiverName, message)
//}
//
//func (n *Network) Receive(senderName string, message interface{}) {
//}
//
//func (n *Network) handlePacket(clientPacket []byte, endPoint int, userAddr *net.UDPAddr) {
//	data, err := udp.ParsePacketByKind(clientPacket, endPoint)
//	if err != nil {
//		log.Panicln("잘못된 요청")
//	}
//
//	if QPort := n.ConnTable[data.GetQPort()]; QPort == nil {
//		n.tempHandleNewConnection(data.GetQPort(), userAddr)
//	}
//	n.throwData(data, userAddr)
//}
//
////func (n *Network) handleNewConnection(QPort uint32, userAddr string) bool {
////	checkUser := api.CheckUserValidation(userAddr)
////	if checkUser {
////		n.IPTable[QPort] = userAddr
////		n.QPortTable[userAddr] = QPort
////		n.ChanTable[QPort] = make(chan *packet.Packet)
////		packetActor := actor.NewPacketActor(1, QPort, userAddr, n.ChanTable[QPort], n.Game.AddPlayer(QPort))
////		n.PacketActorTable[QPort] = packetActor
////		go n.PacketActorTable[QPort].ProcessLoopPacket()
////	}
////	return checkUser
////}
//
//func (n *Network) throwData(data udp.ClientPacketI, userAddr *net.UDPAddr) {
//	if n.ConnTable[data.GetQPort()] != nil || n.nextSEQTable[data.GetQPort()] == data.GetSEQ() {
//		n.nextSEQTable[data.GetQPort()] += 1
//		n.Send("game_tick", internal_type.NewSEQData(data.GetQPort(), data.GetSEQ()))
//		n.chanTable[data.GetQPort()] <- data
//	}
//}
//
//// handleConnection을 실행하도록 한다
//
//func (n *Network) tempHandleNewConnection(qPort uint32, userAddr *net.UDPAddr) {
//	n.chanTable[qPort] = make(chan udp.PacketI)
//	n.ConnTable[qPort] = userAddr
//	n.nextSEQTable[qPort] = 1
//	// 이 부분에 대해서 mediator로 게임 객체에 전달하게 하여 게임 객체를 네트워크 객체가 안 가지도록 할 수 있음
//	//packetActor := actor.NewPacketActor(qPort, userAddr, n.chanTable[qPort], n.game.AddPlayer(qPort))
//	packetActor := actor.NewPacketActor(qPort, userAddr, n.chanTable[qPort])
//	n.Send("game", qPort)
//	if _, err := n.Mediator.Register("actor", packetActor); err != nil {
//		log.Panicln("메디에이터 등록 실패")
//	}
//	n.packetActorTable[qPort] = packetActor
//	packetActor.Send("game_tick", packetActor.QPort)
//	go n.packetActorTable[qPort].ProcessLoopPacket()
//}
//
//func (n *Network) SendUDPPacket(b []byte, udpAddr *net.UDPAddr) (int, error) {
//	status, err := n.ln.WriteToUDP(b, udpAddr)
//	if err != nil {
//		log.Println("Failed to send message:", err)
//		return status, err
//	}
//	return status, nil
//}
