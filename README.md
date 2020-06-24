# kv-to-json

[![Build Status - master](https://travis-ci.com/djschaap/kv-to-json.svg?branch=master)](https://travis-ci.com/djschaap/kv-to-json)

## Overview

kv-to-json accepts a document (typically submitted to AWS API Gateway
via HTTP POST), converts to JSON, and submits to AWS SNS.

## Input Format

The input document is expected to follow a simple key=value format.
Each line consists of a key (alphanumeric characters), followed by
a colon, followed by the value.
Each key-value pair is separated by a newline.
Spaces are not allowed in keys.
The first colon encountered ends the key.
Leading whitespace in the value will be dropped.
Trailing whitespace will be maintained (for now - subject to change).
Invalid lines will be quietly ignored (dropped).

The document is in two sections, header key-values and message
key-values.
The two sections are separated by a blank line.

Example doc:
```
X-secret: abcdef

key1: value1
key2: value number two
```

## Build ZIP for AWS Lambda

```bash
BUILD_DT=`date +%FT%T%z`
COMMIT=`git rev-parse --short HEAD`
FULL_COMMIT=`git log -1`
VER=0.0.0
go get github.com/aws/aws-lambda-go/lambda \
  && GOOS=linux GOARCH=amd64 go build -ldflags \
    "-X main.build_dt=${BUILD_DT} -X main.commit=${COMMIT} -X main.version=${VER}" \
    cmd/lambda/lambda.go \
  && go test all \
  && zip kv-to-json-${VER}.zip lambda
```

### Upload Lambda Code to S3

```bash
aws s3 cp kv-to-json-${VER}.zip s3://PACKAGE_BUCKET_NAME
```

### Manually Update Lambda Code (dev/test)

```bash
aws lambda update-function-code --function-name kv-to-json \
  --zip-file fileb://kv-to-json-${VER}.zip
```

## Run Individual Test (disable cache)

```bash
go test ./pkg/parsedoc -count=1 -v
```
