package sendsns

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var svc *sns.SNS

func OpenSvc() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = sns.New(sess)
	return
}

func SendMessage(topic_arn string, headers map[string]string, message map[string]string) {
	message_json_bytes, _ := json.Marshal(message)

	var message_attributes map[string]*sns.MessageAttributeValue
	message_attributes = make(map[string]*sns.MessageAttributeValue)
	for k, v := range headers {
		message_attributes[k] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
	}
	/* message_attributes["content_type"] = &sns.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String("application/json"),
	} */

	//headers_json_bytes, _ := json.Marshal(headers)
	//fmt.Println("TRACE-Headers:\n", string(headers_json_bytes))
	//fmt.Println("TRACE-Message:\n", string(message_json_bytes))

	result, err := svc.Publish(&sns.PublishInput{
		MessageAttributes: message_attributes,
		Message:           aws.String(string(message_json_bytes)),
		TopicArn:          &topic_arn,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", *result.MessageId)
}
