package main

import (
	"strconv"
	"strings"
)

var trial = 0 // count the trial of calling commandCompare
var count = 0 // count the number of calling call command
var function = ""

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
	var asm_command = addVMCommand(command) // Add VM command as the comment
	// div_line[0]: "push"
	// div_line[1]: storage location such as "constant"
	// div_line[2]: value
	if div_command[1] == "constant" {
		asm_command += "@" + div_command[2] + "\n" // get value
		asm_command += "D=A\n"
	} else {
		asm_command += commandSymbol(command, div_command)
		asm_command += "D=M\n"
	}
	asm_command += "@SP\n"
	asm_command += "A=M\n"
	asm_command += "M=D\n"
	asm_command += "@SP\n"
	asm_command += "M=M+1\n"
	return asm_command
}

// commandPop converts "pop" lines in vm file to lines in asm file and returns lines in asm file.
func commandPop(command string, div_command []string) string {
	var asm_command = addVMCommand(command)
	asm_command += commandSymbol(command, div_command)
	asm_command += "D=A\n"
	asm_command += "@R13\n"
	asm_command += "M=D\n"
	asm_command += "@SP\n"
	asm_command += "M=M-1\n"
	asm_command += "A=M\n"
	asm_command += "D=M\n"
	asm_command += "@R13\n"
	asm_command += "A=M\n"
	asm_command += "M=D\n"
	return asm_command
}

// commandSymbol converts any symbol such as temp, local, etc. lines in vm file to lines in asm file and returns lines in asm file.
func commandSymbol(command string, div_command []string) string {
	if div_command[1] == "temp" {
		return useTemp(div_command[2])
	} else if div_command[1] == "pointer" {
		return usePointer(div_command[2])
	} else if div_command[1] == "static" {
		return useStatic(div_command[2])
	} else {
		return useSymbol(command, div_command)
	}
}

// useStatic converts "static" lines in vm file to lines in asm file and returns lines in asm file.
func useStatic(val string) string {
	return "@" + function + "." + val + "\n"
}

// usePointer converts "pointer" lines in vm file to lines in asm file and returns lines in asm file.
func usePointer(val string) string {
	if val == "0" {
		return "@THIS\n"
	} else {
		return "@THAT\n"
	}
}

// useTemp converts "temp" lines in vm file to lines in asm file and returns lines in asm file.
func useTemp(index string) string {
	i, _ := strconv.ParseInt(index, 10, 32)
	i += 5
	str := strconv.Itoa(int(i))
	return "@R" + str + "\n"
}

// useSymbol converts "local", "arg", "this", "that", etc. lines in vm file to lines in asm file and returns lines in asm file.
func useSymbol(command string, div_command []string) string {
	var asm_command string
	if div_command[1] == "local" {
		asm_command += "@LCL\n"
	} else if div_command[1] == "argument" {
		asm_command += "@ARG\n"
	} else if div_command[1] == "this" {
		asm_command += "@THIS\n"
	} else if div_command[1] == "that" {
		asm_command += "@THAT\n"
	}
	asm_command += "D=M\n"
	asm_command += "@" + div_command[2] + "\n"
	asm_command += "A=D+A\n"
	return asm_command
}

// commandAdd converts "add" lines in vm file to lines in asm file and returns lines in asm file.
func commandAdd(command string) string {
	return commandCalc(command, "add")
}

// commandSub converts "sub" lines in vm file to lines in asm file and returns lines in asm file.
func commandSub(command string) string {
	return commandCalc(command, "sub")
}

// commandAnd converts "and" lines in vm file to lines in asm file and returns lines in asm file.
func commandAnd(command string) string {
	return commandCalc(command, "and")
}

// commandOr converts "or" lines in vm file to lines in asm file and returns lines in asm file.
func commandOr(command string) string {
	return commandCalc(command, "or")
}

// commandCalc converts lines such as "add", "sub", "and", and "or" in vm file to lines in asm file and returns lines in asm file.
func commandCalc(command string, types string) string {
	var asm_command = addVMCommand(command) // Add VM command as the comment
	asm_command += "@SP\n"
	asm_command += "M=M-1\n"
	asm_command += "A=M\n"
	asm_command += "D=M\n"
	asm_command += "@SP\n"
	asm_command += "M=M-1\n"
	asm_command += "@SP\n"
	asm_command += "A=M\n"
	if types == "add" {
		asm_command += "M=M+D\n"
	} else if types == "sub" {
		asm_command += "M=M-D\n"
	} else if types == "and" {
		asm_command += "M=M&D\n"
	} else if types == "or" {
		asm_command += "M=M|D\n"
	}
	asm_command += "@SP\n"
	asm_command += "M=M+1\n"
	return asm_command
}

// commandNeg converts "neg" lines in vm file to lines in asm file and returns lines in asm file.
func commandNeg(command string) string {
	return commandNegNot(command, "neg")
}

// commandNot converts "not" lines in vm file to lines in asm file and returns lines in asm file.
func commandNot(command string) string {
	return commandNegNot(command, "not")
}

// commandNegNot converts "neg" or "not" lines in vm file to lines in asm file and returns lines in asm file.
func commandNegNot(command string, types string) string {
	var asm_command = addVMCommand(command) // Add VM command as the comment
	asm_command += "@SP\n"
	asm_command += "M=M-1\n"
	asm_command += "A=M\n"
	asm_command += "D=M\n"
	if types == "neg" {
		asm_command += "M=-D\n"
	} else if types == "not" {
		asm_command += "M=!D\n"
	}
	asm_command += "@SP\n"
	asm_command += "M=M+1\n"
	return asm_command
}

// commandEq converts "eq" lines in vm file to lines in asm file and returns lines in asm file.
func commandEq(command string) string {
	return commandCompare(command, "eq")
}

// commandLt converts "lt" lines in vm file to lines in asm file and returns lines in asm file.
func commandLt(command string) string {
	return commandCompare(command, "lt")
}

// commandGt converts "gt" lines in vm file to lines in asm file and returns lines in asm file.
func commandGt(command string) string {
	return commandCompare(command, "gt")
}

// commandCompare converts lines such as "eq", "lt", and "gt" in vm file to lines in asm file and returns lines in asm file.
func commandCompare(command string, types string) string {
	var asm_command = addVMCommand(command) // Add VM command as the comment
	trial_str := strconv.Itoa(trial)        // Convert trial into string

	asm_command += "@SP\n"
	asm_command += "M=M-1\n"
	asm_command += "A=M\n"
	asm_command += "D=M\n"
	asm_command += "@SP\n"
	asm_command += "M=M-1\n"
	asm_command += "@SP\n"
	asm_command += "A=M\n"
	asm_command += "D=M-D\n"
	asm_command += "@BOOL" + trial_str + "\n"
	if types == "eq" {
		asm_command += "D;JEQ\n"
	} else if types == "lt" {
		asm_command += "D;JLT\n"
	} else if types == "gt" {
		asm_command += "D;JGT\n"
	}
	asm_command += "@SP\n"
	asm_command += "A=M\n"
	asm_command += "M=0\n"
	asm_command += "@ENDBOOL" + trial_str + "\n"
	asm_command += "0;JMP\n"
	asm_command += "(BOOL" + trial_str + ")\n"
	asm_command += "@SP\n"
	asm_command += "A=M\n"
	asm_command += "M=-1\n"
	asm_command += "(ENDBOOL" + trial_str + ")\n"
	asm_command += "@SP\n"
	asm_command += "M=M+1\n"
	trial += 1
	return asm_command
}

func commandLabel(command string) string {
	return "(" + command + ")" + "\n"
}

func commandIfGoto(command string) string {
	return "@SP\nM=M-1\nA=M\nD=M\n" + "@" + command + "\n" + "D;JNE\n"
}

func commandGoto(command string) string {
	return "@" + command + "\n" + "0;JMP\n"
}

func initialPush(num string) string {
	var asm_command string
	asm_command += "D=" + num + "\n"
	asm_command += "@SP\n"
	asm_command += "A=M\n"
	asm_command += "M=D\n"
	asm_command += "@SP\n"
	asm_command += "M=M+1\n"
	return asm_command
}

func commandCall(function_name string, n_arg int) string {
	var asm_command = addVMCommand("call " + function_name + strconv.Itoa(n_arg))
	label_name := function_name + strconv.Itoa(count)
	asm_command += "@" + label_name + "\n"
	// push return-address
	asm_command += initialPush("A")
	// push LCL
	asm_command += "@LCL\n"
	asm_command += initialPush("M")
	// push ARG
	asm_command += "@ARG\n"
	asm_command += initialPush("M")
	// push THIS
	asm_command += "@THIS\n"
	asm_command += initialPush("M")
	// push THAT
	asm_command += "@THAT\n"
	asm_command += initialPush("M")
	// LCL = SP
	asm_command += "@SP\n"
	asm_command += "D=M\n"
	asm_command += "@LCL\n"
	asm_command += "M=D\n"
	// ARG = SP-n-5
	asm_command += "@SP\n"
	asm_command += "D=M\n"
	asm_command += "@" + strconv.Itoa(n_arg+5) + "\n"
	asm_command += "D=D-A\n"
	asm_command += "@ARG\n"
	asm_command += "M=D\n"
	// goto f
	asm_command += commandGoto(function_name)
	// (return-address)
	asm_command += commandLabel(label_name)
	count += 1
	return asm_command
}

func commandInit() string {
	var asm_command = ""
	asm_command += "@256\n"
	asm_command += "D=A\n"
	asm_command += "@SP\n"
	asm_command += "M=D\n"
	asm_command += commandCall("Sys.init", 0)
	return asm_command
}

func commandFunction(command string, div_command []string) (string, error) {
	var asm_command = addVMCommand(command)
	function = strings.Split(div_command[1], ".")[0]
	asm_command += "(" + div_command[1] + ")\n"
	repeat_time, err := strconv.Atoi(div_command[2])
	if err != nil {
		return "", err
	}
	for i := 0; i < repeat_time; i++ {
		asm_command += initialPush("0")
	}
	return asm_command, nil
}

func restoreAddress(pointer string, num string) string {
	var asm_command string
	asm_command += "@R13\n"
	asm_command += "D=M\n"
	asm_command += "@" + num + "\n"
	asm_command += "D=D-A\n"
	asm_command += "A=D\n"
	asm_command += "D=M\n"
	asm_command += "@" + pointer + "\n"
	asm_command += "M=D\n"
	return asm_command
}

func commandReturn(command string) string {
	var asm_command = addVMCommand(command)
	// FRAME = LCL
	asm_command += "@LCL\n"
	asm_command += "D=M\n"
	// RET = *(FRAME-5)
	asm_command += "@R13\n"
	asm_command += "M=D\n"
	asm_command += "@R13\n"
	asm_command += "D=M\n"
	asm_command += "@5\n"
	asm_command += "D=D-A\n"
	asm_command += "A=D\n"
	asm_command += "D=M\n"
	asm_command += "@R14\n"
	asm_command += "M=D\n"
	// *ARG = pop()
	asm_command += "@SP\n"
	asm_command += "M=M-1\n"
	asm_command += "A=M\n"
	asm_command += "D=M\n"
	asm_command += "@ARG\n"
	asm_command += "A=M\n"
	asm_command += "M=D\n"
	// SP = ARG+1
	asm_command += "@ARG\n"
	asm_command += "D=M\n"
	asm_command += "@SP\n"
	asm_command += "M=D+1\n"
	// THAT = *(FRAME-1)
	asm_command += restoreAddress("THAT", "1")
	// THIS = *(FRAME-2)
	asm_command += restoreAddress("THIS", "2")
	// ARG = *(FRAME-3)
	asm_command += restoreAddress("ARG", "3")
	// ARG = *(FRAME-4)
	asm_command += restoreAddress("LCL", "4")
	// goto RET
	asm_command += "@R14\n"
	asm_command += "A=M\n"
	asm_command += "0;JMP\n"
	return asm_command
}

// commandEnd returns the end statement.
// "(END)\n@END\n0;JMP\n" matches the end statement of asm file.
func commandEnd() string {
	return "(END)\n@END\n0;JMP\n"
}
