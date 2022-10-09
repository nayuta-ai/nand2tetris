package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// This function parses the argument of filename and create filepath variable
func parseArgs(args ...string) (string, error) {
	flg := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f := flg.String("filename", "", "filename")
	flg.Parse(args)
	if *f == "" {
		return "", fmt.Errorf("error: filename is empty")
	}
	return fmt.Sprintf("test/%s", *f), nil
}

// This function scans the contents in filepath and converts to binary string each line
func scanline(fp *os.File) ([]string, error) {
	scanner := bufio.NewScanner(fp)
	lines := make([]string, 0)
	count := 0
	preDict()
	// The first loop removes comments and creates a label table
	for scanner.Scan() {
		line := removeComment(scanner.Text()) // It removes comments
		if len(line) != 0 {
			char := lParser(line)
			if char != "" {
				dict[char] = count
				count -= 1
			}
			lines = append(lines, line)
			count += 1 // It counts the row
		}
	}
	if err := scanner.Err(); err != nil {
		return make([]string, 0), err
	}
	byteInstraction := make([]string, 0)
	// The second loop converts the assembler's command
	for _, line := range lines {
		instraction, err := convertAsmCommand(line)
		if err != nil {
			return make([]string, 0), err
		}
		byteInstraction = append(byteInstraction, instraction)
	}
	return byteInstraction, nil
}

// This function creates Hack file converting from [][]byte
func createHack(dataset []string, filepath string) error {
	file, _ := os.Create(filepath + ".hack")
	defer file.Close()
	for _, data := range dataset {
		_, err := file.WriteString(data)
		if err != nil {
			return err
		}
	}
	return nil
}
