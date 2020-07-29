package main

import (
	"flag"
	"fmt"
	"github.com/djschaap/kv-to-json" // kvtojson
	"github.com/djschaap/kv-to-json/internal/parsedoc"
	"github.com/djschaap/logevent/fromenv"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type webServer struct{}

var (
	app     *kvtojson.Sess
	buildDt string
	commit  string
	version string
)

func (s *webServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// assume health probe; return 200 OK
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"OK"}`))
	case "POST":
		// parse r, send to logevent
		w.Header().Set("Content-Type", "application/json")
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// 406 Not Acceptable
			log.Println("unable to read body:", err)
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(`{"message":"unable to read body"}`))
		}
		headers, message, err := parsedoc.ParseDoc(string(reqBody))
		if err != nil {
			// 406 Not Acceptable
			log.Println("unable to parse body:", err)
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(`{"message":"unable to parse body"}`))
		}
		logEvent := parsedoc.ConvertToLogEvent(headers, message)
		err = app.SendOne(logEvent)
		if err != nil {
			// 503 Service Unavailable
			log.Println("send failed:", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"message":"send failed"}`))
		} else {
			// 200 OK
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"sent"}`))
		}
	default:
		// 405 Method Not Allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(``))
	}
}

func printVersion() {
	fmt.Println("kv-to-json cli  Version:",
		version, " Commit:", commit,
		" Built at:", buildDt)
}

func main() {
	printVersion()

	printVersion := flag.Bool("v", false, "print version and exit")
	webPort := flag.Int("p", 0, "web port (else, read stdin and exit)")
	flag.Parse()
	if *printVersion {
		os.Exit(0)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	sender, err := fromenv.GetMessageSenderFromEnv()
	if err != nil {
		log.Fatal("Error initializing output:", err)
	}
	err = sender.OpenSvc()
	if err != nil {
		log.Fatal("ERROR calling sender.OpenSvc():", err)
	}

	app = kvtojson.New(sender)

	if *webPort > 0 {
		// start web server
		s := &webServer{}
		http.Handle("/state/alert/v1", s)
		err := http.ListenAndServe(fmt.Sprintf(":%d", *webPort), nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// read stdin and exit
		var data []byte
		data, _ = ioutil.ReadAll(os.Stdin)
		headers, message, _ := parsedoc.ParseDoc(string(data))
		logEvent := parsedoc.ConvertToLogEvent(headers, message)

		err := app.SendOne(logEvent)
		if err != nil {
			log.Fatal(err)
		}
	}
}
