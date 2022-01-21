package ooxml

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"log"

	"github.com/go-xmlfmt/xmlfmt"
)

func ReadPackage(fileName string) error {
	log.Println("Opening document...")

	// Open a zip archive for reading.
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	// Iterate through the files in the archive.
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		xml_data, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		rc.Close()

		err = WritePackage(f, xml_data)
		if err != nil {
			return err
		}
	}

	return nil
}

func WritePackage(f *zip.File, xml_data []byte) error {
	file_filter_list := map[string]bool{
		"word/_rels/document.xml.rels": true,
		"word/document.xml":            true,
		"word/styles.xml":              true,
	}

	xml_formatted := xmlfmt.FormatXML(string(xml_data[:]), "", "	", true)

	if file_filter_list[f.Name] {
		switch f.Name {
		case c_document_rels:
			log.Printf("\nContents of %s: %s\n", f.Name, xml_formatted)

			// xml_structure := &Relationships{
			// 	Relationships: []Relation{
			// 		Id:         "rId5",
			// 		Type:       "http://schemas.openxmlformats.org/officeDocument/2006/relationships/subDocument",
			// 		Target:     `file:///\\biscoito.eu\test\`,
			// 		TargetMode: "External",
			// 	},
			// }

			var rl Relationships

			xml.Unmarshal(xml_data, &rl)

		case c_document:
		case c_styles:
		default:
		}
	}

	return nil
}
