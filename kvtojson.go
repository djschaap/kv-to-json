package kvtojson

import (
	"github.com/djschaap/logevent"
)

// Sess stores kvtojson session state.
type Sess struct {
	messageSender logevent.MessageSender
	traceoutput   bool
}

// SendOne sends a single LogEvent message.
func (sessObj *Sess) SendOne(logEvent logevent.LogEvent) error {
	return sessObj.messageSender.SendMessage(logEvent)
}

// New creates a new kvtojson object/sesion.
func New(
	messageSender logevent.MessageSender,
) *Sess {
	s := Sess{
		messageSender: messageSender,
	}
	return &s
}
