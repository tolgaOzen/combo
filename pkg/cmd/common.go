package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// LoadConfig loads key-value pairs from a configuration file
func LoadConfig(filePath string) (map[string]string, error) {
	// Sanitize and validate the file path
	filePath = filepath.Clean(filePath)
	trustedDir, err := getConfigDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted directory: %w", err)
	}
	if !strings.HasPrefix(filePath, trustedDir) {
		return nil, fmt.Errorf("file path is outside the trusted directory: %s", filePath)
	}

	config := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid configuration line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// EnsureConfig ensures the directory and configuration file exist
func EnsureConfig(filePath, defaultContent string) error {
	// Sanitize and validate the file path
	filePath = filepath.Clean(filePath)
	trustedDir, err := getConfigDirectory()
	if err != nil {
		return fmt.Errorf("failed to get trusted directory: %w", err)
	}
	if !strings.HasPrefix(filePath, trustedDir) {
		return fmt.Errorf("file path is outside the trusted directory: %s", filePath)
	}

	dir := filepath.Dir(filePath)

	// Create directory if it does not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create file if it does not exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create config file %s: %w", filePath, err)
		}
		defer file.Close()

		// Write default content to the file
		if _, err := file.WriteString(defaultContent); err != nil {
			return fmt.Errorf("failed to write to config file %s: %w", filePath, err)
		}
	}

	return nil
}

func getConfigDirectory() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(dir, ".combo"), nil
}

// runGitCommit executes the git commit command
func runGitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runGitCommit executes the git commit command
func runGitBranch(name string) error {
	cmd := exec.Command("git", "checkout", "-b", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
