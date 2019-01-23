package smol

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

const (
	TypeLz10 = 0x10
	TypeLz11 = 0x11
)

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
		result, compressionErr := CompressLZ10(fileToExtract)
		if compressionErr != nil {
			panic(compressionErr)
		}
		ioutil.WriteFile(savePath, result, os.ModePerm)
	} else {
		result, err := DecompressDetect(decompressed)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(savePath, result, os.ModePerm)
	}

	log.Print("Done! Saved to " + savePath)
}
