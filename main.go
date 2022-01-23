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
	var fileName string

	flag.StringVar(&fileName, "input", "", "Target document")
	flag.Parse()

	// If argument invalid / not supplied
	if len(fileName) <= 0 {
		fmt.Printf("Usage: %s -input target.doc\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	filePath, err := filepath.Abs(fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = ooxml.ReadPackage(filePath)
	if err != nil {
		log.Fatal(err)
	}
}
