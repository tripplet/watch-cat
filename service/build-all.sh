#!/bin/sh

mkdir bin
env GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o bin/watchcat-arm
env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/watchcat-linux
env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/watchcat-macos
env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/watchcat-win
