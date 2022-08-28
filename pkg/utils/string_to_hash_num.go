package utils

import "encoding/binary"

func StrToHashInt(originStr string) (uint16, uint16) {
	if len(originStr) < 1 {
		return 0, 0
	}
	byteSeque := []byte(originStr)
	if len(byteSeque) < 2 {
		return 0, uint16(uint8(byteSeque[0]))
	}
	if len(byteSeque) < 3 {
		return 0, binary.LittleEndian.Uint16(byteSeque[0:])
	}
	if len(byteSeque) < 4 {
		return uint16(uint8(byteSeque[0])), binary.LittleEndian.Uint16(byteSeque[1:])
	}

	high := binary.LittleEndian.Uint16(byteSeque[0:])
	low := binary.LittleEndian.Uint16(byteSeque[2:])
	return high, low
	// return uint32(high)<<16 + uint32(low)
}
