package main

import (
	"flag"
	"fmt"
	"os"
)

// parseArgs returns the file path as the string type.
func parseArgs() (string, error) {
	f := flag.String("filename", "", "filename")
	flag.Parse()
	if *f == "" {
		return "", fmt.Errorf("error: filename is empty")
	}
	return fmt.Sprintf(*f), nil
}

// createAsm returns .asm file as the stdoutput.
func createAsm(commands string, filepath string) error {
	file, _ := os.Create(filepath + ".asm")
	defer file.Close()
	_, err := file.WriteString(commands)
	if err != nil {
		return err
	}
	return nil
}
