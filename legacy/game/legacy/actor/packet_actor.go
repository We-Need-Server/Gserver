package actor

import (
	"WeNeedGameServer/game/internal/command"
	"WeNeedGameServer/legacy/mediator"
	"WeNeedGameServer/protocol/udp"
	"WeNeedGameServer/protocol/udp/uclient"
	"WeNeedGameServer/util"
	"fmt"
	"math"
	"net"
)

type PacketActor struct {
	QPort      uint32
	UserAddr   *net.UDPAddr
	packetChan chan udp.PacketI
	//actorPlayer *player.Player
	Mediator *mediator.Mediator
}

func NewPacketActor(QPort uint32, UserAddr *net.UDPAddr, packetChan chan udp.PacketI) *PacketActor {
	return &PacketActor{QPort, UserAddr, packetChan, nil}
}

func (a *PacketActor) Register(m *mediator.Mediator) {
	a.Mediator = m
}

func (a *PacketActor) Send(receiverName string, message interface{}) {
	a.Mediator.Notify("actor", receiverName, message)
}

func (a *PacketActor) Receive(senderName string, message interface{}) {
}

func (a *PacketActor) ProcessLoopPacket() {
	for {
		pkt := <-a.packetChan

		switch pkt.GetPacketKind() {
		case 'N':
			a.processEventPacket(pkt.(*uclient.EventPacket))
		}

	}
}

func (a *PacketActor) processEventPacket(packet *uclient.EventPacket) {
	if packet.GetPacketKind() == 'N' {
		a.processCommandPayload(packet.Payload, packet.PayloadEndpoint)
	}

	fmt.Printf("패킷 수신 - 사용자: %s, QPort: %d\n", a.UserAddr, a.QPort)
	fmt.Printf("패킷 내용: %+v\n", packet)
}

func (a *PacketActor) processCommandPayload(payload []byte, payLoadEndpoint int) {
	for i := 0; i < payLoadEndpoint; {
		payloadCommand := command.Command(payload[i : i+2])
		switch payloadCommand {
		case command.FB:
			zDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			fmt.Println("FB", zDelta)
			//a.actorPlayer.MoveForward(zDelta)
			i += 6
			break
		case command.LR:
			xDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			fmt.Println(xDelta)
			//a.actorPlayer.MoveSide(xDelta)
			fmt.Println("LB", xDelta)
			i += 6
			break
		case command.YW:
			yawDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			//a.actorPlayer.TransferYaw(yawDelta)
			fmt.Println("YW", yawDelta)
			i += 6
			break
		case command.PT:
			ptDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			//a.actorPlayer.TransferPT(ptDelta)
			fmt.Println("PT", ptDelta)
			i += 6
			break
		case command.JP:
			jp := util.ByteToBool(payload[i+2])
			//a.actorPlayer.TurnJP(jp)
			fmt.Println("JP", jp)
			i += 3
			break
		case command.SH:
			//a.actorPlayer.TurnIsShoot()
			fmt.Println("SH")
			i += 2
			break
		case command.HT:
			userQPort := util.ConvertBinaryToUint32(payload[i+2 : i+6])
			hpDelta := util.ConvertBinaryToInt16(payload[i+6 : i+8])
			//a.actorPlayer.StoreHitInformation(userQPort, hpDelta)
			fmt.Println("HT", userQPort, hpDelta)
			i += 8
			break
		}
	}
}
