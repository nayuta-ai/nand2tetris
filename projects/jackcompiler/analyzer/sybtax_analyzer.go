package analyzer

import (
	"bufio"
	"os"
)

func SyntaxAnalyzer(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	token, err := tokenizer(scanner)
	if err != nil {
		return err
	}
	tree, err := parser(token)
	if err != nil {
		return err
	}
	file := fetchFileName(path)
	err = convertToXMLformat(tree, file[:len(file)-5])
	// err = convertToTXMLformat(token)
	if err != nil {
		return err
	}
	return nil
}
