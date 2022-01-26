package ooxml

import "encoding/xml"

type Relationships struct {
	XMLName       xml.Name   `xml:"http://schemas.openxmlformats.org/package/2006/relationships Relationships"`
	Relationships []Relation `xml:"Relationship"`
}

type Relation struct {
	Id         string `xml:"Id,attr"`
	Type       string `xml:"Type,attr"`
	Target     string `xml:"Target,attr"`
	TargetMode string `xml:"TargetMode,attr,omitempty"`
}
