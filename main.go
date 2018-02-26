package main

import (
	"flag"
	"os"
	"io/ioutil"
	"log"
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

	log.Print("Reading " + extractedFilePath + "...")
	if extractedFilePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

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
		result, compressionErr := Decompress(fileToExtract)
		if compressionErr != nil {
			panic(compressionErr)
		}
		ioutil.WriteFile(savePath, result, os.ModePerm)
	}

	log.Print("Done! Saved to " + savePath)
}
