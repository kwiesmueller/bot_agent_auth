package message_handler

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/bot_agent/api"
)

func TestImplementsAgent(t *testing.T) {
	c := New("")
	var i *api.MessageHandler
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}
