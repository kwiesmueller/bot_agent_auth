package rest

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsRest(t *testing.T) {
	c := New("", "", "", nil, nil)
	var i *Rest
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}
