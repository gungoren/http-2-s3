#!/bin/sh

set -euo pipefail

output_file="http-2-s3"

if [[ -e $output_file ]]; then
  echo "rebuilding: $output_file"
else
  echo "$output_file does not exist"
fi

GOOS=linux GOARCH=amd64 go build -v -o $output_file ./

zip $output_file.zip $output_file
