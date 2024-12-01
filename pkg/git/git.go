package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// DiffResult encapsulates the staged diff results
type DiffResult struct {
	Files []string
	Diff  string
}

// GetDifferences retrieves and prints staged differences.
func GetDifferences() (string, error) {
	// Fetch staged differences
	result, err := FetchStagedDiff()
	if err != nil {
		return "", fmt.Errorf("error fetching staged differences: %w", err)
	}

	// Handle no staged changes
	if result == nil || len(result.Files) == 0 {
		return "", fmt.Errorf("no staged changes found. Stage your changes manually, or use the `--all` flag")
	}

	// Format and return the file names
	return result.Diff, nil
}

// FetchStagedDiff retrieves staged changes, including file names and diff content.
func FetchStagedDiff() (*DiffResult, error) {
	// Fetch the list of staged files
	filesOut, err := runGitCommand([]string{"diff", "--cached", "--name-only"})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve staged file names: %w", err)
	}

	// If no files are staged, return nil
	files := strings.Split(strings.TrimSpace(filesOut), "\n")
	if len(files) == 1 && files[0] == "" {
		return nil, nil
	}

	// Fetch the full staged diff
	diffOut, err := runGitCommand([]string{"diff", "--cached", "--diff-algorithm=minimal"})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve staged diff: %w", err)
	}

	return &DiffResult{
		Files: files,
		Diff:  diffOut,
	}, nil
}

// runGitCommand executes a Git command and returns the output as a string.
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
