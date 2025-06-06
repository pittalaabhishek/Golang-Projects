package commands

import (
	"cli_todo_application/storage"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task to your TODO list",
	Long:  `Add a new task to your TODO list with the provided description.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		description := strings.Join(args, " ")
		if strings.TrimSpace(description) == "" {
			return fmt.Errorf("task description cannot be empty")
		}

		store, err := storage.NewTaskStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		defer store.Close()

		task, err := store.AddTask(description)
		if err != nil {
			return fmt.Errorf("failed to add task: %w", err)
		}

		fmt.Printf("Added \"%s\" to your task list.\n", task.Description)
		return nil
	},
}