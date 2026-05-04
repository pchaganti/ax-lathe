package verify_test

import (
	"os"
	"testing"

	"github.com/devenjarvis/lathe/internal/verify"
)

func TestSpawnVerifierMissingClaude(t *testing.T) {
	tutDir := t.TempDir()

	// Remove PATH so claude binary can't be found
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	err := verify.SpawnVerifier("test-slug", tutDir)
	if err == nil {
		t.Error("SpawnVerifier() expected error when claude not in PATH, got nil")
	}
}
