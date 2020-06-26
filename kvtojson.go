package kvtojson

import (
	"github.com/djschaap/logevent"
)

type Sess struct {
	messageSender logevent.MessageSender
	traceoutput   bool
}

func (self *Sess) SendOne(logEvent logevent.LogEvent) error {
	return self.messageSender.SendMessage(logEvent)
}

func New(
	messageSender logevent.MessageSender,
) *Sess {
	s := Sess{
		messageSender: messageSender,
	}
	return &s
}
