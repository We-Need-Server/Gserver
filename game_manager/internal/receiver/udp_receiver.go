package receiver

import (
	"WeNeedGameServer/game_manager/internal/actor"
	"WeNeedGameServer/protocol/udp"
	"fmt"
	"log"
	"net"
)

type UdpReceiver struct {
	chanTable         map[uint32]chan udp.PacketI
	connTable         map[uint32]*net.UDPAddr
	networkActorTable map[uint32]*actor.UdpActor
	nextSeqTable      map[uint32]uint32
	nChan             chan udp.PacketI
	nChanManager      *UdpChanManager
	udpConn           *net.UDPConn
}

func NewUdpReceiver(connTable map[uint32]*net.UDPAddr, nextSeqTable map[uint32]uint32, nChan chan udp.PacketI, udpConn *net.UDPConn) *UdpReceiver {
	return &UdpReceiver{
		chanTable:         make(map[uint32]chan udp.PacketI),
		connTable:         connTable,
		networkActorTable: make(map[uint32]*actor.UdpActor),
		nextSeqTable:      nextSeqTable,
		nChan:             nChan,
		nChanManager:      NewUdpChanManager(nChan),
		udpConn:           udpConn,
	}
}

func (r *UdpReceiver) StartUdp() {
	go r.nChanManager.StartChanManager()
	readBuffer := make([]byte, 2048)
	for {
		readCount, addr, err := r.udpConn.ReadFromUDP(readBuffer)
		if err != nil {
			log.Panicln("잘못된 요청")
		}
		r.handlePacket(readBuffer, readCount, addr)
	}
}

func (r *UdpReceiver) handlePacket(clientPacket []byte, endPoint int, userAddr *net.UDPAddr) {
	data, err := udp.ParsePacketByKind(clientPacket, endPoint)
	if err != nil {
		log.Panicln("잘못된 요청")
	}

	if QPort := r.connTable[data.GetQPort()]; QPort == nil {
		r.tempHandleNewConnection(data.GetQPort(), userAddr)
	}
	r.throwData(data)
}

func (r *UdpReceiver) throwData(data udp.ClientPacketI) {
	if r.connTable[data.GetQPort()] != nil || r.nextSeqTable[data.GetQPort()] == data.GetSEQ() {
		r.nextSeqTable[data.GetQPort()] += 1
		fmt.Println("R 패킷 왔니")
		fmt.Println(data.GetPacketKind())
		if data.GetPacketKind() == 'N' {
			r.chanTable[data.GetQPort()] <- data
		} else {
			r.nChan <- data
		}
	}
}

func (r *UdpReceiver) tempHandleNewConnection(qPort uint32, userAddr *net.UDPAddr) {
	r.chanTable[qPort] = make(chan udp.PacketI)
	r.connTable[qPort] = userAddr
	r.nextSeqTable[qPort] = 1
	networkActor := actor.NewUdpActor(qPort, userAddr, r.chanTable[qPort], &r.nChanManager.CmChan)
	r.networkActorTable[qPort] = networkActor
	go r.networkActorTable[qPort].ProcessLoopPacket()
}
