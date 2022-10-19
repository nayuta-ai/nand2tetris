package main

import (
	"strconv"
	"strings"
)

// parser returns the string of parse a line in .vm file.
func parser(command string) string {
	remove_comment := removeComment(command)
	div_command := divideCommand(remove_comment) // Divide the command by the space
	if len(div_command) == 0 {
		return ""
	} else if len(div_command) == 1 {
		return callCommand(div_command[0])
	} else if len(div_command) == 2 {
		return callLabel(remove_comment, div_command)
	} else if len(div_command) == 3 {
		if div_command[0] == "push" {
			return commandPush(remove_comment, div_command)
		} else if div_command[0] == "pop" {
			return commandPop(remove_comment, div_command)
		} else if div_command[0] == "function" {
			asm_command, _ := commandFunction(remove_comment, div_command)
			return asm_command
		} else if div_command[0] == "call" {
			function_name := div_command[1]
			n_arg, _ := strconv.Atoi(div_command[2])
			return commandCall(function_name, n_arg)
		}
	}
	return ""
}

func callLabel(command string, div_command []string) string {
	var asm_command = addVMCommand(command)
	if div_command[0] == "label" {
		asm_command += commandLabel(div_command[1])
	} else if div_command[0] == "goto" {
		asm_command += commandGoto(div_command[1])
	} else if div_command[0] == "if-goto" {
		asm_command += commandIfGoto(div_command[1])
	}
	return asm_command
}

// callCommand returns the string of parse a line in .vm file such as
// add, sub, neg, eq, lt, gt, and, or, not
func callCommand(command string) string {
	if command == "add" {
		return commandAdd(command)
	} else if command == "sub" {
		return commandSub(command)
	} else if command == "neg" {
		return commandNeg(command)
	} else if command == "eq" {
		return commandEq(command)
	} else if command == "lt" {
		return commandLt(command)
	} else if command == "gt" {
		return commandGt(command)
	} else if command == "and" {
		return commandAnd(command)
	} else if command == "or" {
		return commandOr(command)
	} else if command == "not" {
		return commandNot(command)
	} else if command == "return" {
		return commandReturn(command)
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
