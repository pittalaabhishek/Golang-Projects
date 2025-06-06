package commands

import (
	"cli_todo_application/storage"
	"fmt"

	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all tasks completed today",
	Long:  `Display all tasks that were completed today.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewTaskStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		defer store.Close()

		tasks, err := store.GetCompletedTasksToday()
		if err != nil {
			return fmt.Errorf("failed to get completed tasks: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("You have not completed any tasks today.")
			return nil
		}

		fmt.Println("You have finished the following tasks today:")
		for _, task := range tasks {
			completedTime := "unknown time"
			if task.CompletedAt != nil {
				completedTime = task.CompletedAt.Format("15:04")
			}
			fmt.Printf("- %s (completed at %s)\n", task.Description, completedTime)
		}

		return nil
	},
}