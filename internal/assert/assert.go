package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func SliceContains(t *testing.T, sl []string, target string) {
	t.Helper()
	for _, val := range sl {
		if val == target {
			return
		}
	}
	t.Errorf("The slice does not contain %s", target)

}

func SliceDoesNotContain(t *testing.T, sl []string, target string) {
	t.Helper()
	for _, val := range sl {
		if val == target {
			t.Errorf("The slice contains %s", target)
		}
	}
}
