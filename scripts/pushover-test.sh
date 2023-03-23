#!/bin/sh

go build -o pushover ./cmd/...
./pushover --help
