package packet

import (
	"WeNeedGameServer/util"
	"fmt"
)

type PacketI interface {
	GetPacketKind() uint32
	GetQPort() uint32
}

func ParsePacketByKind(np []byte, endPoint int) (PacketI, error) {
	if len(np) < 4 {
		return nil, fmt.Errorf("packet too small: %d bytes", len(np))
	}

	PKind := util.ConvertBinaryToUint32(np[0:4])

	switch PKind {
	case 41:
		return (&TickIPacket{}).ParsePacket(np, endPoint), nil
	case 46:
		return (&EventPacket{}).ParsePacket(np, endPoint), nil
	case 50:
		return (&TickRPacket{}).ParsePacket(np, endPoint), nil
	default:
		return nil, fmt.Errorf("unknown packet kind: %d", PKind)
	}
}
