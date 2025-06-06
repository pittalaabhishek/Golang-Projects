package commands

import (
	"cli_todo_application/storage"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do [task number]",
	Short: "Mark a task on your TODO list as complete",
	Long:  `Mark a task as complete by providing its number from the list.`,
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

		task, err := store.CompleteTask(taskNum)
		if err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		fmt.Printf("You have completed the \"%s\" task.\n", task.Description)
		return nil
	},
}