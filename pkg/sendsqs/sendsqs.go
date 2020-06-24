package sendsqs

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var svc *sqs.SQS

func OpenSvc() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = sqs.New(sess)
	return
}

func SendMessage(queue_url string, headers map[string]string, message map[string]string) {
	message_json_bytes, _ := json.Marshal(message)

	var message_attributes map[string]*sqs.MessageAttributeValue
	message_attributes = make(map[string]*sqs.MessageAttributeValue)
	for k, v := range headers {
		message_attributes[k] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
	}
	/* message_attributes["content_type"] = &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String("application/json"),
	} */

	//headers_json_bytes, _ := json.Marshal(headers)
	//fmt.Println("TRACE-Headers:\n", string(headers_json_bytes))
	//fmt.Println("TRACE-Message:\n", string(message_json_bytes))

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(10),
		MessageAttributes: message_attributes,
		MessageBody:       aws.String(string(message_json_bytes)),
		QueueUrl:          &queue_url,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", *result.MessageId)
}
