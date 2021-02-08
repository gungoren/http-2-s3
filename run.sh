#!/bin/sh

set -euo pipefail

output_file="http-2-s3"

export S3_BUCKET_NAME=<bucket-name>
export AWS_ACCESS_KEY_ID=<access-key>
export AWS_SECRET_ACCESS_KEY=<secret-key>

./$output_file

