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

func ByteToBool(b byte) bool {
	return b != 0
}

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
