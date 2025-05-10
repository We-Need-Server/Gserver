package packet

import (
	"WeNeedGameServer/packet/eventpacket"
	"WeNeedGameServer/packet/tickipacket"
	"WeNeedGameServer/packet/tickrpacket"
	"WeNeedGameServer/util"
	"fmt"
)

type PacketI interface {
	GetPacketKind() uint32
	GetQPort() uint32
}

type FieldType int

const (
	TypeUint8 FieldType = iota
	TypeUint16
	TypeUint32
	TypeString
	TypeFloat32
	TypeFloat64
	TypeBytes
)

type Field struct {
	PropertyName string
	Offset       uint32
	PropertyType FieldType
	// 타입에 대해서
}

type PropertyMap map[string]Field

func ParsePacketByKind(np []byte, endPoint int) (PacketI, error) {
	if len(np) < 4 {
		return nil, fmt.Errorf("packet too small: %d bytes", len(np))
	}

	PKind := util.ConvertBinaryToUint32(np[0:4])

	switch PKind {
	case 41:
		return eventpacket.ParseEventPacket(np, endPoint), nil
	case 46:
		return tickipacket.ParseTickIPacket(np, endPoint), nil
	case 50:
		return tickrpacket.ParseTickRPacket(np, endPoint), nil
	default:
		return nil, fmt.Errorf("unknown packet kind: %d", PKind)
	}
}
