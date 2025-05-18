package receiver

import (
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/network/actor"
	"WeNeedGameServer/packet"
	"log"
	"net"
)

type Receiver struct {
	chanTable         map[uint32]chan packet.PacketI
	connTable         *map[uint32]*net.UDPAddr
	networkActorTable map[uint32]*actor.NetworkActor
	nextSeqTable      *map[uint32]uint32
	nQueue            *internal_type.Queue[packet.PacketI]
	nQueueManager     *QueueManager
	udpConn           *net.UDPConn
}

func NewReceiver(connTable *map[uint32]*net.UDPAddr, nextSeqTable *map[uint32]uint32, nQueue *internal_type.Queue[packet.PacketI], udpConn *net.UDPConn) *Receiver {
	return &Receiver{
		chanTable:         make(map[uint32]chan packet.PacketI),
		connTable:         connTable,
		networkActorTable: make(map[uint32]*actor.NetworkActor),
		nextSeqTable:      nextSeqTable,
		nQueue:            nQueue,
		nQueueManager:     NewQueueManager(nQueue),
		udpConn:           udpConn,
	}
}

func (r *Receiver) StartUDP() {
	go r.nQueueManager.StartQueueManager()
	readBuffer := make([]byte, 2048)
	for {
		readCount, addr, err := r.udpConn.ReadFromUDP(readBuffer)
		if err != nil {
			log.Panicln("잘못된 요청")
		}
		r.handlePacket(readBuffer, readCount, addr)
	}
}

func (r *Receiver) handlePacket(clientPacket []byte, endPoint int, userAddr *net.UDPAddr) {
	data, err := packet.ParsePacketByKind(clientPacket, endPoint)
	if err != nil {
		log.Panicln("잘못된 요청")
	}

	if QPort := (*r.connTable)[data.GetQPort()]; QPort == nil {
		r.tempHandleNewConnection(data.GetQPort(), userAddr)
	}
	r.throwData(data)
}

func (r *Receiver) throwData(data packet.ClientPacketI) {
	if (*r.connTable)[data.GetQPort()] != nil || (*r.nextSeqTable)[data.GetQPort()] == data.GetSEQ() {
		(*r.nextSeqTable)[data.GetQPort()] += 1
		r.chanTable[data.GetQPort()] <- data
	}
}

func (r *Receiver) tempHandleNewConnection(qPort uint32, userAddr *net.UDPAddr) {
	r.chanTable[qPort] = make(chan packet.PacketI)
	(*r.connTable)[qPort] = userAddr
	(*r.nextSeqTable)[qPort] = 1
	networkActor := actor.NewNetworkActor(qPort, userAddr, r.chanTable[qPort], &r.nQueueManager.QmChan)
	r.networkActorTable[qPort] = networkActor
	go r.networkActorTable[qPort].ProcessLoopPacket()
}
