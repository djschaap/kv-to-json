# kv-to-json

[![Build Status - master](https://travis-ci.com/djschaap/kv-to-json.svg?branch=master)](https://travis-ci.com/djschaap/kv-to-json)

## Overview

kv-to-json accepts a document (typically submitted to AWS API Gateway
via HTTP POST), converts to JSON, and submits to AWS SQS.

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
go get github.com/aws/aws-lambda-go/lambda \
  && GOOS=linux GOARCH=amd64 go build cmd/lambda/lambda.go \
  && go test all \
  && zip lambda.zip lambda \
  && aws lambda update-function-code --function-name kv-to-json \
    --zip-file fileb://lambda.zip
```

## Run Individual Test (disable cache)

```bash
go test ./pkg/parsedoc -count=1 -v
```
