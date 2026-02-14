#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: ./genhash.sh <password>"
    exit 1
fi
go run tools/genhash.go "$1"
