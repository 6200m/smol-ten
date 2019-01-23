package smol

import "encoding/binary"

// toNDS24 takes a buffer and an offset. It takes 3 bytes (u24) and converts to uint32.
func toNDS24(buffer []byte) uint32 {
	// Add a null in to make it 4 bytes.
	return binary.BigEndian.Uint32(append([]byte{0x00}, buffer[0:3]...))
}
