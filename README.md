# go-subdoc
A tool made in Go that injects a hidden malicious subdoc in a Word document.

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

Please input only either a domain or an IP Address as the target.
