package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand - Creates new root command
func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "combo",
		Short: "Generate commit messages based on git changes effortlessly.",
		Long: `Combo is a CLI tool designed to generate concise and descriptive commit messages automatically. 
It analyzes git changes and provides commit messages adhering to conventional commit standards or other formats of your choice. 
Customize the language, length, and format to fit your workflow.`,
	}
}
