package packet

import (
	"WeNeedGameServer/packet/client"
	"WeNeedGameServer/util"
	"fmt"
	"math"
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

func ParsePacketByKind(np []byte, endPoint int) (uint32, PacketI, error) {
	if len(np) < 4 {
		return uint32(math.MaxUint32), nil, fmt.Errorf("packet too small: %d bytes", len(np))
	}

	pKind := util.ConvertBinaryToUint32(np[0:4])

	switch pKind {
	case 41:
		return pKind, client.ParseEventPacket(np, endPoint), nil
	case 46:
		return pKind, client.ParseTickIPacket(np, endPoint), nil
	case 50:
		return pKind, client.ParseTickRPacket(np, endPoint), nil
	default:
		return uint32(math.MaxUint32), nil, fmt.Errorf("unknown packet kind: %d", pKind)
	}
}
