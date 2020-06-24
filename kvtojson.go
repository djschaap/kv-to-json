package kvtojson

import (
	"github.com/djschaap/logevent"
)

type Sess struct {
	messageSender logevent.MessageSender
	topicArn      string
	traceoutput   bool
}

func (self *Sess) SendOne(logEvent logevent.LogEvent) error {
	return self.messageSender.SendMessage(self.topicArn, logEvent)
}

func New(
	messageSender logevent.MessageSender,
	snsTopicArn string,
) *Sess {
	s := Sess{
		messageSender: messageSender,
		topicArn:      snsTopicArn,
	}
	return &s
}
