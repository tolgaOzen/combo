package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewConfigCommand - returns a new cobra command for config
func NewConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Manage application configuration",
	}

	// Add subcommands
	command.AddCommand(newConfigSetCommand())
	command.AddCommand(newConfigGetCommand())

	return command
}

// newConfigSetCommand - returns a cobra command for setting config
func newConfigSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration key-value pair",
		Args:  cobra.ExactArgs(2), // Requires exactly 2 arguments: key and value
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			configPath, err := getConfigPath()
			if err != nil {
				return err
			}

			// Load or create the config file
			config, err := loadOrCreateConfig(configPath)
			if err != nil {
				return err
			}

			// Update the key-value pair
			config[key] = value

			// Save the updated config
			if err := saveConfig(configPath, config); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			fmt.Printf("Configuration set: %s=%s\n", key, value)
			return nil
		},
	}
}

// newConfigGetCommand - returns a cobra command for getting config values
func newConfigGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value by key",
		Args:  cobra.ExactArgs(1), // Requires exactly 1 argument: key
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			configPath, err := getConfigPath()
			if err != nil {
				return err
			}

			// Load the config
			config, err := loadOrCreateConfig(configPath)
			if err != nil {
				return err
			}

			// Get the value for the key
			value, exists := config[key]
			if !exists {
				return fmt.Errorf("key %s not found in configuration", key)
			}

			fmt.Printf("%s=%s\n", key, value)
			return nil
		},
	}
}

// getConfigPath - retrieves the configuration file path
func getConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(dir, ".combo", "config"), nil
}

// loadOrCreateConfig - loads the config file or creates a new one if it doesn't exist
func loadOrCreateConfig(filePath string) (map[string]string, error) {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Initialize the config map
	config := make(map[string]string)

	// Read the file if it exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File does not exist, return an empty config
		return config, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Parse the file line by line
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
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return config, nil
}

// saveConfig - saves the configuration map to the file
func saveConfig(filePath string, config map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	for key, value := range config {
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return fmt.Errorf("failed to write to config file: %w", err)
		}
	}

	return nil
}
