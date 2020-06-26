package main

import (
	"flag"
	"fmt"
	"github.com/djschaap/kv-to-json" // kvtojson
	"github.com/djschaap/kv-to-json/internal/parsedoc"
	"github.com/djschaap/logevent/sendsns"
	//sendsns "github.com/djschaap/logevent/sendstdout" // DEV HACK
	"io/ioutil"
	"log"
	"os"
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

	var data []byte
	data, _ = ioutil.ReadAll(os.Stdin)
	headers, message, _ := parsedoc.ParseDoc(string(data))
	logEvent := parsedoc.ConvertToLogEvent(headers, message)

	var traceOutput bool
	if len(os.Getenv("TRACE_OUTPUT")) > 0 {
		fmt.Println("*** TRACE_OUTPUT is enabled ***")
		traceOutput = true
	}

	topicArn := os.Getenv("TOPIC_ARN")
	sender := sendsns.New(topicArn)
	if traceOutput {
		sender.SetTrace(true)
	}
	app := kvtojson.New(sender)
	err := sender.OpenSvc()
	if err != nil {
		log.Fatal(err)
	}
	err = app.SendOne(logEvent)
	if err != nil {
		log.Fatal(err)
	}
}
