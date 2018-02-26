package main

import (
	"errors"
	"encoding/binary"
	"log"
	"strconv"
)

// Decompress takes input and verifies against a magic byte.
// If the check is successful, it then reads the length and attempts
// its best to decompress.
func Decompress(input []byte) (decompress []byte, err error) {
	if input[0] != MAGICBYTE {
		return nil, errors.New("invalid magic byte")
	}

	// Next, there's a u24 or a u32 depending on how the file is.
	// If the next 3 bytes are 0 as an int, we can read one more byte to get the length.
	// Otherwise the first 3 are all we need.
	definedLength := input[1:3]
	totalLength := toNDS24(definedLength)
	if totalLength == 0 {
		// That means the total length is a u32 after all.
		// We can go ahead and read it as such.
		totalLength = binary.BigEndian.Uint32(input[1:4])
	}
	log.Print("your binary is ", strconv.Itoa(int(totalLength)), " big")

	return nil, nil
}