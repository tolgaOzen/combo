package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/tolgaOzen/combo/internal"
	"github.com/tolgaOzen/combo/pkg/git"
	"github.com/tolgaOzen/combo/pkg/prompt"
)

// Define the Bubble Tea model
type branchModel struct {
	message  string
	choice   string
	quitting bool
}

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

// Init Initial model setup
func (m branchModel) Init() tea.Cmd {
	return nil
}

// Update handles user input and state changes
func (m branchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "", tea.KeyEnter.String():
			return commitModel{message: m.message, choice: "yes", quitting: true}, tea.Quit
		case "n", "N":
			return commitModel{message: m.message, choice: "no", quitting: true}, tea.Quit
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the UI for branch name confirmation
func (m branchModel) View() string {
	if m.quitting {
		if m.choice == "yes" {
			return fmt.Sprintf(
				"%s\n\n%s\n",
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render("✔ Branch created successfully!"),
				lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Italic(true).Render("You can now switch to your new branch and start working."),
			)
		}
		return fmt.Sprintf(
			"%s\n\n%s\n",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render("✘ Branch creation aborted."),
			lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Italic(true).Render("No branch was created. You can revise and try again."),
		)
	}

	// Define styles
	brandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Padding(0, 0)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		Underline(true)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Italic(true).
		PaddingLeft(2)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("3")).
		PaddingTop(1)

	// Render sections
	brand := brandStyle.Render("Generating your branch name...")
	header := headerStyle.Render("Here’s your suggested branch name:")
	message := messageStyle.Render(fmt.Sprintf("➤ %s", m.message))
	prompt := promptStyle.Render("Would you like to create this branch? (Y/n):")

	// Combine output
	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s", brand, header, message, prompt)
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
		program := tea.NewProgram(&branchModel{message: message})
		mod, err := program.Run()
		if err != nil {
			return fmt.Errorf("bubble tea program encountered an error: %w", err)
		}

		// Check user choice
		if result, ok := mod.(branchModel); ok && result.choice == "yes" {
			// Run git commit command
			if err := runGitBranch(message); err != nil {
				return fmt.Errorf("failed to run git commit: %w", err)
			}
		}

		return nil
	}
}
