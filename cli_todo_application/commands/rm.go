package commands

import (
	"cli_todo_application/storage"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [task number]",
	Short: "Remove a task from your TODO list",
	Long:  `Remove a task from your TODO list by providing its number from the list.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskNum, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task number: %s", args[0])
		}

		if taskNum < 1 {
			return fmt.Errorf("task number must be greater than 0")
		}

		store, err := storage.NewTaskStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		defer store.Close()

		task, err := store.DeleteTask(taskNum)
		if err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		fmt.Printf("You have deleted the \"%s\" task.\n", task.Description)
		return nil
	},
}