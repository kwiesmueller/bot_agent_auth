package handler

import (
	"testing"

	"fmt"

	. "github.com/bborbe/assert"
	"github.com/bborbe/bot_agent/message"
	h "github.com/bborbe/bot_agent_auth/handler"
)

func TestImplementsHandler(t *testing.T) {
	c := New("", nil)
	var i *h.Handler
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchTrue(t *testing.T) {
	c := New("/auth", nil)
	match := c.Match(&message.Request{
		Message: "/auth register bborbe",
	})
	if err := AssertThat(match, Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchFalse(t *testing.T) {
	c := New("/auth", nil)
	match := c.Match(&message.Request{
		Message: "/auth register",
	})
	if err := AssertThat(match, Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageSuccess(t *testing.T) {
	token := "abc"
	username := "testuser"
	counter := 0
	c := New("/auth", func(authToken string, userName string) error {
		if err := AssertThat(authToken, Is(token)); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(userName, Is(username)); err != nil {
			t.Fatal(err)
		}
		counter++
		return nil
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	responses, err := c.HandleMessage(&message.Request{
		Message:   fmt.Sprintf("/auth register %s", username),
		AuthToken: token,
	})
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(responses, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(responses), Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageFailure(t *testing.T) {
	token := "abc"
	username := "testuser"
	counter := 0
	c := New("/auth", func(authToken string, userName string) error {
		if err := AssertThat(authToken, Is(token)); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(userName, Is(username)); err != nil {
			t.Fatal(err)
		}
		counter++
		return fmt.Errorf("foo")
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	responses, err := c.HandleMessage(&message.Request{
		Message:   fmt.Sprintf("/auth register %s", username),
		AuthToken: token,
	})
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(responses, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(responses), Is(1)); err != nil {
		t.Fatal(err)
	}
}
