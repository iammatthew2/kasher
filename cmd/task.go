package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"kasher/internal/config"

	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage kasher tasks",
}

func init() {
	taskCmd.AddCommand(createCmd)
	taskCmd.AddCommand(updateCmd)
	taskCmd.AddCommand(deleteCmd)
	taskCmd.AddCommand(listCmd)
	taskCmd.AddCommand(clearAllCmd)
	taskCmd.AddCommand(infoCmd)
}

// promptForTaskName prompts for a new task name, ensuring it is not empty and not already in use.
func promptForTaskName(cfg config.KasherConfig, message string) (string, error) {
	for {
		var name string
		fmt.Print(message + " ")
		fmt.Scanln(&name)
		if name == "" {
			fmt.Println("Task name cannot be empty.")
			continue
		}
		if _, exists := cfg[name]; exists {
			fmt.Printf("Task '%s' already exists. Please choose another name.\n", name)
			continue
		}
		return name, nil
	}
}

var infoCmd = &cobra.Command{
	Use:   "_info",
	Short: "Print the location of the kasher config and cache directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := config.GetConfigPath()
		if err != nil {
			return err
		}
		fmt.Println("Kasher config file location:")
		fmt.Println(configPath)

		// Print cache directory
		cacheDir, err := getCacheDir()
		if err != nil {
			return err
		}
		fmt.Println("Kasher cache directory:")
		fmt.Println(cacheDir)
		return nil
	},
}

// getCacheDir returns the kasher cache directory path
func getCacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "kasher"), nil
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		taskName, err := promptForTaskName(cfg, "Enter a name for the new task:")
		if err != nil {
			return err
		}
		if _, exists := cfg[taskName]; exists {
			return fmt.Errorf("task '%s' already exists", taskName)
		}
		task, err := PromptTaskDetails(&config.TaskConfig{})
		if err != nil {
			return err
		}
		if err := cfg.AddTask(taskName, task); err != nil {
			return err
		}
		if err := config.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Printf("Task '%s' created.\n", taskName)
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing task",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		taskName, err := PromptTaskName(cfg, "Select a task to update:")
		if err != nil {
			return err
		}
		existing := cfg[taskName]
		task, err := PromptTaskDetails(&existing)
		if err != nil {
			return err
		}
		if err := cfg.UpdateTask(taskName, task); err != nil {
			return err
		}
		if err := config.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Printf("Task '%s' updated.\n", taskName)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		taskName, err := PromptTaskName(cfg, "Select a task to delete:")
		if err != nil {
			return err
		}
		if err := cfg.DeleteTask(taskName); err != nil {
			return err
		}
		if err := config.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Printf("Task '%s' deleted.\n", taskName)
		return nil
	},
}

var clearAllCmd = &cobra.Command{
	Use:   "clearAll",
	Short: "Delete all tasks settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		var confirm string
		fmt.Print("Are you sure you want to delete all tasks and clear the config? (y/N): ")
		fmt.Scanln(&confirm)
		if confirm == "y" || confirm == "Y" {
			if err := config.ClearConfig(); err != nil {
				return fmt.Errorf("failed to clear config: %w", err)
			}
			fmt.Println("Config cleared.")
		} else {
			fmt.Println("Aborted.")
		}
		return nil
	},
}
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		if len(cfg) == 0 {
			fmt.Println("No tasks found.")
			return nil
		}
		fmt.Println("Tasks:")
		for name, task := range cfg {
			fmt.Printf("- %s: %s (expires: %s)\n", name, task.Command, task.Expiration)
			if task.Notes != "" {
				fmt.Printf("    Notes: %s\n", task.Notes)
			}
		}
		return nil
	},
}
