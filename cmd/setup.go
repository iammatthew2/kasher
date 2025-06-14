package cmd

import (
	"fmt"
	"kasher/internal/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up a new cached task",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		var action string
		prompt := &survey.Select{
			Message: "What would you like to do?",
			Options: []string{"Create new task", "Update existing task", "Delete task", "Cancel"},
		}
		if err := survey.AskOne(prompt, &action); err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		switch action {
		case "Create new task":
			// Stub task setup
			task := config.TaskConfig{
				Command:    "qdt deployment list",
				Expiration: "1d",
				Notes:      "Query deployments",
			}
			cfg["qdt"] = task

			if err := config.SaveConfig(cfg); err != nil {
				fmt.Println("Failed to save config:", err)
				return
			}
			fmt.Println("Task 'qdt' created.")

		case "Update existing task":
			fmt.Println("TODO: update")
		case "Delete task":
			fmt.Println("TODO: delete")
		case "Cancel":
			fmt.Println("Aborted.")
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
