package main

import (
	"flag"
	"fmt"
	"github.com/djschaap/kv-to-json" // kvtojson
	"github.com/djschaap/kv-to-json/internal/parsedoc"
	"github.com/djschaap/logevent"
	"github.com/djschaap/logevent/senddump"
	"github.com/djschaap/logevent/sendhec"
	"github.com/djschaap/logevent/sendsns"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var (
	buildDt string
	commit  string
	version string
)

func getenvBool(k string) bool {
	v := os.Getenv(k)
	if len(v) > 0 {
		return true
	}
	return false
}

func printVersion() {
	fmt.Println("kv-to-json cli  Version:",
		version, " Commit:", commit,
		" Built at:", buildDt)
}

func main() {
	printVersion()

	printVersion := flag.Bool("v", false, "print version and exit")
	flag.Parse()
	if *printVersion {
		os.Exit(0)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	var data []byte
	data, _ = ioutil.ReadAll(os.Stdin)
	headers, message, _ := parsedoc.ParseDoc(string(data))
	logEvent := parsedoc.ConvertToLogEvent(headers, message)

	var traceOutput bool
	if len(os.Getenv("TRACE_OUTPUT")) > 0 {
		fmt.Println("*** TRACE_OUTPUT is enabled ***")
		traceOutput = true
	}

	snsTopicArn := os.Getenv("TOPIC_ARN")
	hasSnsTopic, _ := regexp.MatchString(`^arn:`, snsTopicArn)
	hecToken := os.Getenv("HEC_TOKEN")
	var sender logevent.MessageSender
	if hasSnsTopic {
		log.Println("Destination: sendsns", snsTopicArn)
		snsSender := sendsns.New(snsTopicArn)
		sender = snsSender
	} else if len(hecToken) > 0 {
		hecUrl := os.Getenv("HEC_URL")
		log.Println("Destination: sendhec", hecUrl)
		hecSender := sendhec.New(hecUrl, hecToken)
		if len(os.Getenv("HEC_INSECURE")) > 0 {
			hecSender.SetHecInsecure(true)
		}
		sender = hecSender
	} else {
		log.Println("WARNING: using senddump; forcing TRACE_OUTPUT")
		dumpSender := senddump.New()
		dumpSender.SetTrace(true)
		sender = dumpSender
	}
	if traceOutput {
		sender.SetTrace(true)
	}
	err = sender.OpenSvc()
	if err != nil {
		log.Fatal("ERROR calling sender.OpenSvc():", err)
	}

	app := kvtojson.New(sender)
	err = app.SendOne(logEvent)
	if err != nil {
		log.Fatal(err)
	}
}
