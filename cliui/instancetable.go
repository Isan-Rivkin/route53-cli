package cliui

import (
	"fmt"
)

// SelectInstanceFromList is return selected instance from the list by prompt question
func SelectInstanceFromList(instances []string) (string, error) {

	if len(instances) == 1 {
		return instances[0], nil
	}
	tablePrompt := NewTable()

	header := []string{
		"#",
		"Instance ID",
		"IP",
		"Role",
		"Environment",
		"Instance Type",
		"Platform",
		"Region",
	}

	tablePrompt.AddHeaders(header)

	for i, row := range instances {
		tablePrompt.AddRow(i+1, 0, fmt.Sprintf("%d", i+1))
		tablePrompt.AddRow(i+1, 1, row)
		tablePrompt.AddRow(i+1, 2, row)
		tablePrompt.AddRow(i+1, 3, row)
		tablePrompt.AddRow(i+1, 4, row)
		tablePrompt.AddRow(i+1, 5, row)
		tablePrompt.AddRow(i+1, 6, row)
		tablePrompt.AddRow(i+1, 7, row)

	}

	selectedInstance, err := tablePrompt.Render()
	if err != nil {
		return "", err
	}

	return instances[selectedInstance-1], nil
}
