package main

import (
	"fmt"
	"strconv"
)

// create binary for a command
func aBinary(command aCommand) (string, error) {
	b := "0"
	if _, notInt := strconv.Atoi(command.Data); notInt == nil {
		i, err := strconv.ParseInt(command.Data, 10, 16)
		if err != nil {
			return "", err
		}
		b += fmt.Sprintf("%015b", i) + "\n"
	} else {
		if val, ok := dict[command.Data]; ok {
			i := val
			b += fmt.Sprintf("%015b", i) + "\n"
		} else {
			dict[command.Data] = nxt
			i := nxt
			nxt += 1
			b += fmt.Sprintf("%015b", i) + "\n"
		}
	}
	return b, nil
}

// convert string to binary in the religion of comp
func compToBinary(c cCommand) string {
	m := make(map[string]string)
	m["0"] = "0101010"
	m["1"] = "0111111"
	m["-1"] = "0111010"
	m["D"] = "0001100"
	m["A"] = "0110000"
	m["!D"] = "0001101"
	m["!A"] = "0110001"
	m["-D"] = "0001111"
	m["-A"] = "0110011"
	m["D+1"] = "0011111"
	m["A+1"] = "0110111"
	m["D-1"] = "0001110"
	m["A-1"] = "0110010"
	m["D+A"] = "0000010"
	m["D-A"] = "0010011"
	m["A-D"] = "0000111"
	m["D&A"] = "0000000"
	m["D|A"] = "0010101"
	m["M"] = "1110000"
	m["!M"] = "1110001"
	m["-M"] = "1110011"
	m["M+1"] = "1110111"
	m["M-1"] = "1110010"
	m["D+M"] = "1000010"
	m["D-M"] = "1010011"
	m["M-D"] = "1000111"
	m["D&M"] = "1000000"
	m["D|M"] = "1010101"
	return m[c.Comp]
}

// convert string to binary in the religion of dest
func destToBinary(c cCommand) string {
	m := make(map[string]string)
	m[""] = "000"
	m["M"] = "001"
	m["D"] = "010"
	m["MD"] = "011"
	m["A"] = "100"
	m["AM"] = "101"
	m["AD"] = "110"
	m["AMD"] = "111"
	return m[c.Dest]
}

// convert string to binary in the religion of jump
func jumpToBinary(c cCommand) string {
	m := make(map[string]string)
	m[""] = "000"
	m["JGT"] = "001"
	m["JEQ"] = "010"
	m["JGE"] = "011"
	m["JLT"] = "100"
	m["JNE"] = "101"
	m["JLE"] = "110"
	m["JMP"] = "111"
	return m[c.Jump]
}

// create binary for c command
func cBinary(command cCommand) (string, error) {
	cString := "111"
	for _, c := range compToBinary(command) {
		cString += string(c)
	}
	if len(cString) != 10 {
		return "", fmt.Errorf("error: invalid cCommand.comp to binary")
	}
	for _, c := range destToBinary(command) {
		cString += string(c)
	}
	if len(cString) != 13 {
		return "", fmt.Errorf("error: invalid cCommand.comp to binary")
	}
	for _, c := range jumpToBinary(command) {
		cString += string(c)
	}
	if len(cString) != 16 {
		return "", fmt.Errorf("error: invalid cCommand.comp to binary")
	}
	cString += "\n"
	return cString, nil
}
