package message_handler

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/api"
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent/message_handler"
)

func TestImplementsAgent(t *testing.T) {
	c := New(nil, nil, nil)
	var i *message_handler.MessageHandler
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageCreateApplication(t *testing.T) {
	counter := 0
	c := New(func(applicationName string) (*api.ApplicationPassword, error) {
		counter++
		pw := api.ApplicationPassword("bar")
		return &pw, nil
	}, nil, nil)
	c.HandleMessage(&message.Request{
		Message: "/auth application create foo",
	})
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageDeleteApplication(t *testing.T) {
	counter := 0
	c := New(nil, func(applicationName string) error {
		counter++
		return nil
	}, nil)
	c.HandleMessage(&message.Request{
		Message: "/auth application delete foo",
	})
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageExistsApplication(t *testing.T) {
	counter := 0
	c := New(nil, nil, func(applicationName string) (*bool, error) {
		counter++
		result := true
		return &result, nil
	})
	c.HandleMessage(&message.Request{
		Message: "/auth application exists foo",
	})
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
