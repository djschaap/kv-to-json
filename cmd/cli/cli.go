package main

import (
	"encoding/json"
	"fmt"
	"github.com/djschaap/kv-to-json/pkg/parsedoc"
	"github.com/djschaap/kv-to-json/pkg/sendsns"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	build_dt string
	commit   string
	version  string
)

func main() {
	fmt.Println("kv-to-json cli  Version:", version, " Commit:", commit,
		" Built at:", build_dt)
	var data []byte
	data, _ = ioutil.ReadAll(os.Stdin)
	headers, message, _ := parsedoc.ParseDoc(string(data))
	headers_json_bytes, _ := json.Marshal(headers)
	fmt.Println("Headers:\n", string(headers_json_bytes))
	message_json_bytes, _ := json.Marshal(message)
	fmt.Println("Message:\n", string(message_json_bytes))

	topic_arn := os.Getenv("TOPIC_ARN")
	has_queue, _ := regexp.MatchString(`^arn:`, topic_arn)
	if has_queue {
		sendsns.OpenSvc()
		sendsns.SendMessage(topic_arn, headers, message)
	}
}
