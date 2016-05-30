package handler

import (
	"testing"

	"fmt"

	. "github.com/bborbe/assert"
	"github.com/bborbe/bot_agent/api"
	h "github.com/bborbe/bot_agent/handler"
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
	match := c.Match(&api.Request{
		Message: "/auth token remove abc",
	})
	if err := AssertThat(match, Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchFalse(t *testing.T) {
	c := New("/auth", nil)
	match := c.Match(&api.Request{
		Message: "/auth token remove",
	})
	if err := AssertThat(match, Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageSuccess(t *testing.T) {
	authToken := "abc"
	token := "edf"
	counter := 0
	c := New("/auth", func(_authToken string, _token string) error {
		if err := AssertThat(_authToken, Is(authToken)); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(_token, Is(token)); err != nil {
			t.Fatal(err)
		}
		counter++
		return nil
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	responses, err := c.HandleMessage(&api.Request{
		Message:   fmt.Sprintf("/auth token remove %s", token),
		AuthToken: authToken,
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
	authToken := "abc"
	token := "edf"
	counter := 0
	c := New("/auth", func(_authToken string, _token string) error {
		if err := AssertThat(_authToken, Is(authToken)); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(_token, Is(token)); err != nil {
			t.Fatal(err)
		}
		counter++
		return fmt.Errorf("foo")
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	responses, err := c.HandleMessage(&api.Request{
		Message:   fmt.Sprintf("/auth token remove %s", token),
		AuthToken: authToken,
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
