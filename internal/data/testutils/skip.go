package testutils

import (
	"os"
	"testing"
)

// CheckSkipTestContainers checks if the SKIP_TESTCONTAINERS environment variable is set to true.
// Look for this variable in the GitHub Actions workflow (.github/workflows/test.yaml).
func CheckSkipTestContainers(t *testing.T) {
	if os.Getenv("SKIP_TESTCONTAINERS") == "true" {
		t.Skip("Skipping TestContainers tests")
	}
}
