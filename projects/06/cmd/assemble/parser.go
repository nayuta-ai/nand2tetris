package main

import (
	"fmt"
	"strings"
)

// remove comment
// example: "@32 // A instruction" -> "@32"
func removeComment(line string) string {
	for i := 1; i < len(line); i++ {
		if line[i-1:i+1] == "//" {
			return strings.TrimSpace(line[:i-1])
		}
	}
	return strings.TrimSpace(line)
}

// a command
type aCommand struct {
	Data string
}

// c command
type cCommand struct {
	Comp string
	Dest string
	Jump string
}

// check the command type
// if line has "@", it is a command
// elif line has "()", it is l command
// else, it is c command
func commandType(line string) (string, error) {
	for _, char := range line {
		if string([]rune{char}) == "@" {
			return aBinary(aParse(line))
		} else if string([]rune{char}) == "(" {
			return "", nil // lParse(line) // l command
		} else {
			return cBinary(cParse(line)) // c command
		}
	}
	return "", fmt.Errorf("unsupported")
}

// a command parser extracts strings after "@"
// example: "@32" -> aCommand.Data = "32"
func aParse(line string) aCommand {
	var val = make([]byte, 0, 100)
	for i, char := range line {
		if i > 0 {
			val = append(val, byte(char))
		}
	}
	return aCommand{
		Data: string(val),
	}
}

// l command parser extracts strings after "(" and before ")"
// example: "(LOOP)" -> lCommand.Data = "LOOP"
func lParse(line string) string {
	var val = make([]byte, 0, 100)
	for i, char := range line {
		if i == 0 && string(char) != "(" {
			return string(val)
		}
		if 0 < i && i < len(line)-1 {
			val = append(val, byte(char))
		}
	}
	return string(val)
}

// c command parser extracts some strings; one called dest is before "=",
// second called comp is within "=" and ";", the other called jump is after ";"
// example: "A=D-A;JEQ" -> cCommand.Comp = "D-A", cCommand.Dest = "A", cCommand.Jump = "JEQ"
func cParse(line string) cCommand {
	var dest = make([]byte, 0, 100)
	var jump = make([]byte, 0, 100)
	var comp = make([]byte, 0, 100)
	var pointer = 0
	for _, char := range line {
		if string([]rune{char}) == "=" {
			pointer = 1
		} else if string([]rune{char}) == ";" {
			pointer = 2
		} else {
			if pointer == 0 {
				dest = append(dest, byte(char))
			} else if pointer == 1 {
				comp = append(comp, byte(char))
			} else {
				jump = append(jump, byte(char))
			}
		}
	}
	if string(comp) == "" {
		comp, dest = dest, comp
	}
	return cCommand{
		Comp: string(comp),
		Dest: string(dest),
		Jump: string(jump),
	}
}

// create dict contained pre symbol
func preDict() {
	dict["SP"] = 0
	dict["LCL"] = 1
	dict["ARG"] = 2
	dict["THIS"] = 3
	dict["THAT"] = 4
	dict["R0"] = 0
	dict["R1"] = 1
	dict["R2"] = 2
	dict["R3"] = 3
	dict["R4"] = 4
	dict["R5"] = 5
	dict["R6"] = 6
	dict["R7"] = 7
	dict["R8"] = 8
	dict["R9"] = 9
	dict["R10"] = 10
	dict["R11"] = 11
	dict["R12"] = 12
	dict["R13"] = 13
	dict["R14"] = 14
	dict["R15"] = 15
	dict["SCREEN"] = 16384
	dict["KBD"] = 24576
}
