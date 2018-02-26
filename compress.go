package main

func Compress(input []byte) (compressedOutput []byte, err error) {
	// We need to go ahead and set up the header.
	// We start with the file's magic.
	var compressed []byte
	var header [4]byte
	header[0] = MAGICBYTE

	// The next u24 is the length.
	// We split the length in 3 bytes.
	length := len(input)
	header[1] = byte(length) & 0xFF
	header[2] = byte(length >> 8) & 0xFF
	header[3] = byte(length >> 16) & 0xFF

	// Now to add the header.
	compressed = append(compressed, header[:]...)

	// The fun part: compressing!
	// We'll start by chunking everything into 8 bytes.
	// A buffer, assuming every byte is different and cannot be compressed:

	return compressed, nil
}