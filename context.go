package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/brain-cli/internal/brain"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Show relevant notes for your current context",
	Long: `Analyzes your current directory, git repo, and recent files
to surface relevant notes from your brain.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := brain.New()
		if err != nil {
			return fmt.Errorf("failed to initialize brain: %w", err)
		}

		context := detectContext()
		
		fmt.Printf("📍 Current context: %s\n", context.Description)
		if context.Project != "" {
			fmt.Printf("📦 Project: %s\n", context.Project)
		}
		fmt.Println()

		results, err := b.GetContextualNotes(context)
		if err != nil {
			return fmt.Errorf("failed to get contextual notes: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("No relevant notes found for this context.")
			fmt.Println("Try adding some notes with: brain add \"your insight here\"")
			return nil
		}

		fmt.Printf("Found %d relevant note(s) for this context:\n\n", len(results))
		for i, result := range results {
			fmt.Printf("%d. %s\n", i+1, result.Note.Content)
			if len(result.Note.Tags) > 0 {
				fmt.Printf("   Tags: %v\n", result.Note.Tags)
			}
			fmt.Printf("   Relevance: %.2f%%\n\n", result.Similarity*100)
		}

		return nil
	},
}

type Context struct {
	Directory   string
	Project     string
	Description string
	Keywords    []string
}

func detectContext() Context {
	cwd, _ := os.Getwd()
	
	ctx := Context{
		Directory:   cwd,
		Description: filepath.Base(cwd),
		Keywords:    []string{},
	}

	// Try to get git repo info
	if gitRoot := getGitRoot(); gitRoot != "" {
		ctx.Project = filepath.Base(gitRoot)
		ctx.Description = ctx.Project
		
		// Get recent commit messages for context
		if commits := getRecentCommits(); len(commits) > 0 {
			ctx.Keywords = append(ctx.Keywords, commits...)
		}
	}

	// Add directory name as keyword
	ctx.Keywords = append(ctx.Keywords, filepath.Base(cwd))

	return ctx
}

func getGitRoot() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getRecentCommits() []string {
	cmd := exec.Command("git", "log", "-5", "--pretty=format:%s")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	
	commits := strings.Split(string(output), "\n")
	var keywords []string
	for _, commit := range commits {
		if commit != "" {
			keywords = append(keywords, commit)
		}
	}
	return keywords
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
