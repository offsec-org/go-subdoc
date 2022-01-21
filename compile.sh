#!/bin/bash

# https://go.dev/doc/install/source#environment

read -p 'Target OS (windows, linux, ...): ' TARGET_OS
GOOS=$TARGET_OS GOARCH=amd64 go build -ldflags "-s -w" biscoito/go-subdoc