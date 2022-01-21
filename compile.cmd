@echo off
:: https://go.dev/doc/install/source#environment

set /p TARGET_OS=Target OS (windows, linux, ...): 

set GOOS=%TARGET_OS%
set GOARCH=amd64

go build -ldflags "-s -w" biscoito/go-subdoc

pause
exit /b 0