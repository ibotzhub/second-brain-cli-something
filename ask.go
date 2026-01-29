package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/brain-cli/internal/brain"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask a question about your notes",
	Long: `Ask a natural language question and get an answer based on your notes.

Examples:
  brain ask "What have I learned about database optimization?"
  brain ask "How should I handle errors in Go?"
  brain ask "What are the team's coding standards?"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		question := args[0]

		b, err := brain.New()
		if err != nil {
			return fmt.Errorf("failed to initialize brain: %w", err)
		}

		// Search for relevant notes
		results, err := b.Search(question, 5, nil)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("I don't have any notes that might answer that question.")
			fmt.Println("Try adding some notes first with: brain add \"your insight\"")
			return nil
		}

		fmt.Printf("💡 Based on your notes, here's what I found:\n\n")
		fmt.Printf("Question: %s\n\n", question)
		
		fmt.Println("Relevant notes:")
		for i, result := range results {
			fmt.Printf("%d. %s\n", i+1, result.Note.Content)
			if len(result.Note.Tags) > 0 {
				fmt.Printf("   Tags: %v\n", result.Note.Tags)
			}
			fmt.Printf("   Relevance: %.0f%%\n\n", result.Similarity*100)
		}

		// TODO: In the future, we could use an LLM to synthesize these notes
		// into a natural language answer. For now, just show the relevant notes.

		return nil
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}
