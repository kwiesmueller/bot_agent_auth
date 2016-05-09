package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestCreateRequestConsumer(t *testing.T) {
	createRequestConsumer, err := createRequestConsumer("prefix", "nsqd", "nsqlookupd", "testbot", "auth-api", "auth-app-name", "auth-app-pw", "asdfasdf")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(createRequestConsumer, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
