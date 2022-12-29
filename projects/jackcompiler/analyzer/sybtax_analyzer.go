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
	token, err := Tokenizer(scanner)
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

func GetTokens(path string) (*Tokens, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	tokens, err := Tokenizer(scanner)
	if err != nil {
		return nil, err
	}
	return &tokens, nil
}
