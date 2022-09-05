package main

// addVMCommand returns the comment line which is formatted the VM command.
func addVMCommand(command string) string {
	return "//" + command + "\n"
}

// addSpace returns the asm command which is added the new line.
func addSpace(asm_command string) string {
	return asm_command + "\n"
}

// commandPush converts "push" lines in vm file to lines in asm file and returns lines in asm file.
func commandPush(command string, div_command []string) string {
	var asm_command = addVMCommand(command)
	// div_line[0]: "push"
	// div_line[1]: storage location such as "constant"
	// div_line[2]: value
	if div_command[1] == "constant" {
		asm_command += "@" + div_command[2] + "\n" // get value
		asm_command += "D=A\n"
		asm_command += "@SP\n"
		asm_command += "A=M\n"
		asm_command += "M=D\n"
		asm_command += "@SP\n"
		asm_command += "M=M+1\n"
	}
	return asm_command
}

// commandAdd converts "add" lines in vm file to lines in asm file and returns lines in asm file.
func commandAdd(command string) string {
	var asm_command = addVMCommand(command) // Add VM command as the comment
	asm_command += "@SP\n"
	asm_command += "AM=M-1\n"
	asm_command += "D=M\n"
	asm_command += "A=A-1\n"
	asm_command += "D=D+M\n"
	asm_command += "M=D\n"
	return asm_command
}

// commandEnd returns the end statement.
// "(END)\n@END\n0;JMP\n" matches the end statement of asm file.
func commandEnd() string {
	return "(END)\n@END\n0;JMP\n"
}
