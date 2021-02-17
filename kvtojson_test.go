package kvtojson

import (
	"github.com/djschaap/logevent"
	"github.com/djschaap/logevent/fromenv"
	"testing"
)

func Test_New(t *testing.T) {
	// env is a local object, but stores state GLOBALLY
	env := fromenv.NewFakeEnv()
	env.Setenv("SENDER_PACKAGE", "senddump")

	var sender logevent.MessageSender
	sender, err := fromenv.GetMessageSenderFromEnv()
	if err != nil {
		t.Errorf("expected success but got error: %s", err)
	}

	var sess *Sess
	sess = New(sender, 0)
	if sess == nil {
		t.Error("expected kvtojson session but got nil")
	}
}
