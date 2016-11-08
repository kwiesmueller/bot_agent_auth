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

func TestCreateConfig(t *testing.T) {
	config := createConfig()
	if err := AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
