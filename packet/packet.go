package packet

import (
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

type PacketField struct {
	PropertyName string
	Offset       uint32
	PropertyType FieldType
	// 타입에 대해서
}

func ParsePacketByKind(np []byte, endPoint int) (PacketI, error) {
	if len(np) < 4 {
		return nil, fmt.Errorf("packet too small: %d bytes", len(np))
	}

	PKind := util.ConvertBinaryToUint32(np[0:4])

	switch PKind {
	case 41:
		return ParseEventPacket(np, endPoint), nil
	case 46:
		return ParseTickIPacket(np, endPoint), nil
	case 50:
		return ParseTickRPacket(np, endPoint), nil
	default:
		return nil, fmt.Errorf("unknown packet kind: %d", PKind)
	}
}
