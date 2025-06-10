package actor

import (
	"WeNeedGameServer/command"
	"WeNeedGameServer/game_type"
	"WeNeedGameServer/protocol/udp"
	"WeNeedGameServer/protocol/udp/uclient"
	"WeNeedGameServer/protocol/udp/userver"
	"WeNeedGameServer/util"
	"math"
	"net"
)

type UdpActor struct {
	qPort      uint32
	userAddr   *net.UDPAddr
	packetChan chan udp.PacketI
	qmChan     chan udp.PacketI
}

func NewUdpActor(qPort uint32, userAddr *net.UDPAddr, packetChan chan udp.PacketI, qmChan chan udp.PacketI) *UdpActor {
	return &UdpActor{qPort, userAddr, packetChan, qmChan}
}

func (na *UdpActor) ProcessLoopPacket() {
	for {
		pkt := <-na.packetChan

		switch pkt.GetPacketKind() {
		case 'N':
			na.processEventPacket(pkt.(*uclient.EventPacket))
		}

	}
}

func (na *UdpActor) processEventPacket(packet *uclient.EventPacket) {
	if packet.GetPacketKind() == 'N' {
		na.processCommandPayload(packet.Payload, packet.PayloadEndpoint)
	}
}

func (na *UdpActor) processCommandPayload(payload []byte, payLoadEndpoint int) {
	playerPosition := game_type.NewPlayerStateDefault()
	hitInformationMap := make(map[uint32]int16)
	for i := 0; i < payLoadEndpoint; {
		payloadCommand := command.Command(payload[i : i+2])
		switch payloadCommand {
		case command.FB:
			zDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			playerPosition.PositionZ += zDelta
			i += 6
			break
		case command.LR:
			xDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			playerPosition.PositionX += xDelta
			i += 6
			break
		case command.YW:
			yawDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			playerPosition.YawAngle += yawDelta
			i += 6
			break
		case command.PT:
			ptDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			playerPosition.PtAngle += ptDelta
			i += 6
			break
		case command.JP:
			playerPosition.Jp = true
			i += 2
			break
		case command.SH:
			playerPosition.IsShoot = true
			i += 2
			break
		case command.HT:
			userQPort := util.ConvertBinaryToUint32(payload[i+2 : i+6])
			hpDelta := util.ConvertBinaryToInt16(payload[i+6 : i+8])
			hitInformationMap[userQPort] += hpDelta
			i += 8
			break
		case command.RE:
			playerPosition.IsReload = true
			i += 2
			break
		}
	}

	na.qmChan <- userver.NewDeltaPacket(na.qPort, playerPosition, hitInformationMap)
}
