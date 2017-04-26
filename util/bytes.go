package util

import (
	"encoding/binary"
)

func BytesToUInt32(buf []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(buf))
}

func BytesToString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func BytesTrim(c []byte) []byte {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return c[:n+1]
}
