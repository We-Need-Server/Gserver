package actor

import (
	"WeNeedGameServer/command"
	"WeNeedGameServer/game/player"
	"WeNeedGameServer/packet"
	"WeNeedGameServer/util"
	"fmt"
	"math"
	"net"
)

type PacketActor struct {
	NextSEQ     uint32
	QPort       uint32
	UserAddr    *net.UDPAddr
	packetChan  chan *packet.Packet
	actorPlayer *player.Player
}

func NewPacketActor(NextSEQ uint32, QPort uint32, UserAddr *net.UDPAddr, packetChan chan *packet.Packet, actorPlayer *player.Player) *PacketActor {
	return &PacketActor{NextSEQ, QPort, UserAddr, packetChan, actorPlayer}
}

func (a *PacketActor) ProcessLoopPacket() {
	for {
		pkt := <-a.packetChan
		fmt.Println("시퀀스 번호 전", a.NextSEQ, pkt.SEQ)
		if pkt.SEQ == a.NextSEQ {
			a.NextSEQ += uint32(1)
			a.processCommandPayload(pkt.Payload, pkt.PayloadEndpoint)
		}

		// 패킷 정보 출력
		fmt.Println("시퀀스 번호 후", a.NextSEQ, pkt.SEQ)
		fmt.Printf("패킷 수신 - 사용자: %s, QPort: %d\n", a.UserAddr, a.QPort)
		fmt.Printf("패킷 내용: %+v\n", pkt)
	}
}

func (a *PacketActor) processCommandPayload(payload []byte, payLoadEndpoint int) {
	for i := 0; i < payLoadEndpoint; {
		payloadCommand := string(payload[i : i+2])
		switch payloadCommand {
		case command.FB:
			xDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			fmt.Println(xDelta)
			a.actorPlayer.MoveFoward(xDelta)
			i += 6
			break
		case command.RL:
			zDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			fmt.Println(zDelta)
			a.actorPlayer.MoveSide(zDelta)
			i += 6
			break
		case command.YW:
			yawDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			fmt.Println(yawDelta)
			a.actorPlayer.TransferYaw(yawDelta)
			i += 6
			break
		case command.PT:
			ptDelta := math.Float32frombits(util.ConvertBinaryToUint32(payload[i+2 : i+6]))
			fmt.Println(ptDelta)
			a.actorPlayer.TransferPT(ptDelta)
			i += 6
			break
		}
	}
}
