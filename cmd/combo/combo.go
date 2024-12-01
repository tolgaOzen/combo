package main

import (
	"os"

	"github.com/tolgaOzen/combo/pkg/cmd"
)

func main() {
	root := cmd.NewRootCommand()

	commit := cmd.NewCommitCommand()
	root.AddCommand(commit)

	branch := cmd.NewBranchCommand()
	root.AddCommand(branch)

	version := cmd.NewVersionCommand()
	root.AddCommand(version)

	config := cmd.NewConfigCommand()
	root.AddCommand(config)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
