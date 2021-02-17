package kvtojson

import (
	"github.com/djschaap/logevent"
	"log"
	"time"
)

// Sess stores kvtojson session state.
type Sess struct {
	messageSender logevent.MessageSender
	retries       int
	retryDelay    time.Duration
	traceoutput   bool
}

// SendOne sends a single LogEvent message.
func (sessObj *Sess) SendOne(logEvent logevent.LogEvent) error {
	err := sessObj.messageSender.SendMessage(logEvent)
	if err != nil {
		log.Println("retrying after error from SendMessage:", err)
		for i := 0; i < sessObj.retries; i++ {
			time.Sleep(sessObj.retryDelay)
			err = sessObj.messageSender.SendMessage(logEvent)
			if err == nil {
				log.Printf("success on retry i=%d", i)
				break
			} else {
				log.Printf("error from SendMessage i=%d: %s", i, err)
			}
		}
	}
	return err
}

// New creates a new kvtojson object/sesion.
func New(
	messageSender logevent.MessageSender,
	retries int,
) *Sess {
	s := Sess{
		messageSender: messageSender,
		retries:       retries,
		retryDelay:    2 * time.Second,
	}
	return &s
}
