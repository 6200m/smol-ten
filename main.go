package smol_ten

import (
	"flag"
	"os"
	"io/ioutil"
	"log"
)

var extractedFilePath string
var savePath string

// BUG(spotlightishere): Does not handle compression.
func main() {
	flag.StringVar(&extractedFilePath, "extract", "", "Path to file to decompress.")
	flag.StringVar(&savePath, "output", "", "Path to save extracted file")
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
	decrypted, err := Extract(fileToExtract, pubkFile)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(savePath, decrypted, os.ModePerm)
	log.Print("Done! Saved to " + savePath)
}
