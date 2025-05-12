package packet

import (
	"WeNeedGameServer/packet/client"
	"fmt"
)

type PacketI interface {
	GetPacketKind() uint8
	GetQPort() uint32
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

func ParsePacketByKind(np []byte, endPoint int) (PacketI, error) {
	if len(np) < 1 {
		return nil, fmt.Errorf("packet too small: %d bytes", len(np))
	}

	pKind := np[0]

	fmt.Println("pkind", pKind)

	switch pKind {
	case 41:
		return client.ParseEventPacket(np, endPoint), nil
	case 46:
		return client.ParseTickIPacket(np, endPoint), nil
	case 50:
		return client.ParseTickRPacket(np, endPoint), nil
	default:
		return nil, fmt.Errorf("unknown packet kind: %d", pKind)
	}
}
