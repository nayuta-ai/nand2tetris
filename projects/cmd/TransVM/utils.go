package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// parseArgs returns the file path as the string type.
func parseArgs() (string, string, error) {
	fp := flag.String("filepath", "", "filepath")
	fn := flag.String("filename", "", "filename")
	flag.Parse()
	if *fn == "" {
		return "", "", fmt.Errorf("error: filename is empty")
	}
	if *fp == "" {
		return "", "", fmt.Errorf("error: filepath doesn't exist")
	}
	return fmt.Sprintf(*fp), fmt.Sprintf(*fn), nil
}

func checkSys(files []string) bool {
	for _, file := range files {
		file_split := strings.Split(file, "/")
		if file_split[len(file_split)-1] == "Sys.vm" {
			return true
		}
	}
	return false
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
