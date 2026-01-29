package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/brain-cli/internal/brain"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search your notes semantically",
	Long: `Search your notes using semantic similarity, not just keywords.
	
Examples:
  brain search "making APIs faster"
  brain search "database optimization" --limit 10
  brain search "performance" --tags go`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		limit, _ := cmd.Flags().GetInt("limit")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		b, err := brain.New()
		if err != nil {
			return fmt.Errorf("failed to initialize brain: %w", err)
		}

		results, err := b.Search(query, limit, tags)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("No matching notes found.")
			return nil
		}

		fmt.Printf("Found %d relevant note(s):\n\n", len(results))
		for i, result := range results {
			fmt.Printf("%d. [%s] %s\n", i+1, result.Note.Timestamp.Format("2006-01-02"), result.Note.Content)
			if len(result.Note.Tags) > 0 {
				fmt.Printf("   Tags: %v\n", result.Note.Tags)
			}
			if result.Note.Project != "" {
				fmt.Printf("   Project: %s\n", result.Note.Project)
			}
			fmt.Printf("   Relevance: %.2f%%\n\n", result.Similarity*100)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntP("limit", "l", 5, "Maximum number of results to return")
	searchCmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags")
}
