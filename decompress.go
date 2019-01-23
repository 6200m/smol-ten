package smol

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"log"
	"os"
)

var decompressed []byte
var compressed *bytes.Buffer

func DecompressDetect(fileToExtract []byte) ([]byte, error) {
	compressionType := fileToExtract[0]

	switch compressionType {
	case TypeLz10:
		// There's a u24 or a u32 depending on how the file is.
		// If the next 3 bytes are 0 as an int, we can read one more byte to get the length.
		// Otherwise the first 3 are all we need.
		var offset int
		definedLength := fileToExtract[1:3]
		offset = 4
		uncompressedLength := toNDS24(definedLength)
		if uncompressedLength == 0 {
			// That means the total length is a u32 after all.
			// We can go ahead and read it as such.
			uncompressedLength = binary.BigEndian.Uint32(fileToExtract[1:4])
			offset = 5
		}

		result, compressionErr := DecompressLZ10(fileToExtract[offset:], int(uncompressedLength))
		if compressionErr != nil {
			return nil, compressionErr
		}

		return result, nil
	case TypeLz11:
		// TODO: For the future.
		return DecompressLZ11(fileToExtract, 0)
	default:
		log.Fatalf("Type 0x%x is not implemented.", compressionType)
		os.Exit(1)
		return nil, nil
	}

	return nil, nil
}

// Decompress takes input, reads the length and attempts ts best to decompress under LZ10.
// See also: https://github.com/profi200/splashtool/blob/7bb5342839866bfecb16537cc8caf1d966dcd545/tex3ds/source/lzss.c
func DecompressLZ10(input []byte, decompressedLength int) ([]byte, error) {
	compressed = bytes.NewBuffer(input)

	var decompressed []byte
	println("Expected length:", decompressedLength)

	var flags byte = 0x00
	var mask byte = 0x01

	for decompressedLength > 0 {
		//log.Printf("0x%x", byteToHandle)

		// If there's no mask, reset.
		if mask == 0x01 {
			flags = nextByte()
			log.Printf("flag: 0x%x", flags)
			mask = 0x80
		} else {
			mask = mask >> 1
		}

		if (flags & mask) > 0x00 {
			// Seems like we've a compressed block
			firstByte := uint8(nextByte())
			//length := firstByte & 0xF0
			length := firstByte >> 4
			length += 3

			disp := uint8(firstByte) & 0x0F
			disp = disp<<8 | uint8(nextByte())
			disp += 0x01

			if int(length) > decompressedLength {
				length = uint8(decompressedLength)
			}

			decompressedLength -= int(length)

			log.Printf("%d from disp 0x%x (current length %d)", length, disp, len(decompressed))
			for length > 0 {
				// Get location to copy byte from
				offset := len(decompressed) - int(int8(disp)) - 1
				log.Printf("Reading from 0x%x...", offset)

				// Go back for it and add to end
				decompressed = append(decompressed, decompressed[offset])
				length--
			}
		} else {
			// Not compressed. 1:1
			decompressed = append(decompressed, nextByte())
			decompressedLength--
		}
	}

	return decompressed, nil
}

func DecompressLZ11(input []byte, decompressedLength int) ([]byte, error) {
	return nil, errors.New("lz11 is unimplemented")
}

// nextByte returns a single byte out of the buffer array given.
func nextByte() byte {
	return compressed.Next(1)[0]
}
