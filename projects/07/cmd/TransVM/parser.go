package main

import (
	"strings"
)

// parser returns the string of parse a line in .vm file.
func parser(command string) string {
	div_command := divideCommand(removeComment(command)) // Divide the command by the space
	if len(div_command) == 0 {
		return ""
	} else if len(div_command) == 1 {
		if div_command[0] == "add" {
			return commandAdd(command)
		}
	} else if len(div_command) == 3 {
		if div_command[0] == "push" {
			return commandPush(command, div_command)
		}
	}
	return ""
}

// divideCommand returns the list of the command.
// For example, command = "push constant 7" returns ["push", "constant", "7"]
func divideCommand(command string) []string {
	return strings.Fields(command)
}

// removeComment returns the string which of comment line is removed.
// For example, command = "// push the command" returns ""
func removeComment(command string) string {
	for i := 1; i < len(command); i++ {
		// Check if the comment line is existed in command
		if command[i-1:i+1] == "//" {
			return strings.TrimSpace(command[:i-1]) // Return the string before "//" and removes the space.
		}
	}
	return strings.TrimSpace(command) // Remove the space.
}
