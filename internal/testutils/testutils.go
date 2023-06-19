package testutils

import (
	"os"
	"testing"
)

func RunAsIntegTest(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("Skipping this integration test. To run this set the env variable INTEGRATION=true")
	}
}
