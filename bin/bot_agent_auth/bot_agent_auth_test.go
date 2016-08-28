package main

import (
	"testing"

	"os"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestCreateRequestConsumer(t *testing.T) {
	createRequestConsumer, err := createRequestConsumer("prefix", "nsqd", "nsqlookupd", "testbot", "auth-api", "auth-app-name", "auth-app-pw", "asdfasdf", "")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(createRequestConsumer, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
