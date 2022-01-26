package ooxml

const (
	// Edited version of xml.Header to match OOXML's standards
	Header = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
)

type Globals struct {
	FileName string
	Target   string
}

var GlobalVar Globals
