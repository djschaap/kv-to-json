package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djschaap/kv-to-json/pkg/parsedoc"
	"github.com/djschaap/kv-to-json/pkg/sendsns"
	"os"
	"regexp"
)

func handle_request(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	enable_trace := os.Getenv("ENABLE_TRACE")
	var doc string = request.Body
	if enable_trace == "1" {
		request_json, _ := json.Marshal(request)
		fmt.Println("TRACE-received-request:\n", string(request_json))
		//fmt.Println("TRACE-received-doc:\n", doc)
	}
	headers, message, _ := parsedoc.ParseDoc(doc)

	topic_arn := os.Getenv("TOPIC_ARN")
	has_queue, _ := regexp.MatchString(`^arn:`, topic_arn)
	if has_queue {
		sendsns.SendMessage(topic_arn, headers, message)
	} else {
		headers_json_bytes, _ := json.Marshal(headers)
		fmt.Println("CANNOT SEND headers:\n", string(headers_json_bytes))
		fmt.Println("CANNOT SEND message:\n", message)
	}
	response := events.APIGatewayProxyResponse{
		IsBase64Encoded: false,
		StatusCode:      200,
		//Headers: {
		//	"x-customer-header": "value"
		//},
		Body: "ok",
	}
	return response, nil
}

func main() {
	sendsns.OpenSvc()
	lambda.Start(handle_request)
}
