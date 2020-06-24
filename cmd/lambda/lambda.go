package main

import (
	"context"
	//"encoding/json"
	//"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djschaap/kv-to-json/pkg/parsedoc"
	"github.com/djschaap/kv-to-json/pkg/sendsqs"
	"os"
	"regexp"
)

func handle_request(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//request_json, _ := json.Marshal(request)
	//fmt.Println("TRACE-received-request:\n", string(request_json))
	var doc string = request.Body
	//fmt.Println("TRACE-received-doc:\n", doc)
	headers, message, _ := parsedoc.ParseDoc(doc)

	queue_url := os.Getenv("DEST_QUEUE")
	has_queue, _ := regexp.MatchString(`^https`, queue_url)
	if has_queue {
		sendsqs.SendMessage(queue_url, headers, message)
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
	sendsqs.OpenSvc()
	lambda.Start(handle_request)
}
