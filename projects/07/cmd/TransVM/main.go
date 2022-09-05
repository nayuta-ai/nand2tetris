/*
Package main implements the VM transrator.

This transrator assign the vm file, convert vm file to asm file, and output it.
*/
package main

import (
	"bufio"
	"log"
	"os"
)

// main returns .asm file which is converted .vm file into.
func main() {
	filepath, err := parseArgs()
	if err != nil {
		log.Println(err)
	}
	fp, err := os.Open(filepath + ".vm") // Get the .vm file
	if err != nil {
		log.Println(err)
		return
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	commands := "" // Initialize the content of the .asm file
	for scanner.Scan() {
		command := parser(scanner.Text()) // Parse command
		if command != "" {
			commands += addSpace(command) // Add the blank line
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
		return
	}
	commands += commandEnd()            // Add the end statement
	err = createAsm(commands, filepath) // Create .asm file
	if err != nil {
		log.Println(err)
		return
	}
}
