/*
Package main implements the VM transrator.

This transrator assign the vm file, convert vm file to asm file, and output it.
*/
package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

// main returns .asm file which is converted .vm file into.
func main() {
	path, filename, err := parseArgs()
	if err != nil {
		log.Println(err)
	}
	files, err := filepath.Glob(path + "/*.vm") // Get the .vm file
	if err != nil {
		log.Println(err)
		return
	}
	commands := "" // Initialize the content of the .asm file

	if checkSys(files) {
		commands += commandInit()
	}
	for _, file := range files {
		fp, _ := os.Open(file)
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
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
	}
	commands += commandEnd()                     // Add the end statement
	err = createAsm(commands, path+"/"+filename) // Create .asm file
	if err != nil {
		log.Println(err)
		return
	}
}
