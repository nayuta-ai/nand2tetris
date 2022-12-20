package analyzer

import (
	"bufio"
	"os"
)

func SyntaxAnalyzer(fp *os.File) ([]Token, error) {
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	token, err := tokenizer(scanner)
	if err != nil {
		return nil, err
	}
	return token, nil
}
