package util

import (
	"encoding/binary"
)

func BytesToUInt32(buf []byte) uint32 {
	return uint32(binary.LittleEndian.Uint32(buf))
}
