package smol

import (
	"encoding/binary"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

const MAGICBYTE = 0x10
var extractedFilePath string
var savePath string
var compress bool

// BUG(spotlightishere): Does not handle compression.
func main() {
	flag.StringVar(&extractedFilePath, "extract", "", "Path to file to work with.")
	flag.StringVar(&savePath, "output", "", "Path to save operated upon file")
	flag.BoolVar(&compress, "compress", true, "If true, compresses the file. If false, decompresses.")
	flag.Parse()

	if extractedFilePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Print("Reading " + extractedFilePath + "...")

	fileToExtract, err := ioutil.ReadFile(extractedFilePath)
	if err != nil {
		panic(err)
	}

	if savePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if compress {
		result, compressionErr := Compress(fileToExtract)
		if compressionErr != nil {
			panic(compressionErr)
		}
		ioutil.WriteFile(savePath, result, os.ModePerm)
	} else {
		if fileToExtract[0] != MAGICBYTE {
			log.Panicf("invalid magic byte")
		}

		// There's a u24 or a u32 depending on how the file is.
		// If the next 3 bytes are 0 as an int, we can read one more byte to get the length.
		// Otherwise the first 3 are all we need.
		definedLength := fileToExtract[1:3]
		uncompressedLength := toNDS24(definedLength)
		if uncompressedLength == 0 {
			// That means the total length is a u32 after all.
			// We can go ahead and read it as such.
			uncompressedLength = binary.BigEndian.Uint32(fileToExtract[1:4])
		}

		result, compressionErr := Decompress(fileToExtract, int(uncompressedLength))
		if compressionErr != nil {
			panic(compressionErr)
		}
		ioutil.WriteFile(savePath, result, os.ModePerm)
	}

	log.Print("Done! Saved to " + savePath)
}
