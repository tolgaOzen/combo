package cmd

import (
	`context`
	`fmt`
	tea `github.com/charmbracelet/bubbletea`
	`github.com/spf13/cobra`
	`os`
	`path/filepath`
	`time`

	`github.com/tolgaOzen/combo/internal`
	`github.com/tolgaOzen/combo/pkg/git`
	`github.com/tolgaOzen/combo/pkg/prompt`
)

// NewBranchCommand -
func NewBranchCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "branch",
		Short: "Create new branch",
		RunE:  branch(),
		Args:  cobra.NoArgs,
	}
	return command
}

func branch() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		dir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}

		// Configuration file path
		configPath := filepath.Join(dir, ".combo", "config")

		// Ensure the config directory and file exist
		defaultContent := `# Default configuration
openai_api_key=
prompt_locale=en-US
prompt_max_length=72
`
		if err := EnsureConfig(configPath, defaultContent); err != nil {
			return fmt.Errorf("failed to ensure configuration: %w", err)
		}

		// Load configuration
		config, err := LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Retrieve values from the config
		apiKey, ok := config["openai_api_key"]
		if !ok || apiKey == "" {
			return fmt.Errorf("missing or empty 'openai_api_key' in configuration")
		}

		locale, ok := config["prompt_locale"]
		if !ok || locale == "" {
			locale = "en-US" // Default locale
		}

		// Initialize the OpenAI client
		client := internal.NewOpenAIClient(apiKey)

		// Generate a prompt
		p, err := prompt.GenerateBranchNamePrompt(
			prompt.WithLocale(prompt.Locale(locale)),
			prompt.WithMaxLength(30),
		)
		if err != nil {
			return fmt.Errorf("failed to generate prompt: %w", err)
		}

		diff, err := git.GetDifferences()
		if err != nil {
			return fmt.Errorf("failed to get git differences: %w", err)
		}

		// Prepare the chat completion request
		request := internal.CreateChatCompletionRequest(p, diff)

		// Set up a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Send the chat completion request and handle the response
		response, err := client.SendChatCompletionRequest(ctx, request)
		if err != nil {
			return fmt.Errorf("chat completion request failed: %w", err)
		}

		// Ensure the chat completion response and commit message are valid
		if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
			return fmt.Errorf("no commit message generated from the OpenAI response")
		}

		// Generated commit message
		message := response.Choices[0].Message.Content

		// Bubble Tea program setup
		program := tea.NewProgram(&model{message: message})
		mod, err := program.Run()
		if err != nil {
			return fmt.Errorf("bubble tea program encountered an error: %w", err)
		}

		// Check user choice
		if result, ok := mod.(model); ok && result.choice == "yes" {
			// Run git commit command
			if err := runGitBranch(message); err != nil {
				return fmt.Errorf("failed to run git commit: %w", err)
			}
		}

		return nil
	}
}
