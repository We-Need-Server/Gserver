package actor

import (
	"WeNeedGameServer/command"
	"WeNeedGameServer/game/player"
	"WeNeedGameServer/packet"
	"WeNeedGameServer/packet/client"
	"WeNeedGameServer/packet/server"
	"WeNeedGameServer/util"
	"fmt"
	"math"
	"net"
)

type NetworkActor struct {
	qPort      uint32
	userAddr   *net.UDPAddr
	packetChan chan packet.PacketI
	qmChan     *chan packet.PacketI
}

func NewNetworkActor(qPort uint32, userAddr *net.UDPAddr, packetChan chan packet.PacketI, qmChan *chan packet.PacketI) *NetworkActor {
	return &NetworkActor{qPort, userAddr, packetChan, qmChan}
}

func (na *NetworkActor) ProcessLoopPacket() {
	for {
		pkt := <-na.packetChan

		switch pkt.GetPacketKind() {
		case 'N':
			na.processEventPacket(pkt.(*client.EventPacket))
		}

	}
}

func (na *NetworkActor) processEventPacket(packet *client.EventPacket) {
	if packet.GetPacketKind() == 'N' {
		na.processCommandPayload(packet.Payload, packet.PayloadEndpoint)
	}

	fmt.Printf("패킷 수신 - 사용자: %s, QPort: %d\n", na.userAddr, na.qPort)
	fmt.Printf("패킷 내용: %+v\n", packet)
}

func (na *NetworkActor) processCommandPayload(payload []byte, payLoadEndpoint int) {
	playerPosition := player.NewPlayerPositionD()
	hitInformationMap := make(map[uint32]int16)
	for i := 0; i < payLoadEndpoint; {
		payloadCommand := command.Command(payload[i : i+2])
		switch payloadCommand {
		case command.FB:
			zDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			fmt.Println("FB", zDelta)
			playerPosition.PositionZ += zDelta
			//a.actorPlayer.MoveForward(zDelta)
			i += 6
			break
		case command.LR:
			xDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			fmt.Println(xDelta)
			playerPosition.PositionX += xDelta
			//a.actorPlayer.MoveSide(xDelta)
			fmt.Println("LB", xDelta)
			i += 6
			break
		case command.YW:
			yawDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			//a.actorPlayer.TransferYaw(yawDelta)
			fmt.Println("YW", yawDelta)
			playerPosition.YawAngle += yawDelta
			i += 6
			break
		case command.PT:
			ptDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			//a.actorPlayer.TransferPT(ptDelta)
			fmt.Println("PT", ptDelta)
			playerPosition.PtAngle += ptDelta
			i += 6
			break
		case command.JP:
			//jp := util.ByteToBool(payload[i+2])
			//a.actorPlayer.TurnJP(jp)
			fmt.Println("JP")
			playerPosition.Jp = true
			i += 2
			break
		case command.SH:
			//a.actorPlayer.TurnIsShoot()
			fmt.Println("SH")
			playerPosition.IsShoot = true
			i += 2
			break
		case command.HT:
			userQPort := util.ConvertBinaryToUint32(payload[i+2 : i+6])
			hpDelta := util.ConvertBinaryToInt16(payload[i+6 : i+8])
			//a.actorPlayer.StoreHitInformation(userQPort, hpDelta)
			fmt.Println("HT", userQPort, hpDelta)
			hitInformationMap[userQPort] += hpDelta
			i += 8
			break
		case command.RE:
			playerPosition.IsReload = true
			i += 2
			break
		}
	}

	*na.qmChan <- server.NewDeltaPacket(na.qPort, playerPosition, &hitInformationMap)
}
