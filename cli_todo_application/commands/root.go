package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "task",
	Short: "A CLI for managing your TODOS",
	Long: `task is a CLI for managing your TODOs.

This application allows you to add, list, complete, and manage your tasks
from the command line with persistent storage.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(completedCmd)
}