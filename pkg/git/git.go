package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	MaxDiffSize = 4096 // Limit to 4KB to prevent overwhelming LLM
)

// DiffResult encapsulates the staged diff results
type DiffResult struct {
	Files []string
	Diff  string
}

// GetDifferences retrieves staged differences, truncating if needed.
func GetDifferences() (string, error) {
	result, err := FetchStagedDiff()
	if err != nil {
		return "", fmt.Errorf("error fetching staged differences: %w", err)
	}

	if result == nil || len(result.Files) == 0 {
		return "", fmt.Errorf("no staged changes found. Stage your changes manually, or use the `--all` flag")
	}

	diff := result.Diff
	if len(diff) > MaxDiffSize {
		diff = diff[:MaxDiffSize] + "\n[...truncated]"
	}

	return diff, nil
}

// FetchStagedDiff retrieves staged changes using `--patch --compact-summary` for better output.
func FetchStagedDiff() (*DiffResult, error) {
	filesOut, err := runGitCommand([]string{"diff", "--cached", "--name-only"})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve staged file names: %w", err)
	}

	files := strings.Split(strings.TrimSpace(filesOut), "\n")
	if len(files) == 1 && files[0] == "" {
		return nil, nil
	}

	diffOut, err := runGitCommand([]string{"diff", "--cached", "--patch", "--compact-summary"})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve staged diff summary: %w", err)
	}

	return &DiffResult{
		Files: files,
		Diff:  diffOut,
	}, nil
}

// runGitCommand executes a Git command and returns the output.
func runGitCommand(args []string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git command failed: %s, error: %w", out.String(), err)
	}

	return out.String(), nil
}
