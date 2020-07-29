package parsedoc

import (
	"bufio"
	"github.com/djschaap/logevent"
	"regexp"
	"strings"
)

// ConvertToLogEvent converts ParseDoc output to a LogEvent.
func ConvertToLogEvent(headers, innerMessage map[string]string) logevent.LogEvent {
	attr := logevent.Attributes{}
	content := logevent.MessageContent{
		Event: innerMessage,
	}

	if headers["customer_code"] != "" {
		attr.CustomerCode = headers["customer_code"]
	}
	if headers["host"] != "" {
		attr.Host = headers["host"]
		content.Host = headers["host"]
	}
	if headers["index"] != "" {
		content.Index = headers["index"]
	}
	if headers["source"] != "" {
		attr.Source = headers["source"]
		content.Source = headers["source"]
	}
	if headers["source_environment"] != "" {
		attr.SourceEnvironment = headers["source_environment"]
	}
	if headers["sourcetype"] != "" {
		attr.Sourcetype = headers["sourcetype"]
		content.Sourcetype = headers["sourcetype"]
	}
	if headers["type"] != "" {
		attr.Type = headers["type"]
	}

	event := logevent.LogEvent{
		attr,
		content,
	}
	return event
}

// ParseDoc produces headers & message maps, given an input document.
func ParseDoc(doc string) (map[string]string, map[string]string, error) {
	var headersDone bool
	var headers, message map[string]string
	headers = make(map[string]string)
	message = make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(doc))
	blankLineRegex := regexp.MustCompile(`^\s*$`)
	re := regexp.MustCompile(`^(\S+):\s*(.*)`)
	for scanner.Scan() {
		if blankLineRegex.MatchString(scanner.Text()) {
			headersDone = true
			continue
		}
		kv := re.FindStringSubmatch(scanner.Text())
		if len(kv) < 2 {
			continue
		}
		if headersDone {
			message[kv[1]] = kv[2]
		} else {
			headers[kv[1]] = kv[2]
		}
	}
	return headers, message, nil
}
