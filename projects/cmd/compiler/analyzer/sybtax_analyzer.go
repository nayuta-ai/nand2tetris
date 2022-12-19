package analyzer

import (
	"bufio"
	"fmt"
	"os"
)

func SyntaxAnalyzer(fp *os.File) error {
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
