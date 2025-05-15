package util

import (
	"encoding/binary"
	"runtime"
)

func ConvertBinaryToUint32(b []byte) uint32 {
	if runtime.GOARCH == "arm64" {
		return binary.BigEndian.Uint32(b)
	} else {
		return binary.LittleEndian.Uint32(b)
	}
}

func ConvertBinaryToInt16(b []byte) int16 {
	if runtime.GOARCH == "arm64" {
		return int16(binary.BigEndian.Uint16(b))
	} else {
		return int16(binary.LittleEndian.Uint16(b))
	}
}

func ByteToBool(b byte) bool {
	return b != 0
}

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
