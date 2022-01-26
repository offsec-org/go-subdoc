package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"biscoito/go-subdoc/ooxml"
)

func main() {
	// Parse cmd arguments
	flag.StringVar(&ooxml.GlobalVar.FileName, "input", "", "Target document")
	flag.StringVar(&ooxml.GlobalVar.Target, "target", "", "Target server (only domain / ip address)")
	flag.Parse()

	// If argument invalid / not supplied
	if len(ooxml.GlobalVar.FileName) <= 0 || len(ooxml.GlobalVar.Target) <= 0 {
		fmt.Printf("Usage: %s -input target.docx/docm -target example.com/127.0.0.1\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Get absolute file path of the target file
	filePath, err := filepath.Abs(ooxml.GlobalVar.FileName)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize routine
	err = ooxml.Initialize(filePath)
	if err != nil {
		log.Fatal(err)
	}
}
