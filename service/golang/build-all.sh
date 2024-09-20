#!/bin/sh


mkdir -p bin
env GOOS=linux GOARCH=arm go build -ldflags "-s -w" -trimpath -o bin/watchcat-arm
env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o bin/watchcat-linux
env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o bin/watchcat-macos
env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o bin/watchcat-win
