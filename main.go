package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"biscoito/go-subdoc/ooxml"
)

type Globals struct {
	FileName string
}

func main() {
	var vars Globals

	flag.StringVar(&vars.FileName, "input", "", "Target document")
	flag.Parse()

	// If argument invalid / not supplied
	if len(vars.FileName) <= 0 {
		fmt.Printf("Usage: %s -input target.doc\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	file_path, err := filepath.Abs(vars.FileName)
	if err != nil {
		log.Fatal(err)
	}

	err = ooxml.ReadPackage(file_path)
	if err != nil {
		log.Fatal(err)
	}
}
