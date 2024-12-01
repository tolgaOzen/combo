package prompt

import (
	"encoding/json"
	"fmt"
)

// Locale defines the supported languages for commit messages.
type Locale string

const (
	EnUS Locale = "en-US" // Default locale for English (United States).
	EnGB Locale = "en-GB" // English (United Kingdom).
	FrFR Locale = "fr-FR" // French (France).
	EsES Locale = "es-ES" // Spanish (Spain).
	DeDE Locale = "de-DE" // German (Germany).
	ItIT Locale = "it-IT" // Italian (Italy).
	KoKR Locale = "ko-KR" // Korean (South Korea).
	JaJP Locale = "ja-JP" // Japanese (Japan).
	ZhCN Locale = "zh-CN" // Chinese (Simplified, China).
	ZhTW Locale = "zh-TW" // Chinese (Traditional, Taiwan).
	PtBR Locale = "pt-BR" // Portuguese (Brazil).
	RuRU Locale = "ru-RU" // Russian (Russia).
	ArSA Locale = "ar-SA" // Arabic (Saudi Arabia).
	HiIN Locale = "hi-IN" // Hindi (India).
)

func (c Locale) String() string {
	return string(c)
}

// CommitType defines the type of commit messages.
type CommitType string

const (
	Build    CommitType = "build"
	Chore    CommitType = "chore"
	CI       CommitType = "ci"
	Docs     CommitType = "docs"
	Feat     CommitType = "feat"
	Fix      CommitType = "fix"
	Perf     CommitType = "perf"
	Refactor CommitType = "refactor"
	Revert   CommitType = "revert"
	Style    CommitType = "style"
	Test     CommitType = "test"
)

func (c CommitType) String() string {
	return string(c)
}

// commitTypeDescriptions maps CommitType to its description.
var commitTypeDescriptions = map[CommitType]string{
	Build:    "Changes that affect the build system or external dependencies (e.g., gulp, broccoli, npm).",
	Chore:    "Other changes that don't modify src or test files.",
	CI:       "Changes to CI configuration files and scripts (e.g., Travis, Circle, BrowserStack, SauceLabs).",
	Docs:     "Documentation only changes.",
	Feat:     "A new feature.",
	Fix:      "A bug fix.",
	Perf:     "A code change that improves performance.",
	Refactor: "A code change that neither fixes a bug nor adds a feature.",
	Revert:   "Reverts a previous commit.",
	Style:    "Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc.).",
	Test:     "Adding missing tests or correcting existing tests.",
}

// CommitStyle defines the type of commit message formats.
type CommitStyle string

const (
	Empty        CommitStyle = ""
	Conventional CommitStyle = "conventional"
)

// commitStyleFormats maps CommitStyle to its message format.
var commitStyleFormats = map[CommitStyle]string{
	Empty:        "<commit message>",
	Conventional: "<type>(<optional scope>): <commit message>",
}

// SpecifyCommitFormat returns the format specification for a given CommitStyle.
func SpecifyCommitFormat(style CommitStyle) (string, error) {
	format, exists := commitStyleFormats[style]
	if !exists {
		return "", fmt.Errorf("invalid commit style: %s", style)
	}
	return fmt.Sprintf("The output response must be in format:\n%s", format), nil
}

// generateCommitTypeDescriptions converts the commitTypeDescriptions map to a JSON string
// and formats it for inclusion in the prompt.
func generateCommitTypeDescriptions() (string, error) {
	descriptionMap := make(map[string]string)
	for key, value := range commitTypeDescriptions {
		descriptionMap[string(key)] = value
	}

	jsonDescriptions, err := json.MarshalIndent(descriptionMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to generate commit descriptions: %v", err)
	}
	return fmt.Sprintf(
		"Choose a type from the type-to-description JSON below that best describes the git diff:\n%s",
		string(jsonDescriptions),
	), nil
}

// Config holds configuration for generating a commit message prompt.
type Config struct {
	Locale             Locale // The language of the commit message.
	MaxLength          int    // Maximum allowed character length for the message.
	CommitDescriptions string // Description of available commit types.
	CommitFormat       string // Commit message format (e.g., "<type>(<scope>): <message>").
}

// Option defines a functional option for configuring the prompt generation.
type Option func(*Config)

// WithLocale sets the locale in the configuration.
func WithLocale(locale Locale) Option {
	return func(cfg *Config) {
		cfg.Locale = locale
	}
}

// WithMaxLength sets the maximum length in the configuration.
func WithMaxLength(maxLength int) Option {
	return func(cfg *Config) {
		cfg.MaxLength = maxLength
	}
}

// GeneratePrompt generates a concise prompt for creating git commit messages.
func GeneratePrompt(style CommitStyle, opts ...Option) (string, error) {
	// Default configuration
	config := &Config{
		Locale:    EnUS, // Default to en-US
		MaxLength: 72,   // Default max length
	}

	// Apply functional options
	for _, opt := range opts {
		opt(config)
	}

	commitDescriptions := ""

	// Generate commit type descriptions if the commit style is conventional.
	if style == Conventional {
		var err error
		commitDescriptions, err = generateCommitTypeDescriptions()
		if err != nil {
			return "", err
		}
	}

	commitFormat, err := SpecifyCommitFormat(style)
	if err != nil {
		return "", err
	}

	config.CommitDescriptions = commitDescriptions
	config.CommitFormat = commitFormat

	// Build the prompt using the configuration.
	return buildPrompt(config)
}

// buildPrompt constructs the final prompt string based on the given configuration.
func buildPrompt(config *Config) (string, error) {
	// Validate configuration
	if config.Locale.String() == "" {
		return "", fmt.Errorf("locale cannot be empty")
	}
	if config.MaxLength <= 0 {
		return "", fmt.Errorf("maxLength must be greater than 0")
	}
	if config.CommitFormat == "" {
		return "", fmt.Errorf("commitFormat cannot be empty")
	}

	// Construct the prompt
	return fmt.Sprintf(
		`Write a concise and relevant git commit message for the given code diff:
Language: %s
Maximum length: %d characters.
Focus: Only include details about the code changes. Avoid unnecessary information such as translations or extra explanations.
Format: Use the specified commit message format:
%s
%s
`,
		config.Locale.String(),
		config.MaxLength,
		config.CommitDescriptions,
		config.CommitFormat,
	), nil
}
