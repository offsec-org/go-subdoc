package ooxml

import "encoding/xml"

type Relationships struct {
	XMLName       xml.Name   `xml:"http://schemas.openxmlformats.org/package/2006/relationships Relationships"`
	Relationships []Relation `xml:"Relationship"`
}

type Relation struct {
	Id         string `xml:",attr"`
	Type       string `xml:",attr"`
	Target     string `xml:",attr"`
	TargetMode string `xml:",attr,omitempty"`
}
