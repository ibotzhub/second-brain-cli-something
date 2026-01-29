package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "brain",
	Short: "A second brain CLI for storing and retrieving ideas semantically",
	Long: `Brain is a CLI tool that helps you save thoughts, ideas, and learnings,
then surface them when you need them using semantic search and context awareness.

Save anything worth remembering, and let your brain remind you when it's relevant.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can go here
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.brain/config.yaml)")
}
