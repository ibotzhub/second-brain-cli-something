package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yourusername/brain-cli/internal/brain"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long: `List all notes in your brain, sorted by most recent first.

Examples:
  brain list
  brain list --limit 10
  brain list --tags go,performance`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		b, err := brain.New()
		if err != nil {
			return fmt.Errorf("failed to initialize brain: %w", err)
		}

		notes, err := b.ListNotes(tags)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(notes) == 0 {
			fmt.Println("No notes found.")
			fmt.Println("Add your first note with: brain add \"your insight here\"")
			return nil
		}

		// Sort by timestamp, most recent first
		sort.Slice(notes, func(i, j int) bool {
			return notes[i].Timestamp.After(notes[j].Timestamp)
		})

		// Apply limit
		if limit > 0 && limit < len(notes) {
			notes = notes[:limit]
		}

		fmt.Printf("Found %d note(s):\n\n", len(notes))
		for i, note := range notes {
			fmt.Printf("%d. [%s] %s\n", i+1, note.Timestamp.Format("2006-01-02 15:04"), note.Content)
			if len(note.Tags) > 0 {
				fmt.Printf("   Tags: %v\n", note.Tags)
			}
			if note.Project != "" {
				fmt.Printf("   Project: %s\n", note.Project)
			}
			fmt.Printf("   ID: %s\n\n", note.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntP("limit", "l", 0, "Maximum number of notes to show (0 = all)")
	listCmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags")
}
