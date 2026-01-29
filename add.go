package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/brain-cli/internal/brain"
)

var addCmd = &cobra.Command{
	Use:   "add [note]",
	Short: "Add a new note to your brain",
	Long: `Add a new note, idea, or learning to your brain.
	
Examples:
  brain add "Redis caching reduced API latency by 60%"
  brain add "Use context.WithTimeout for API calls" --tags go,best-practices
  brain add "Team prefers tabs over spaces" --project myapp`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content := args[0]
		tags, _ := cmd.Flags().GetStringSlice("tags")
		project, _ := cmd.Flags().GetString("project")

		b, err := brain.New()
		if err != nil {
			return fmt.Errorf("failed to initialize brain: %w", err)
		}

		note := &brain.Note{
			Content:   content,
			Tags:      tags,
			Project:   project,
			Timestamp: time.Now(),
		}

		if err := b.AddNote(note); err != nil {
			return fmt.Errorf("failed to add note: %w", err)
		}

		fmt.Printf("✓ Note added successfully (ID: %s)\n", note.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringSliceP("tags", "t", []string{}, "Tags for the note")
	addCmd.Flags().StringP("project", "p", "", "Project this note belongs to")
}
