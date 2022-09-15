package main

import (
	"log"
	"os"
)

var dict = make(map[string]int) // dictionary for acommand, lcommand
var nxt = 16                    // address for undefined variable

// read the file and write the binary represents to *.hack
func main() {
	filepath, err := parseArgs(os.Args[1:]...)
	if err != nil {
		log.Println(err)
	}
	fp, err := os.Open(filepath + ".asm") // To be resolved later with argparse
	if err != nil {
		log.Println(err)
		return
	}
	defer fp.Close()

	byteInstraction, err := scanline(fp)
	if err != nil {
		log.Println(err)
	}
	err = createHack(byteInstraction, filepath)
	if err != nil {
		log.Println(err)
	}
}
