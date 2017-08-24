#!/bin/sh

env GOOS=linux GOARCH=arm go build -o watch-cat-arm
env GOOS=linux GOARCH=amd64 go build -o watch-cat-linux
env GOOS=darwin GOARCH=amd64 go build -o watch-cat-macos
env GOOS=windows GOARCH=amd64 go build -o watch-cat-win
