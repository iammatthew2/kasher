package cmd

import (
	"fmt"
	"strings"
	"time"

	"kasher/internal/config"

	"github.com/AlecAivazis/survey/v2"
)

// PromptTaskDetails interactively asks the user for task details.
func PromptTaskDetails(existing *config.TaskConfig, skipCommandPrompt bool) (config.TaskConfig, error) {
	var task config.TaskConfig
	if skipCommandPrompt && existing.Command != "" {
		task.Command = existing.Command
	} else {
		commandPrompt := &survey.Input{Message: "Shell command:", Default: existing.Command}
		var command string
		survey.AskOne(commandPrompt, &command, survey.WithValidator(survey.Required))
		if command == "q" || command == "quit" || command == "exit" || command == "?" || command == "help" {
			if command == "?" || command == "help" {
				fmt.Println("Enter the shell command to run for this task. Example: echo 'Hello, world!'")
				return task, fmt.Errorf("user requested help")
			}
			return task, fmt.Errorf("user exited prompt")
		}
		task.Command = command
	}
	// Prompt for expiration with help and validation loop
	for {
		var expiration string
		defaultExpiration := existing.Expiration
		if defaultExpiration == "" {
			defaultExpiration = "24h"
		}
		prompt := &survey.Input{
			Message: "Cache expiration (e.g. 10m, 1h, 2h30m, 45s, 1.5h, 500ms):",
			Default: defaultExpiration,
		}
		survey.AskOne(prompt, &expiration)
		if expiration == "q" || expiration == "quit" || expiration == "exit" || expiration == "?" || expiration == "help" {
			if expiration == "?" || expiration == "help" {
				fmt.Println("Examples of valid durations:")
				fmt.Println("  10m      (10 minutes)")
				fmt.Println("  1h       (1 hour)")
				fmt.Println("  2h30m    (2 hours, 30 minutes)")
				fmt.Println("  45s      (45 seconds)")
				fmt.Println("  1.5h     (1.5 hours)")
				fmt.Println("  500ms    (500 milliseconds)")
				fmt.Println("See https://pkg.go.dev/time#ParseDuration for all valid formats.")
				continue
			}
			return task, fmt.Errorf("user exited prompt")
		}
		if expiration == "" {
			fmt.Println("Expiration is required.")
			continue
		}
		if _, err := time.ParseDuration(expiration); err != nil {
			fmt.Printf("Invalid duration: %v\n", err)
			continue
		}
		task.Expiration = expiration
		break
	}
	// Prompt for notes
	notesPrompt := &survey.Input{Message: "Notes (optional):", Default: existing.Notes}
	var notes string
	survey.AskOne(notesPrompt, &notes)
	if notes == "q" || notes == "quit" || notes == "exit" || notes == "?" || notes == "help" {
		if notes == "?" || notes == "help" {
			fmt.Println("You can add any notes or description for this task. Leave blank if not needed.")
			return task, fmt.Errorf("user requested help")
		}
		return task, fmt.Errorf("user exited prompt")
	}
	task.Notes = notes
	return task, nil
}

// PromptTaskName prompts the user to select a task name from the config.
func PromptTaskName(cfg config.KasherConfig, message string) (string, error) {
	if len(cfg) == 0 {
		return "", fmt.Errorf("no tasks available")
	}
	var names []string
	for name := range cfg {
		names = append(names, name)
	}
	var selected string
	prompt := &survey.Select{
		Message: message,
		Options: names,
	}
	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", err
	}
	return selected, nil
}

// reservedTaskNames contains task names that are reserved and cannot be used by the user.
var reservedTaskNames = map[string]struct{}{
	"task": {},
	"quit": {},
	"q":    {},
	"exit": {},
	"?":    {},
	"help": {},
}

// isReservedTaskName checks if a given name is reserved.
func isReservedTaskName(name string) bool {
	_, exists := reservedTaskNames[strings.ToLower(name)]
	return exists
}

// PromptForTaskName interactively prompts for a new task name, ensuring it is not empty and not already in use.
func PromptForTaskName(cfg config.KasherConfig, message string) (string, error) {
	for {
		var name string
		prompt := &survey.Input{Message: message}
		err := survey.AskOne(prompt, &name, survey.WithValidator(survey.Required))
		if err != nil {
			return "", err
		}
		if name == "q" || name == "quit" || name == "exit" {
			return "", fmt.Errorf("user exited prompt")
		}
		if name == "?" || name == "help" {
			fmt.Println("Enter a unique name for your new task. The name must not be empty, reserved, or already in use. Spaces will be replaced with dashes.")
			continue
		}
		if name == "" {
			fmt.Println("Task name cannot be empty.")
			continue
		}
		if isReservedTaskName(name) {
			fmt.Println("That name is reserved and cannot be used. Please choose another name.")
			continue
		}
		// Auto-replace spaces with dashes
		if containsSpace := (len(name) != len(removeSpaces(name))); containsSpace {
			fmt.Printf("Note: Spaces in task name will be replaced with dashes: '%s'\n", replaceSpacesWithDashes(name))
			name = replaceSpacesWithDashes(name)
		}
		if _, exists := cfg[name]; exists {
			fmt.Printf("Task '%s' already exists. Please choose another name.\n", name)
			continue
		}
		return name, nil
	}
}

// replaceSpacesWithDashes replaces all spaces in a string with dashes
func replaceSpacesWithDashes(s string) string {
	result := ""
	for _, r := range s {
		if r == ' ' {
			result += "-"
		} else {
			result += string(r)
		}
	}
	return result
}

// removeSpaces removes all spaces from a string
func removeSpaces(s string) string {
	result := ""
	for _, r := range s {
		if r != ' ' {
			result += string(r)
		}
	}
	return result
}
