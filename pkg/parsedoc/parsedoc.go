package parsedoc

import (
	"bufio"
	//"fmt"
	"regexp"
	"strings"
)

func ParseDoc(doc string) (map[string]string, map[string]string, error) {
	var headers_done bool
	var headers, message map[string]string
	headers = make(map[string]string)
	message = make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(doc))
	blank_line_re := regexp.MustCompile(`^\s*$`)
	re := regexp.MustCompile(`^(\S+):\s*(.*)`)
	for scanner.Scan() {
		if blank_line_re.MatchString(scanner.Text()) {
			headers_done = true
			continue
		}
		kv := re.FindStringSubmatch(scanner.Text())
		if len(kv) < 2 {
			continue
		}
		//fmt.Println("k   ", kv[1])
		//fmt.Println("  v ", kv[2])
		if headers_done {
			message[kv[1]] = kv[2]
		} else {
			headers[kv[1]] = kv[2]
		}
	}
	return headers, message, nil
}
