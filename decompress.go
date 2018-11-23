package smol

import (
	"bytes"
	"log"
)

var decompressed []byte
var compressed *bytes.Buffer

// Decompress takes input, reads the length and attempts ts best to decompress under LZ10.
// See also: https://github.com/profi200/splashtool/blob/7bb5342839866bfecb16537cc8caf1d966dcd545/tex3ds/source/lzss.c
func Decompress(input []byte, decompressedLength int) ([]byte, error) {
	compressed = bytes.NewBuffer(input)
	var decompressed []byte
	println("Expected length:", decompressedLength)

	var flags byte = 0x00
	var mask byte = 0x00

	for decompressedLength > 0 {
		byteToHandle := nextByte()

		// If there's no mask, reset.
		if mask == 0x00 {
			flags = byteToHandle
			mask = 0x80
		}

		if (flags & mask) != 0x00 {
			// Seems like we've a compressed block
			length := uint8(byteToHandle) & 0xF0
			length = length >> 4
			length += 3

			disp := uint8(nextByte()) & 0x0F
			disp = disp << 8 | uint8(nextByte())

			if int(length) > decompressedLength {
				length = uint8(decompressedLength)
			}

			decompressedLength -= int(length)

			for length > 0 {
				// Get location to copy byte from
				offset := len(decompressed) - int(disp) - 1

				log.Println(len(decompressed), "->", offset)

				// Go back for it and add to end
				decompressed = append(decompressed, decompressed[offset])
			}
		} else {
			// Not compressed. 1:1
			decompressed = append(decompressed, nextByte())
			decompressedLength--
		}
		mask >>= 1
	}

	return decompressed, nil
}

// getBits returns an array of the 8 bits in a byte.
func getBits(source byte) [8]byte {
	return [8]byte{
		(source >> 7) & 1,
		(source >> 6) & 1,
		(source >> 5) & 1,
		(source >> 4) & 1,
		(source >> 3) & 1,
		(source >> 2) & 1,
		source & 1,
	}
}

// nextByte returns a single byte out of the buffer array given.
func nextByte() byte {
	return compressed.Next(1)[0]
}

// addToDecompressed modifies the decompressed array
func addToDecompressed(given []byte) {
	decompressed = append(decompressed, given...)
}