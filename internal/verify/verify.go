package verify

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed skills/lathe-verify.md
var verifySkillContent string

func SpawnVerifier(slug, tutorialDir string) error {
	tempDir, err := os.MkdirTemp("", "lathe-verify-"+slug+"-")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	// Write the embedded skill into the temp dir so claude can discover it
	skillDir := filepath.Join(tempDir, ".claude", "skills", "lathe-verify")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		return fmt.Errorf("create skill dir: %w", err)
	}
	skillPath := filepath.Join(skillDir, "lathe-verify.md")
	if err := os.WriteFile(skillPath, []byte(verifySkillContent), 0644); err != nil {
		os.RemoveAll(tempDir)
		return fmt.Errorf("write skill: %w", err)
	}

	prompt := fmt.Sprintf(
		"Use the /lathe-verify skill to verify the tutorial. "+
			"LATHE_TUTORIAL_DIR is set to %q.",
		tutorialDir,
	)

	cmd := exec.Command("claude",
		"--add-dir", tempDir,
		"--dangerously-skip-permissions",
		"-p", prompt,
	)
	cmd.Env = append(os.Environ(), "LATHE_TUTORIAL_DIR="+tutorialDir)

	if err := cmd.Start(); err != nil {
		os.RemoveAll(tempDir)
		return fmt.Errorf("start verifier: %w", err)
	}

	// Detach: clean up temp dir when subprocess exits, don't block
	go func() {
		cmd.Wait()
		os.RemoveAll(tempDir)
	}()

	return nil
}
