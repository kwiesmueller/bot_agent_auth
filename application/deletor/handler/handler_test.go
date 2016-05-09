package handler

import (
	"testing"

	. "github.com/bborbe/assert"
	h "github.com/bborbe/bot_agent_auth/handler"
)

func TestImplementsHandler(t *testing.T) {
	c := New("", "", nil)
	var i *h.Handler
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}
