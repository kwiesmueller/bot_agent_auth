package handler

import (
	"testing"

	"fmt"

	. "github.com/bborbe/assert"
	"github.com/bborbe/bot_agent/message"
	h "github.com/bborbe/bot_agent_auth/handler"
)

func TestImplementsHandler(t *testing.T) {
	c := New("", "", nil)
	var i *h.Handler
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchTrue(t *testing.T) {
	c := New("/auth", "", nil)
	match := c.Match(&message.Request{
		Message: "/auth remove group admin from user tester",
	})
	if err := AssertThat(match, Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchFalse(t *testing.T) {
	c := New("/auth", "", nil)
	match := c.Match(&message.Request{
		Message: "/auth remove group admin from user",
	})
	if err := AssertThat(match, Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageSuccess(t *testing.T) {
	userName := "tester"
	groupName := "admin"
	counter := 0
	c := New("/auth", "", func(_groupName string, _userName string) error {
		if err := AssertThat(_groupName, Is(groupName)); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(_userName, Is(userName)); err != nil {
			t.Fatal(err)
		}
		counter++
		return nil
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	responses, err := c.HandleMessage(&message.Request{
		Message:   fmt.Sprintf("/auth remove group %s from user %s", groupName, userName),
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
	userName := "tester"
	groupName := "admin"
	counter := 0
	c := New("/auth", "", func(_groupName string, _userName string) error {
		if err := AssertThat(_groupName, Is(groupName)); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(_userName, Is(userName)); err != nil {
			t.Fatal(err)
		}
		counter++
		return fmt.Errorf("foo")
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	responses, err := c.HandleMessage(&message.Request{
		Message:   fmt.Sprintf("/auth remove group %s from user %s", groupName, userName),
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
