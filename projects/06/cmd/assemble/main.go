package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var dict = make(map[string]int)
var nxt = 16

func main() {
	f := flag.String("filename", "", "filename")
	flag.Parse()
	if *f == "" {
		fmt.Println("Error: filename is empty")
		return
	}
	filepath := fmt.Sprintf("test/%s", *f)
	fp, err := os.Open(filepath + ".asm") // To be resolved later with argparse
	if err != nil {
		log.Println(err)
		return
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	lines := make([]string, 0)
	count := 0
	preDict()
	for scanner.Scan() {
		line := RemoveComment(scanner.Text()) // remove comment
		if len(line) != 0 {
			char := lParse(line)
			if char != "" {
				dict[char] = count
				count -= 1
			}
			lines = append(lines, line)
			count += 1
		}
	}
	if err = scanner.Err(); err != nil {
		log.Println(err)
	}
	byteInstraction := make([][]byte, 0)
	for _, line := range lines {
		instraction := CommandType(line)
		byteInstraction = append(byteInstraction, instraction)
	}
	err = createHack(byteInstraction, filepath)
	if err != nil {
		log.Println(err)
	}
}
