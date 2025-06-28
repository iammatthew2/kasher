package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"kasher/internal/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var forceRefresh bool
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "kasher [taskName]",
	Short: "kasher - shell task runner with caching",
	Long:  "kasher lets you define, run, and cache named shell tasks.",
	Args:  cobra.ArbitraryArgs, // Accept any arguments (task names)
	RunE: func(cmd *cobra.Command, args []string) error {
		// If a subcommand was called, exit and let subcommand's handler run
		if cmd.CalledAs() == "task" {
			return nil
		}
		// If no args, prompt user to select a task interactively
		if len(args) == 0 {
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if len(cfg) == 0 {
				fmt.Println("No tasks found. Use 'kasher task create' to add one.")
				return nil
			}
			var names []string
			for name := range cfg {
				names = append(names, name)
			}
			sort.Strings(names)
			var selected string
			prompt := &survey.Select{
				Message:  "Select a task to run:",
				Options:  names,
				PageSize: 10,
			}
			if err := survey.AskOne(prompt, &selected); err != nil {
				return err
			}
			args = []string{selected}
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
				fmt.Fprintf(os.Stderr, "Task '%s' not found.\n", taskName)
				fmt.Fprintln(os.Stderr, "Run 'kasher task list' to see available tasks.")
				return nil
			}

			if verbose {
				configPath, err := config.GetConfigPath()
				if err == nil {
					fmt.Println("Kasher config file location:")
					fmt.Println(configPath)
				}
				dir, err := os.UserCacheDir()
				if err == nil {
					fmt.Println("Kasher cache directory:")
					fmt.Println(filepath.Join(dir, "kasher"))
				}
			}

			// Check cache validity, skip if forceRefresh is set
			cacheValid := false
			if !forceRefresh && task.LastFetched != "" && task.Expiration != "" {
				lastFetched, err := time.Parse(time.RFC3339, task.LastFetched)
				if err == nil {
					expDur, err := time.ParseDuration(task.Expiration)
					if err == nil && time.Since(lastFetched) < expDur {
						cacheValid = true
					}
				}
			}

			if cacheValid {
				cached, err := config.ReadCache(taskName)
				if err == nil {
					fmt.Print(cached)
					return nil
				}
				// If cache read fails, fall through to re-run the command
			}

			// Prepare to capture and stream output
			var outBuf, errBuf bytes.Buffer
			command := exec.Command("sh", "-c", task.Command)
			command.Stdout = io.MultiWriter(os.Stdout, &outBuf)
			command.Stderr = io.MultiWriter(os.Stderr, &errBuf)
			command.Stdin = os.Stdin

			fmt.Printf("Running: %s\n", task.Command)
			err = command.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
			}

			// Combine output for caching
			output := outBuf.String() + errBuf.String()

			// Save output to cache file
			_ = config.WriteCache(taskName, output)

			// Update LastFetched and save config
			task.LastFetched = time.Now().Format(time.RFC3339)
			cfg[taskName] = task
			_ = config.SaveConfig(cfg) // handle error as needed

			return err
		}
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(taskCmd)
	rootCmd.SuggestionsMinimumDistance = 2
	rootCmd.PersistentFlags().BoolVarP(&forceRefresh, "force", "f", false, "Force refresh of cached task output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show extra information")
}
