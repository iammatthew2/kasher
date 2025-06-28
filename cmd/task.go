package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"kasher/internal/config"

	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage kasher tasks",
}

func init() {
	taskCmd.AddCommand(createCmd)
	taskCmd.AddCommand(createForCmd)
	taskCmd.AddCommand(updateCmd)
	taskCmd.AddCommand(deleteCmd)
	taskCmd.AddCommand(listCmd)
	taskCmd.AddCommand(clearAllCmd)

	taskCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show extra information")
}

// getCacheDir returns the kasher cache directory path
func getCacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "kasher"), nil
}

func printVerboseInfo() {
	configPath, err := config.GetConfigPath()
	if err == nil {
		fmt.Println("Kasher config file location:")
		fmt.Println(configPath)
	}
	cacheDir, err := getCacheDir()
	if err == nil {
		fmt.Println("Kasher cache directory:")
		fmt.Println(cacheDir)
	}
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new task",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			printVerboseInfo()
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		fmt.Println("** Type '?' for help, 'q' to quit at any prompt **")
		var taskName string
		if len(args) > 0 {
			taskName = args[0]
			if isReservedTaskName(taskName) {
				return fmt.Errorf("the name '%s' is reserved and cannot be used. Please choose another name", taskName)
			}
			if _, exists := cfg[taskName]; exists {
				return fmt.Errorf("task '%s' already exists", taskName)
			}
			// Optionally, validate/normalize taskName as PromptForTaskName would
		} else {
			var err error
			taskName, err = PromptForTaskName(cfg, "Enter a name for the new task:")
			if err != nil {
				return err
			}
			if _, exists := cfg[taskName]; exists {
				return fmt.Errorf("task '%s' already exists", taskName)
			}
		}
		task, err := PromptTaskDetails(&config.TaskConfig{}, false)
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
		if verbose {
			printVerboseInfo()
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		taskName, err := PromptTaskName(cfg, "Select a task to update:")
		if err != nil {
			return err
		}
		existing := cfg[taskName]
		task, err := PromptTaskDetails(&existing, false)
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
		if verbose {
			printVerboseInfo()
		}
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
		if verbose {
			printVerboseInfo()
		}
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
		if verbose {
			printVerboseInfo()
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		if len(cfg) == 0 {
			fmt.Println("No tasks found.")
			return nil
		}
		fmt.Println("Tasks:")

		// Get sorted task names for consistent ordering
		var names []string
		for name := range cfg {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			task := cfg[name]
			fmt.Printf("- %s: %s (expires: %s)\n", name, task.Command, task.Expiration)
			if task.Notes != "" {
				fmt.Printf("    Notes: %s\n", task.Notes)
			}
		}
		return nil
	},
}

var createForCmd = &cobra.Command{
	Use:   "createFor <command>",
	Short: "Create a new task for a given shell command",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			printVerboseInfo()
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		taskName, err := PromptForTaskName(cfg, "Enter a name for the new task:")
		if err != nil {
			return err
		}
		if _, exists := cfg[taskName]; exists {
			return fmt.Errorf("task '%s' already exists", taskName)
		}
		if len(args) > 1 {
			fmt.Println("Warning: It looks like you passed multiple arguments. If your command contains spaces, pipes, or shell operators (like &&), you should quote the command, e.g.:\n  kasher task createFor \"echo starting && sleep 5 && echo ending\"")
		}
		// Join all args as the shell command
		shellCommand := strings.Join(args, " ")
		task, err := PromptTaskDetails(&config.TaskConfig{Command: shellCommand}, true)
		if err != nil {
			return err
		}
		if err := cfg.AddTask(taskName, task); err != nil {
			return err
		}
		if err := config.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Printf("Task '%s' created for command: %s\n", taskName, shellCommand)
		return nil
	},
}
