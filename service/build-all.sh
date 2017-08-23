#!/bin/sh

env GOOS=linux GOARCH=arm go build -o watch-cat-rpi
env GOOS=linux GOARCH=amd64 go build -o watch-cat-amd64
env GOOS=darwin GOARCH=amd64 go build -o watch-cat-macos
