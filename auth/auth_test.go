package auth

import (
	"testing"

	"fmt"

	. "github.com/bborbe/assert"
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent/message_handler"
)

func TestImplementsAgent(t *testing.T) {
	c := New(nil, nil)
	var i *message_handler.MessageHandler
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageAuthList(t *testing.T) {
	token := "abc"
	counter := 0
	c := New(func(authToken string) ([]string, error) {
		counter++
		if err := AssertThat(authToken, Is(token)); err != nil {
			t.Fatal(err)
		}
		return nil, nil
	}, func(authToken string, authName string) error {
		return nil
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	c.HandleMessage(&message.Request{
		AuthToken: token,
		Message:   "/auth list",
	})
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageAuthCreate(t *testing.T) {
	token := "abc"
	name := "myauth"
	counter := 0
	c := New(func(authToken string) ([]string, error) {
		return nil, nil
	}, func(authToken string, authName string) error {
		counter++
		if err := AssertThat(authToken, Is(token)); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(authName, Is(name)); err != nil {
			t.Fatal(err)
		}
		return nil
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	c.HandleMessage(&message.Request{
		AuthToken: token,
		Message:   fmt.Sprintf("/auth create %s", name),
	})
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
