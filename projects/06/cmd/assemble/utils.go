package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// parse the argument of filename and create filepath variable
func parseArgs() (string, error) {
	f := flag.String("filename", "", "filename")
	flag.Parse()
	if *f == "" {
		return "", fmt.Errorf("error: filename is empty")
	}
	return fmt.Sprintf("test/%s", *f), nil
}

// scan the content in filepath and convert to binary string each line
func scanline(fp *os.File) ([]string, error) {
	scanner := bufio.NewScanner(fp)
	lines := make([]string, 0)
	count := 0
	preDict()
	for scanner.Scan() {
		line := removeComment(scanner.Text()) // remove comment
		if len(line) != 0 {
			char := lParse(line)
			if char != "" {
				dict[char] = count
				count -= 1
			}
			lines = append(lines, line)
			count += 1 // count the row
		}
	}
	if err := scanner.Err(); err != nil {
		return make([]string, 0), err
	}
	byteInstraction := make([]string, 0)
	for _, line := range lines {
		instraction, err := commandType(line)
		if err != nil {
			return make([]string, 0), err
		}
		byteInstraction = append(byteInstraction, instraction)
	}
	return byteInstraction, nil
}

// create Hack file converting from [][]byte
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
