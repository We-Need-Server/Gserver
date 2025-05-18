package receiver

import (
	"WeNeedGameServer/game/actor"
	"WeNeedGameServer/internal_type"
	"WeNeedGameServer/packet"
	"log"
	"net"
)

type Receiver struct {
	chanTable        map[uint32]chan packet.PacketI
	connTable        *map[uint32]*net.UDPAddr
	packetActorTable map[uint32]*actor.PacketActor
	nextSEQTable     map[uint32]uint32
	nQueue           *internal_type.Queue[packet.PacketI]
	nQueueManager    *QueueManager
	udpConn          *net.UDPConn
}

func NewReceiver(connTable *map[uint32]*net.UDPAddr, nQueue *internal_type.Queue[packet.PacketI], udpConn *net.UDPConn) *Receiver {
	return &Receiver{make(map[uint32]chan packet.PacketI), connTable, make(map[uint32]*actor.PacketActor), make(map[uint32]uint32), nQueue, NewQueueManager(nQueue), udpConn}
}

func (r *Receiver) StartUDP() {
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
	r.throwData(data, userAddr)
}

func (r *Receiver) throwData(data packet.PacketI, userAddr *net.UDPAddr) {
	if (*r.connTable)[data.GetQPort()] != nil || r.nextSEQTable[data.GetQPort()] == data.GetSEQ() {
		r.nextSEQTable[data.GetQPort()] += 1
		r.chanTable[data.GetQPort()] <- data
	}
}

func (r *Receiver) tempHandleNewConnection(qPort uint32, userAddr *net.UDPAddr) {
	r.chanTable[qPort] = make(chan packet.PacketI)
	(*r.connTable)[qPort] = userAddr
	r.nextSEQTable[qPort] = 1
	packetActor := actor.NewPacketActor(qPort, userAddr, r.chanTable[qPort])
	r.packetActorTable[qPort] = packetActor
	go r.packetActorTable[qPort].ProcessLoopPacket()
}
