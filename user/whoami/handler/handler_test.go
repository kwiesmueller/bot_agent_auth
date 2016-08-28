package handler

import (
	"testing"

	"os"

	. "github.com/bborbe/assert"
	h "github.com/bborbe/bot_agent/message_handler/match"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsHandler(t *testing.T) {
	c := New("", nil)
	var i *h.Handler
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}
