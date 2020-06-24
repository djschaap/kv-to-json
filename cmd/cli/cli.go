package main

import (
	"encoding/json"
	"fmt"
	"github.com/djschaap/kv-to-json/pkg/parsedoc"
	"github.com/djschaap/kv-to-json/pkg/sendsqs"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {
	var data []byte
	data, _ = ioutil.ReadAll(os.Stdin)
	headers, message, _ := parsedoc.ParseDoc(string(data))
	headers_json_bytes, _ := json.Marshal(headers)
	fmt.Println("Headers:\n", string(headers_json_bytes))
	message_json_bytes, _ := json.Marshal(message)
	fmt.Println("Message:\n", string(message_json_bytes))

	queue_url := os.Getenv("DEST_QUEUE")
	has_queue, _ := regexp.MatchString(`^https`, queue_url)
	if has_queue {
		sendsqs.OpenSvc()
		sendsqs.SendMessage(queue_url, headers, message)
	}
}
