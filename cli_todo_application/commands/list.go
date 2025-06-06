package commands

import (
	"cli_todo_application/storage"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Long:  `Display all tasks that are not yet completed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewTaskStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		defer store.Close()

		tasks, err := store.GetIncompleteTasks()
		if err != nil {
			return fmt.Errorf("failed to get tasks: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("You have no incomplete tasks!")
			return nil
		}

		fmt.Println("You have the following tasks:")
		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task.Description)
		}

		return nil
	},
}
