package user_delete

import (
	"testing"

	auth_model "github.com/bborbe/auth/model"

	"fmt"

	"os"

	. "github.com/bborbe/assert"
	"github.com/bborbe/bot_agent/api"
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

func TestMatchTrue(t *testing.T) {
	c := New("/auth", nil)
	match := c.Match(&api.Request{
		Message: "/auth user delete tester",
	})
	if err := AssertThat(match, Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchFalse(t *testing.T) {
	c := New("/auth", nil)
	match := c.Match(&api.Request{
		Message: "/auth user delete",
	})
	if err := AssertThat(match, Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestHandleMessageSuccess(t *testing.T) {
	username := "tester"
	counter := 0
	c := New("/auth", func(_username auth_model.UserName) error {
		if err := AssertThat(_username.String(), Is(username)); err != nil {
			t.Fatal(err)
		}
		counter++
		return nil
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	responses, err := c.HandleMessage(&api.Request{
		Message:   "/auth user delete tester",
		AuthToken: auth_model.AuthToken(username),
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
	username := "tester"
	counter := 0
	c := New("/auth", func(_username auth_model.UserName) error {
		if err := AssertThat(_username.String(), Is(username)); err != nil {
			t.Fatal(err)
		}
		counter++
		return fmt.Errorf("foo")
	})
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	responses, err := c.HandleMessage(&api.Request{
		Message:   "/auth user delete tester",
		AuthToken: auth_model.AuthToken(username),
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
