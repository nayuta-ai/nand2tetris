package analyzer

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var symbolCharacter string = "(){}.,;+-*/=<>[]&|~"

var commentFlags bool = false

// Token types
const (
	TokenTypeKeyword        = "keyword"
	TokenTypeSymbol         = "symbol"
	TokenTypeIntConstant    = "integerConstant"
	TokenTypeStringConstant = "stringConstant"
	TokenTypeIdentifier     = "identifier"
)

var keywordSlice = []string{
	"class",
	"constructor",
	"function",
	"method",
	"field",
	"static",
	"var",
	"int",
	"char",
	"boolean",
	"void",
	"true",
	"false",
	"null",
	"this",
	"let",
	"do",
	"if",
	"else",
	"while",
	"return",
}

type TokensInterface interface {
	Next() (Token, error)
	Lookup() (Token, error)
}

func (t *Tokens) Next() (Token, error) {
	t.Index++
	if t.Index >= len(t.Token) {
		return t.Token[0], errors.New("no more tokens")
	}
	return t.Token[t.Index], nil
}

func (t *Tokens) Lookup(i int) (Token, error) {
	if t.Index+i >= len(t.Token) {
		return t.Token[0], errors.New("no more tokens")
	}
	return t.Token[t.Index+i], nil
}

type Tokens struct {
	Token []Token
	Index int
}

type Token struct {
	Type  string
	Value string
}

func Tokenizer(scanner *bufio.Scanner) (Tokens, error) {
	var token []Token
	for scanner.Scan() {
		content := scanner.Text()
		content = deleteComments(content)
		var currString string
		var stringConstant string
		for _, char := range content {
			if stringConstant != "" {
				if char == '"' {
					token = append(token, Token{"stringConstant", stringConstant[1:]})
					stringConstant = ""
				} else {
					stringConstant += fmt.Sprintf("%c", char)
				}
			} else if char == ' ' || char == '	' {
				token = appendToken(currString, token)
				currString = ""
			} else if char == '"' {
				stringConstant += fmt.Sprintf("%c", char)
			} else if strings.ContainsRune(symbolCharacter, char) {
				token = appendToken(currString, token)
				currString = ""
				token = append(token, Token{"symbol", normalizexXMLSymbol(char)})
			} else {
				currString += fmt.Sprintf("%c", char)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Tokens{}, err
	}
	return Tokens{token, 0}, nil
}

func appendToken(currString string, token []Token) []Token {
	if currString == "" {
		return token
	}
	if Contains(keywordSlice, currString) {
		token = append(token, Token{"keyword", currString})
	} else if _, err := strconv.Atoi(currString); err == nil {
		token = append(token, Token{"integerConstant", currString})
	} else {
		token = append(token, Token{"identifier", currString})
	}
	return token
}

// https://play.golang.org/p/Qg_uv_inCek
// Contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func deleteComments(line string) string {
	var deletedLine string
	// Delete "//" comments.
	commentEntry := strings.Index(line, "//")
	if commentEntry != -1 {
		deletedLine += line[:commentEntry]
		return deletedLine
	}
	// Delete "/** */" comments.
	commentEntry = strings.Index(line, "/**")
	if commentEntry != -1 {
		commentFlags = true
		deletedLine += line[:commentEntry]
	}
	if commentFlags {
		commentEnd := strings.Index(line, "*/")
		if commentEnd != -1 {
			commentFlags = false
			deletedLine += line[commentEnd+2:]
		}
		return deletedLine
	}
	return line
}

func normalizexXMLSymbol(char rune) string {
	if char == '<' {
		return "&lt;"
	} else if char == '>' {
		return "&gt;"
	} else if char == '&' {
		return "&amp;"
	}
	return fmt.Sprintf("%c", char)
}