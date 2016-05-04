package message_handler

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/bot_agent/message_handler"
)

func TestImplementsAgent(t *testing.T) {
	c := New(nil, nil)
	var i *message_handler.MessageHandler
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}
