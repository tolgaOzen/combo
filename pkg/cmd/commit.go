package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/tolgaOzen/combo/internal"
	"github.com/tolgaOzen/combo/pkg/git"
	"github.com/tolgaOzen/combo/pkg/prompt"
)

// Define the Bubble Tea model
type model struct {
	message  string
	choice   string
	quitting bool
}

// Init Initial model setup
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles user input and state changes
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "", tea.KeyEnter.String():
			return model{message: m.message, choice: "yes", quitting: true}, tea.Quit
		case "n", "N":
			return model{message: m.message, choice: "no", quitting: true}, tea.Quit
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	if m.quitting {
		if m.choice == "yes" {
			return fmt.Sprintf("%s\n\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render("✔ Commit executed successfully!"))
		}
		return fmt.Sprintf("%s\n\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render("✘ Commit aborted."))
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
	brand := brandStyle.Render("Welcome to Combo CLI")
	header := headerStyle.Render("Here’s your commit message:")
	message := messageStyle.Render(fmt.Sprintf("➤ %s", m.message))
	prompt := promptStyle.Render("Would you like to use this message? (Y/n):")

	// Combine output
	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s", brand, header, message, prompt)
}

// NewCommitCommand Commit command logic with Bubble Tea integration
func NewCommitCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "commit",
		Short: "Commit the changes",
		RunE:  commit(),
		Args:  cobra.NoArgs,
	}
	return command
}

func commit() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error getting user home directory: %v", err)
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
			log.Fatalf("Error ensuring configuration: %v", err)
		}

		// Load configuration
		config, err := LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Error loading configuration: %v", err)
		}

		// Retrieve values from the config
		apiKey, ok := config["openai_api_key"]
		if !ok {
			log.Fatal("openai_api_key is missing in the configuration file")
		}

		locale, ok := config["prompt_locale"]
		if !ok {
			locale = "" // Default locale
		}

		maxLengthStr, ok := config["prompt_max_length"]
		if !ok {
			maxLengthStr = "72" // Default max length
		}
		maxLength := 72 // Convert string to int
		if ml, err := strconv.Atoi(maxLengthStr); err == nil {
			maxLength = ml
		}

		// Initialize the OpenAI client
		client := internal.NewOpenAIClient(apiKey)

		// Generate a prompt
		p, err := prompt.GeneratePrompt(
			prompt.Conventional,
			prompt.WithLocale(prompt.Locale(locale)),
			prompt.WithMaxLength(maxLength),
		)
		if err != nil {
			log.Fatal("Prompt generation error")
		}

		diff, err := git.GetDifferences()
		if err != nil {
			log.Fatal("Prompt generation error")
		}

		// Prepare the chat completion request
		request := internal.CreateChatCompletionRequest(p, diff)

		// Set up a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Send the chat completion request and handle the response
		response, err := client.SendChatCompletionRequest(ctx, request)
		if err != nil {
			log.Fatalf("ChatCompletion error: %v\n", err)
		}

		// Generated commit message
		message := response.Choices[0].Message.Content

		// Bubble Tea program setup
		program := tea.NewProgram(&model{message: message})
		mod, err := program.Run()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}

		fmt.Println(mod)

		return nil
	}
}

// LoadConfig loads key-value pairs from a configuration file
func LoadConfig(filePath string) (map[string]string, error) {
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
	dir := filepath.Dir(filePath)

	// Create directory if it does not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
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
