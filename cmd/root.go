package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"kasher/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(taskCmd)
	rootCmd.SuggestionsMinimumDistance = 1
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "kasher [taskName]",
	Short: "kasher - shell task runner with caching",
	Long:  "kasher lets you define, run, and cache named shell tasks.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If a subcommand was called, do nothing here.
		if cmd.CalledAs() == "config" {
			return nil
		}
		// Run a task by name
		if len(args) > 0 {
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			taskName := args[0]
			task, exists := cfg[taskName]
			if !exists {
				return fmt.Errorf("task '%s' not found", taskName)
			}
			// TODO: Add caching logic here
			fmt.Printf("Running: %s\n", task.Command)
			cmd := exec.Command("sh", "-c", task.Command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		}
		return cmd.Help()
	},
}
