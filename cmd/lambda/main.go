package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djschaap/kv-to-json" // kvtojson
	"github.com/djschaap/kv-to-json/internal/parsedoc"
	"github.com/djschaap/logevent/fromenv"
	"log"
	"os"
)

var (
	buildDt string
	commit  string
	version string
)

var app *kvtojson.Sess

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	enableTrace := os.Getenv("ENABLE_TRACE")
	var doc string = request.Body
	if enableTrace == "1" {
		requestJSON, _ := json.Marshal(request)
		fmt.Println("TRACE-received-request:\n", string(requestJSON))
		//fmt.Println("TRACE-received-doc:\n", doc)
	}
	headers, message, _ := parsedoc.ParseDoc(doc)
	logEvent := parsedoc.ConvertToLogEvent(headers, message)

	err := app.SendOne(logEvent)
	if err != nil {
		log.Fatal(err)
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

func printVersion() {
	fmt.Println("kv-to-json lambda  Version:",
		version, " Commit:", commit,
		" Built at:", buildDt)
}

func main() {
	printVersion()

	sender, err := fromenv.GetMessageSenderFromEnv()
	if err != nil {
		log.Fatal("Error initializing output:", err)
	}
	err = sender.OpenSvc()
	if err != nil {
		log.Fatal(err)
	}
	app = kvtojson.New(sender, 2)

	lambda.Start(handleRequest)
}
