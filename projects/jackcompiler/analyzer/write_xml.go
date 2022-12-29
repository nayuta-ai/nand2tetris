package analyzer

import (
	"fmt"
	"os"
	"strings"
)

func convertToXMLformat(tree ParseTree, file string) error {
	contents := convertParseTreeToString(&tree, 0)
	f, err := os.Create(file + ".xml")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	if err != nil {
		return err
	}
	return nil
}

func convertParseTreeToString(tree *ParseTree, indent int) string {
	tab := strings.Repeat("  ", indent)
	var line string
	if Contains([]string{"keyword", "symbol", "identifier", "integerConstant", "stringConstant"}, tree.Type) {
		line += " " + tree.Value + " "
	} else {
		line += "\n"
		for _, leaf := range tree.Children {
			line += convertParseTreeToString(leaf, indent+1)
		}
		line += tab
	}
	return tab + generateTag(tree.Type, line) + "\n"
}

func generateTag(tag string, line string) string {
	return "<" + tag + ">" + line + "</" + tag + ">"
}

func convertToTXMLformat(token Tokens) error {
	var contents string
	contents += "<tokens>\n"
	for {
		tk, err := token.Lookup(0)
		if err != nil {
			break
		}
		contents += fmt.Sprintf("<%s> %s </%s>\n", tk.Type, tk.Value, tk.Type)
		token.Next()
	}
	contents += "</tokens>"
	f, err := os.Create("sampleT.xml")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	if err != nil {
		return err
	}
	return nil
}

func fetchFileName(path string) string {
	files := strings.Split(path, "/")
	return files[len(files)-1]
}
