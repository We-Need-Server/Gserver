package udp

import (
	udp_client2 "WeNeedGameServer/protocol/udp/uclient"
	"fmt"
)

type PacketI interface {
	GetPacketKind() uint8
	GetQPort() uint32
}

type ClientPacketI interface {
	PacketI
	GetSEQ() uint32
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

func ParsePacketByKind(np []byte, endPoint int) (ClientPacketI, error) {
	if len(np) < 1 {
		return nil, fmt.Errorf("packet too small: %d bytes", len(np))
	}

	pKind := np[0]

	fmt.Println("pkind", pKind)

	switch pKind {
	case 'N':
		return udp_client2.ParseEventPacket(np, endPoint), nil
	case 'I':
		return udp_client2.ParseTickIPacket(np, endPoint), nil
	case 'R':
		return udp_client2.ParseTickRPacket(np, endPoint), nil
	default:
		return nil, fmt.Errorf("unknown packet kind: %d", pKind)
	}
}
