package cmd

import (
	"fmt"

	"kasher/internal/config"

	"github.com/AlecAivazis/survey/v2"
)

// PromptTaskDetails interactively asks the user for task details.
func PromptTaskDetails(existing *config.TaskConfig) (config.TaskConfig, error) {
	var task config.TaskConfig
	qs := []*survey.Question{
		{
			Name:     "command",
			Prompt:   &survey.Input{Message: "Shell command:", Default: existing.Command},
			Validate: survey.Required,
		},
		{
			Name:     "expiration",
			Prompt:   &survey.Input{Message: "Cache expiration (e.g. 10m, 1h):", Default: existing.Expiration},
			Validate: survey.Required,
		},
		{
			Name:   "notes",
			Prompt: &survey.Input{Message: "Notes (optional):", Default: existing.Notes},
		},
	}
	answers := struct {
		Command    string
		Expiration string
		Notes      string
	}{}
	if err := survey.Ask(qs, &answers); err != nil {
		return task, err
	}
	task.Command = answers.Command
	task.Expiration = answers.Expiration
	task.Notes = answers.Notes
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
