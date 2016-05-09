package matcher

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestMatchReturnTrueIfLengthEqual(t *testing.T) {
	if err := AssertThat(MatchParts([]string{"a"}, []string{"a"}), Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "a"}, []string{"a", "a"}), Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "a", "a"}, []string{"a", "a", "a"}), Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchReturnFalseIfLengthGreater(t *testing.T) {
	if err := AssertThat(MatchParts([]string{"a"}, []string{"a", "a"}), Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "a"}, []string{"a", "a", "a"}), Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "a", "a"}, []string{"a", "a", "a", "a"}), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchReturnFalseIfLengthLess(t *testing.T) {
	if err := AssertThat(MatchParts([]string{"a", "a"}, []string{"a"}), Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "a", "a"}, []string{"a", "a"}), Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "a", "a", "a"}, []string{"a", "a", "a"}), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestMatchReturnFalseIfValueDiffer(t *testing.T) {
	if err := AssertThat(MatchParts([]string{"a"}, []string{"b"}), Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "a"}, []string{"a", "b"}), Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "a"}, []string{"b", "a"}), Is(false)); err != nil {
		t.Fatal(err)
	}
}
func TestMatchReturnTrueIfValueVariate(t *testing.T) {
	if err := AssertThat(MatchParts([]string{"[NAME]"}, []string{"asdf"}), Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "[NAME]"}, []string{"a", "asdf"}), Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(MatchParts([]string{"a", "b", "[NAME]"}, []string{"a", "b", "asdf"}), Is(true)); err != nil {
		t.Fatal(err)
	}
}
