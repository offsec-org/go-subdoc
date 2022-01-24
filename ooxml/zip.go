package ooxml

const (
	CDocumentXMLRels string = "word/_rels/document.xml.rels"
	CDocument        string = "word/document.xml"
	CStyles          string = "word/styles.xml"
)

type ZipFile struct {
	Name     string
	Contents []byte
}

type ZipFiles struct {
	Files []ZipFile
}

var ZipArray ZipFiles
