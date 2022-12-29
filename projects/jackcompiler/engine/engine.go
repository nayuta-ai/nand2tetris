package engine

import (
	"fmt"
	"os"
	"strconv"

	"jackcompiler/analyzer"
	"jackcompiler/validation"
)

type CompilerEngine struct {
	tokens         *analyzer.Tokens
	fp             *os.File
	st             *SymbolTable
	className      string
	subroutineKind string
	subroutineName string
	labelCounter   int
}

func Compiler(path string) error {
	newFile := path[:len(path)-5] + ".vm"
	f, err := os.Create(newFile)
	if err != nil {
		return err
	}
	defer f.Close()
	token, err := analyzer.GetTokens(path)
	if err != nil {
		return err
	}
	table := New()
	ce := &CompilerEngine{token, f, table, "", "", "", 0}
	err = ce.compileClass()
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileClass() error {
	_, err := ce.tokens.Next() // class
	if err != nil {
		return err
	}
	token, err := ce.tokens.Lookup(0)
	if err != nil {
		return err
	}
	ce.className = token.Value // Save class name
	_, err = ce.tokens.Next()  // class name
	if err != nil {
		return err
	}
	err = writeComment(ce.fp, "Class "+ce.className) // write comment
	if err != nil {
		return err
	}
	_, err = ce.tokens.Next() // {
	if err != nil {
		return err
	}
	err = ce.compileClassVarDec()
	if err != nil {
		return err
	}
	err = ce.compileSubroutineDec()
	if err != nil {
		return err
	}
	err = ce.FinishCompile(validation.Match("}")) // }
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileClassVarDec() error {
	// early stopping
	if !analyzer.Contains([]string{"static", "field"}, ce.getCurrentTokenValue()) {
		return nil
	}
	kind, err := ce.NextValidatedValue(validation.OneOf(IdentifierTypeField, IdentifierTypeStatic))
	if err != nil {
		return err
	}
	varType := ce.getCurrentTokenValue()
	_, err = ce.tokens.Next()
	if err != nil {
		return err
	}
	for {
		name, err := ce.NextValidatedType(analyzer.TokenTypeIdentifier)
		if err != nil {
			return err
		}
		ce.st.Define(name, varType, kind)
		if ce.getCurrentTokenValue() != "," {
			break
		}
		_, err = ce.NextValidatedValue(validation.Match(",")) // ','
		if err != nil {
			return err
		}
	}
	_, err = ce.NextValidatedValue(validation.Match(";")) // ';'
	if err != nil {
		return err
	}
	err = ce.compileClassVarDec()
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileSubroutineDec() error {
	// early stopping
	if !analyzer.Contains([]string{"method", "function", "constructor"}, ce.getCurrentTokenValue()) {
		return nil
	}
	ce.st.StartSubroutine()
	token, err := ce.NextValidatedValue(validation.OneOf("method", "function", "constructor"))
	if err != nil {
		return err
	}
	ce.subroutineKind = token
	// (void | type)
	_, err = ce.tokens.Next()
	if err != nil {
		return err
	}
	token = ce.getCurrentTokenValue()
	ce.subroutineName = token // Save subroutine name
	_, err = ce.tokens.Next() // subroutine name
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("("))
	if err != nil {
		return err
	}
	err = ce.compileParameterList()
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match(")"))
	if err != nil {
		return err
	}
	err = ce.compileSubroutineBody()
	if err != nil {
		return err
	}
	if ce.getCurrentTokenValue() != "}" {
		err = ce.compileSubroutineDec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ce *CompilerEngine) compileParameterList() error {
	if ce.subroutineKind == "method" {
		ce.st.Define("this", ce.className, IdentifierTypeArg)
	}
	if ce.getCurrentTokenValue() == ")" {
		return nil
	}
	for {
		varType := ce.getCurrentTokenValue()
		_, err := ce.tokens.Next()
		if err != nil {
			return err
		}
		name, err := ce.NextValidatedType("identifier")
		if err != nil {
			return err
		}
		ce.st.Define(name, varType, IdentifierTypeArg)
		if ce.getCurrentTokenValue() == "," {
			_, err = ce.NextValidatedValue(validation.Match(","))
			if err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (ce *CompilerEngine) compileSubroutineBody() error {
	_, err := ce.NextValidatedValue(validation.Match("{"))
	if err != nil {
		return err
	}
	for {
		if ce.getCurrentTokenValue() != "var" {
			break
		}
		err = ce.compileVarDec()
		if err != nil {
			return err
		}
	}
	err = writeComment(ce.fp, "Subroutine "+ce.subroutineKind+" "+ce.subroutineName)
	if err != nil {
		return err
	}
	err = writeFunction(ce.fp, ce.className+"."+ce.subroutineName, ce.st.VarCount("local"))
	if err != nil {
		return err
	}
	if ce.subroutineKind == "constructor" {
		writePush(ce.fp, SegmentCONST, ce.st.VarCount("field"))
		writeCall(ce.fp, "Memory.alloc", 1)
		writePop(ce.fp, SegmentPOINT, 0)
	}
	if ce.subroutineKind == "method" {
		writePush(ce.fp, SegmentARG, 0)
		writePop(ce.fp, SegmentPOINT, 0)
	}
	err = ce.compileStatements()
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("}"))
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileVarDec() error {
	// 'var' type varName (',' varName)* ';'
	_, err := ce.NextValidatedValue(validation.Match("var"))
	if err != nil {
		return err
	}
	varType := ce.getCurrentTokenValue()
	// type
	_, err = ce.tokens.Next()
	if err != nil {
		return err
	}
	for {
		// varName
		varName, err := ce.NextValidatedType(analyzer.TokenTypeIdentifier)
		if err != nil {
			return err
		}
		ce.st.Define(varName, varType, IdentifierTypeVar)
		if ce.getCurrentTokenValue() != "," {
			break
		}
		_, err = ce.NextValidatedValue(validation.Match(","))
		if err != nil {
			return err
		}
	}
	_, err = ce.NextValidatedValue(validation.Match(";"))
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileStatements() error {
	var err error
	for {
		switch ce.getCurrentTokenValue() {
		case "do":
			err = ce.compileDo()
			if err != nil {
				return err
			}
		case "let":
			err = ce.compileLet()
			if err != nil {
				return err
			}
		case "while":
			err = ce.compileWhile()
			if err != nil {
				return err
			}
		case "if":
			err = ce.compileIf()
			if err != nil {
				return err
			}
		case "return":
			err = ce.compileReturn()
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (ce *CompilerEngine) compileDo() error {
	err := writeComment(ce.fp, "Do Statement") // write comment
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("do"))
	if err != nil {
		return err
	}
	err = ce.compileSubroutine()
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match(";"))
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileLet() error {
	var err error
	var isArray bool
	err = writeComment(ce.fp, "Let Statement")
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("let"))
	if err != nil {
		return err
	}
	// varName
	identifier, err := ce.NextValidatedType(analyzer.TokenTypeIdentifier)
	if err != nil {
		return err
	}
	segment := ce.st.KindOf(identifier)
	index := ce.st.IndexOf(identifier)
	// ('[' expression ']')?
	if ce.getCurrentTokenValue() == "[" {
		isArray = true
		_, err = ce.NextValidatedValue(validation.Match("["))
		if err != nil {
			return err
		}
		err = ce.compileExpression()
		if err != nil {
			return err
		}
		_, err = ce.NextValidatedValue(validation.Match("]"))
		if err != nil {
			return err
		}
		writePush(ce.fp, segment, index)
		writeArithmetic(ce.fp, ArithmCmdADD)
	}
	_, err = ce.NextValidatedValue(validation.Match("="))
	if err != nil {
		return err
	}
	err = ce.compileExpression()
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match(";"))
	if err != nil {
		return err
	}
	if isArray {
		writePop(ce.fp, SegmentTEMP, 0)
		writePop(ce.fp, SegmentPOINT, 1)
		writePush(ce.fp, SegmentTEMP, 0)
		writePop(ce.fp, SegmentTHAT, 0)
	} else {
		err = writePop(ce.fp, segment, index)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ce *CompilerEngine) compileWhile() error {
	// 'while' '(' expression ')' '{' statements '}'
	err := writeComment(ce.fp, "While Statement")
	if err != nil {
		return err
	}
	ce.labelCounter++
	label := ce.className + ".while." + strconv.Itoa(ce.labelCounter) + ".L1"
	err = writeLabel(ce.fp, label)
	if err != nil {
		return err
	}
	// while
	_, err = ce.NextValidatedValue(validation.Match("while"))
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("("))
	if err != nil {
		return err
	}
	err = ce.compileExpression()
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match(")"))
	if err != nil {
		return err
	}
	// not
	err = writeArithmetic(ce.fp, ArithmCmdNOT)
	if err != nil {
		return err
	}
	gotoLabel := ce.className + ".while." + strconv.Itoa(ce.labelCounter) + ".L2"
	err = writeIf(ce.fp, gotoLabel)
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("{"))
	if err != nil {
		return err
	}
	err = ce.compileStatements()
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("}"))
	if err != nil {
		return err
	}
	// End of while statement
	err = writeComment(ce.fp, "End of While Statement")
	if err != nil {
		return err
	}
	err = writeGoto(ce.fp, label)
	if err != nil {
		return err
	}
	err = writeLabel(ce.fp, gotoLabel)
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileIf() error {
	// 'if' '(' expression ')' '{' statements '}'
	// ('else' '{' statements '}')?
	writeComment(ce.fp, "If Statement")
	ce.labelCounter++

	_, err := ce.NextValidatedValue(validation.Match("if"))
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("("))
	if err != nil {
		return err
	}
	err = ce.compileExpression()
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match(")"))
	if err != nil {
		return err
	}
	writeArithmetic(ce.fp, ArithmCmdNOT)
	label := ce.className + ".if." + strconv.Itoa(ce.labelCounter) + ".L1"
	writeIf(ce.fp, label)

	_, err = ce.NextValidatedValue(validation.Match("{"))
	if err != nil {
		return err
	}
	gotoLabel := ce.className + ".if." + strconv.Itoa(ce.labelCounter) + ".L2"
	err = ce.compileStatements()
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("}"))
	if err != nil {
		return err
	}
	writeGoto(ce.fp, gotoLabel)
	writeLabel(ce.fp, label)

	if ce.getCurrentTokenValue() == "else" {
		writeComment(ce.fp, "Else Statement")
		_, err := ce.NextValidatedValue(validation.Match("else"))
		if err != nil {
			return err
		}
		_, err = ce.NextValidatedValue(validation.Match("{"))
		if err != nil {
			return err
		}
		err = ce.compileStatements()
		if err != nil {
			return err
		}
		_, err = ce.NextValidatedValue(validation.Match("}"))
		if err != nil {
			return err
		}
	}
	writeLabel(ce.fp, gotoLabel)
	return nil
}

func (ce *CompilerEngine) compileReturn() error {
	// Grammar: 'return' expression? ';'
	err := writeComment(ce.fp, "Return Statement")
	if err != nil {
		return err
	}
	_, err = ce.NextValidatedValue(validation.Match("return"))
	if err != nil {
		return err
	}
	if ce.getCurrentTokenValue() != ";" {
		err = ce.compileExpression()
		if err != nil {
			return err
		}
	} else {
		err = writePush(ce.fp, SegmentCONST, 0)
		if err != nil {
			return err
		}
	}
	_, err = ce.NextValidatedValue(validation.Match(";"))
	if err != nil {
		return err
	}
	err = writeReturn(ce.fp)
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileSubroutine() error {
	// Grammar: subroutineName '(' expressionList ')' |
	// (className|varName) '.' subroutineName '(' expressionList ')'
	var funcName string
	var identifierType string
	var nArgs int
	identifier, err := ce.NextValidatedType("identifier")
	if err != nil {
		return err
	}
	className := identifier
	if ce.getCurrentTokenValue() == "." {
		_, err = ce.NextValidatedValue(validation.Match("."))
		if err != nil {
			return err
		}
		funcName, err = ce.NextValidatedType("identifier")
		if err != nil {
			return err
		}
		identifierType = ce.st.TypeOf(identifier)
	} else {
		className = ce.className
		funcName = identifier
		nArgs++
		writePush(ce.fp, SegmentPOINT, 0)
	}
	if identifierType != "" {
		segment := ce.st.KindOf(identifier)
		index := ce.st.IndexOf(identifier)
		err = writePush(ce.fp, segment, index)
		if err != nil {
			return err
		}
		nArgs++
		className = identifierType
	}
	_, err = ce.NextValidatedValue(validation.Match("("))
	if err != nil {
		return err
	}
	nExp, err := ce.compileExpressionList()
	if err != nil {
		return err
	}
	nArgs += nExp
	_, err = ce.NextValidatedValue(validation.Match(")"))
	if err != nil {
		return err
	}
	if funcName == "" {
		err = writeCall(ce.fp, className, nArgs)
	} else {
		err = writeCall(ce.fp, className+"."+funcName, nArgs)
	}
	if err != nil {
		return err
	}
	err = writePop(ce.fp, SegmentTEMP, 0)
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilerEngine) compileExpressionList() (int, error) {
	var nExp int
	var err error
	for {
		// stopping condition
		if ce.getCurrentTokenValue() == ")" {
			return nExp, nil
		} else if ce.getCurrentTokenValue() == "," {
			_, err = ce.NextValidatedValue(validation.Match(","))
			if err != nil {
				return nExp, err
			}
		}
		nExp++
		err = ce.compileExpression()
		if err != nil {
			return nExp, err
		}
	}
}

func (ce *CompilerEngine) compileExpression() error {
	ops := []string{"+", "-", "*", "/", "&amp;", "&lt;", "&gt;", "|", "="}
	err := ce.compileTerm()
	if err != nil {
		return err
	}
	commands := make([]string, 0)
	for {
		if !validation.OneOf(ops...)(ce.getCurrentTokenValue()) {
			break
		}
		cmd, ok := ArithmSymbols[ce.getCurrentTokenValue()]
		if !ok {
			return fmt.Errorf("value error: invalid arthmSymbols")
		}
		commands = append(commands, cmd)
		_, err = ce.NextValidatedValue(validation.OneOf(ops...))
		if err != nil {
			return err
		}
		err = ce.compileTerm()
		if err != nil {
			return err
		}
	}
	for _, c := range commands {
		err = writeArithmetic(ce.fp, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ce *CompilerEngine) compileTerm() error {
	token, err := ce.tokens.Lookup(0)
	if err != nil {
		return err
	}
	nextToken, err := ce.tokens.Lookup(1)
	if err != nil {
		return err
	}
	if token.Type == analyzer.TokenTypeIntConstant {
		// integerConstant
		num, err := strconv.Atoi(token.Value)
		if err != nil {
			return err
		}
		err = writePush(ce.fp, SegmentCONST, num)
		if err != nil {
			return err
		}
		_, err = ce.NextValidatedType(analyzer.TokenTypeIntConstant)
		if err != nil {
			return err
		}
	} else if token.Type == analyzer.TokenTypeStringConstant {
		// stringConstant
		err = writePush(ce.fp, SegmentCONST, len(token.Value))
		if err != nil {
			return err
		}
		err = writeCall(ce.fp, "String.new", 1)
		if err != nil {
			return err
		}
		for _, r := range token.Value {
			err = writePush(ce.fp, SegmentCONST, int(r))
			if err != nil {
				return err
			}
			err = writeCall(ce.fp, "String.appendChar", 2)
			if err != nil {
				return err
			}
		}
		_, err = ce.NextValidatedType(analyzer.TokenTypeStringConstant)
		if err != nil {
			return err
		}
	} else if token.Type == analyzer.TokenTypeKeyword && validation.Match("true")(token.Value) {
		// keywordConstant true
		writePush(ce.fp, SegmentCONST, 1)
		writeArithmetic(ce.fp, ArithmCmdNEG)
		_, err = ce.NextValidatedValue(validation.Match("true"))
		if err != nil {
			return err
		}
	} else if token.Type == analyzer.TokenTypeKeyword && validation.OneOf("null", "false")(token.Value) {
		// keywordConstant false | null
		writePush(ce.fp, SegmentCONST, 0)
		_, err = ce.NextValidatedValue(validation.OneOf("false", "null"))
		if err != nil {
			return err
		}

	} else if token.Type == analyzer.TokenTypeKeyword && validation.Match("this")(token.Value) {
		// keywordConstant this
		writePush(ce.fp, SegmentPOINT, 0)
		_, err = ce.NextValidatedValue(validation.Match("this"))
		if err != nil {
			return err
		}

	} else if token.Value == "(" {
		// '(' expression ')'
		_, err = ce.NextValidatedValue(validation.Match("("))
		if err != nil {
			return err
		}
		err = ce.compileExpression()
		if err != nil {
			return err
		}
		_, err = ce.NextValidatedValue(validation.Match(")"))
		if err != nil {
			return err
		}
	} else if token.Type == analyzer.TokenTypeSymbol && token.Value == "-" {
		// unaryOp term -
		_, err = ce.NextValidatedValue(validation.Match("-"))
		if err != nil {
			return err
		}
		err = ce.compileTerm()
		if err != nil {
			return err
		}
		writeArithmetic(ce.fp, ArithmCmdNEG)
	} else if token.Type == analyzer.TokenTypeSymbol && token.Value == "~" {
		// unaryOp term ~
		_, err = ce.NextValidatedValue(validation.Match("~"))
		if err != nil {
			return err
		}
		err = ce.compileTerm()
		if err != nil {
			return err
		}
		writeArithmetic(ce.fp, ArithmCmdNOT)
	} else if nextToken.Type == analyzer.TokenTypeSymbol && nextToken.Value == "." {
		// subroutineCall
		identifier, err := ce.NextValidatedType(analyzer.TokenTypeIdentifier)
		if err != nil {
			return err
		}
		_, err = ce.NextValidatedValue(validation.Match(".")) // '.'
		if err != nil {
			return err
		}
		funcName, err := ce.NextValidatedType(analyzer.TokenTypeIdentifier)
		if err != nil {
			return err
		}
		nArgs := 0
		className := identifier
		identifierType := ce.st.TypeOf(identifier)
		if identifierType != "" {
			segment := ce.st.KindOf(identifier)
			index := ce.st.IndexOf(identifier)
			writePush(ce.fp, segment, index)
			nArgs++
			className = identifierType
		}
		_, err = ce.NextValidatedValue(validation.Match("(")) // '('
		if err != nil {
			return err
		}
		nExp, err := ce.compileExpressionList() // expressionList
		if err != nil {
			return err
		}
		nArgs += nExp
		_, err = ce.NextValidatedValue(validation.Match(")")) // ')'
		if err != nil {
			return err
		}
		writeCall(ce.fp, className+"."+funcName, nArgs)
	} else {
		// varName
		varName, err := ce.NextValidatedType(analyzer.TokenTypeIdentifier)
		if err != nil {
			return err
		}
		kind := ce.st.KindOf(varName)
		index := ce.st.IndexOf(varName)
		err = writePush(ce.fp, kind, index)
		if err != nil {
			return err
		}

		// '[' expression ']'
		if ce.getCurrentTokenValue() == "[" {
			_, err = ce.NextValidatedValue(validation.Match("["))
			if err != nil {
				return err
			}
			err = ce.compileExpression()
			if err != nil {
				return err
			}
			_, err = ce.NextValidatedValue(validation.Match("]"))
			if err != nil {
				return err
			}
			writeArithmetic(ce.fp, ArithmCmdADD)
			writePop(ce.fp, SegmentPOINT, 1)
			writePush(ce.fp, SegmentTHAT, 0)
		}
	}
	return nil
}

func (ce *CompilerEngine) getCurrentTokenValue() string {
	token, _ := ce.tokens.Lookup(0)
	return token.Value
}

func (ce *CompilerEngine) getCurrentTokenType() string {
	token, _ := ce.tokens.Lookup(0)
	return token.Type
}

func (ce *CompilerEngine) getNextTokenValue() string {
	token, _ := ce.tokens.Lookup(1)
	return token.Value
}

func (ce *CompilerEngine) NextValidatedValue(rules ...validation.Rule) (string, error) {
	token, err := ce.tokens.Lookup(0)
	if err != nil {
		return "", fmt.Errorf("value error: %+v", err)
	}
	for _, r := range rules {
		if r(token.Value) {
			_, err = ce.tokens.Next()
			if err != nil {
				return "", fmt.Errorf("value error: %+v", err)
			}
			return token.Value, nil
		}
	}
	return "", fmt.Errorf("validation error: %s", token.Value)
}

func (ce *CompilerEngine) NextValidatedType(tType string) (string, error) {
	token, err := ce.tokens.Lookup(0)
	if err != nil {
		return "", fmt.Errorf("value error: %+v", err)
	}
	if token.Type == tType {
		_, err = ce.tokens.Next()
		if err != nil {
			return "", fmt.Errorf("value error: %+v", err)
		}
		return token.Value, nil
	}
	return "", fmt.Errorf("validation error %s: expected %s, but %s", token.Value, tType, token.Type)
}

func (ce *CompilerEngine) FinishCompile(rules ...validation.Rule) error {
	token, err := ce.tokens.Lookup(0)
	if err != nil {
		return fmt.Errorf("value error: %+v", err)
	}
	for _, r := range rules {
		if r(token.Value) {
			return nil
		}
	}
	return fmt.Errorf("validation error: %s", token.Value)
}
