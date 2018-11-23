package smol

import (
	"bytes"
	"errors"
)

var decompressed []byte
var compressed *bytes.Buffer

// Decompress takes input, reads the length and attempts
// its best to decompress under LZ10/11/77/all those fun things.
// See also: https://github.com/Barubary/dsdecmp/blob/master/CSharp/DSDecmp/Formats/Nitro/LZ11.cs#L85-L118
func Decompress(input []byte, decompressedLength int) ([]byte, error) {
	compressed = bytes.NewBuffer(input)
	println(decompressedLength)

	for decompressedLength > len(decompressed) {
		byteToHandle := nextByte()

		// If the "flag" byte is 00, we can just do a shortcut
		// and write out the next 8 bytes.
		//if byteToHandle == 0x00 {
		//	println("Skipping byte")
		//	addToDecompressed(compressed.Next(8))
		//} else if byteToHandle == 0x01 {
			// We'll go bit by bit to determine the next action.
			for _, flagBit := range getBits(byteToHandle) {
				//println("Currently at bit", loc)
				//fmt.Printf("Hi from 0x%x\n", byteToHandle)

				if flagBit == 0 {
					// We just write this referenced byte, no compression applied.
					nextByte := compressed.Next(1)
					//fmt.Printf("0x%d\n", nextByte)
					addToDecompressed(nextByte)
				} else if flagBit == 1 {
					workingByte := nextByte()
					indicator := workingByte >> 4
					var count int

					// TODO: better document
					if indicator == 0 {
						count = int(workingByte << 4)

						workingByte = nextByte()
						count += int(workingByte >> 4)
						count += 0x11
					} else if indicator == 1 {
						count = int((workingByte & 0xf) << 12)
						workingByte = nextByte()
						count += int(workingByte >> 4)
						count += 0x111
					} else {
						count = int(workingByte >> 4)
						count += 1
					}

					disp := ((workingByte & 0xf) << 8) + nextByte()
					disp += 1

					for count != 0 {
						//fmt.Println("Organically sourcing from", int(disp))
						//println(hex.EncodeToString(decompressed))
						addToDecompressed(decompressed[len(decompressed)-int(disp):])
						count--
					}
				} else {
					return nil, errors.New("unknown flag bit")
				}
			}
		//} else {
		//	// The block flag should not be anything other than 0 or 1.
		//	return nil, errors.New("invalid byte flag received, found " + fmt.Sprintf("0x%x", byteToHandle))
		//}
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