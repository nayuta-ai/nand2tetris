package analyzer

type ParseTree struct {
	Type     string
	Value    string
	Children []*ParseTree
}

func parser(tokens Tokens) (ParseTree, error) {
	return parseClass(tokens)
}

func parseClass(tokens Tokens) (ParseTree, error) {
	tree := newParseTree("class", "")
	// keyword
	tokens, err := tree.addChild(tokens)
	if err != nil {
		return ParseTree{}, err
	}
	// identifier
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return ParseTree{}, err
	}
	// sybmol
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return ParseTree{}, err
	}
	// classVarDec
	for {
		if !contains([]string{"static", "field"}, tokens.Token[tokens.Index].Value) {
			break
		}
		tokens, err = tree.parseClassVarDec(tokens)
		if err != nil {
			return ParseTree{}, err
		}
	}
	// subroutineDec
	for {
		if !contains([]string{"constructor", "function", "method"}, tokens.Token[tokens.Index].Value) {
			break
		}
		tokens, err = tree.parseSubroutineDec(tokens)
		if err != nil {
			return ParseTree{}, err
		}
	}
	// symbol
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return ParseTree{}, err
	}
	return *tree, nil
}

func (pt *ParseTree) parseClassVarDec(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("classVarDec", "")
	for {
		if tokens.Token[tokens.Index].Value == ";" {
			tokens, err = tree.addChild(tokens)
			if err != nil {
				return tokens, err
			}
			break
		}
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseSubroutineDec(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("subroutineDec", "")
	// keyword
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// keyword
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// identifier
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// parameterList
	tokens, err = tree.parseParameterList(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// subroutineBody
	tokens, err = tree.parseSubroutineBody(tokens)
	if err != nil {
		return tokens, err
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseSubroutineBody(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("subroutineBody", "")
	// symbol
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	for {
		if tokens.Token[tokens.Index].Value == "}" {
			// symbol
			tokens, err = tree.addChild(tokens)
			if err != nil {
				return tokens, err
			}
			break
		} else if tokens.Token[tokens.Index].Value == "var" {
			// varDec
			tokens, err = tree.parseVarDec(tokens)
			if err != nil {
				return tokens, nil
			}
		} else {
			// statements
			tokens, err = tree.parseStatements(tokens)
			if err != nil {
				return tokens, nil
			}
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseVarDec(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("varDec", "")
	for {
		if tokens.Token[tokens.Index].Value == ";" {
			tokens, err = tree.addChild(tokens)
			if err != nil {
				return tokens, err
			}
			break
		}
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseStatements(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("statements", "")
	for {
		if tokens.Token[tokens.Index].Value == "do" {
			tokens, err = tree.parseDo(tokens)
			if err != nil {
				return tokens, err
			}
		} else if tokens.Token[tokens.Index].Value == "let" {
			tokens, err = tree.parseLet(tokens)
			if err != nil {
				return tokens, err
			}
		} else if tokens.Token[tokens.Index].Value == "return" {
			tokens, err = tree.parseReturn(tokens)
			if err != nil {
				return tokens, err
			}
		} else if tokens.Token[tokens.Index].Value == "if" {
			tokens, err = tree.parseIf(tokens)
			if err != nil {
				return tokens, err
			}
		} else if tokens.Token[tokens.Index].Value == "while" {
			tokens, err = tree.parseWhile(tokens)
			if err != nil {
				return tokens, err
			}
		} else {
			break
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseDo(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("doStatement", "")
	// do
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// subroutineCall
	tokens, err = tree.parseSubroutineCall(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol(;)
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseSubroutineCall(tokens Tokens) (Tokens, error) {
	var err error
	for {
		if tokens.Token[tokens.Index].Value == "(" {
			tokens, err = pt.addChild(tokens)
			if err != nil {
				return tokens, err
			}
			tokens, err = pt.parseExpressionList(tokens)
			if err != nil {
				return tokens, err
			}
			tokens, err = pt.addChild(tokens)
			if err != nil {
				return tokens, err
			}
			break
		}
		tokens, err = pt.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	}
	return tokens, nil
}

func (pt *ParseTree) parseLet(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("letStatement", "")
	// symbol
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// identifier(varName)
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	if tokens.Token[tokens.Index].Value == "[" {
		// symbol("[")
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
		// expression
		tokens, err = tree.parseExpression(tokens)
		if err != nil {
			return tokens, err
		}
		// symbol("]")
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	}
	// symbol(=)
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// expression
	tokens, err = tree.parseExpression(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol(;)
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseReturn(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("returnStatement", "")
	// return
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	if tokens.Token[tokens.Index].Value != ";" {
		tokens, err = tree.parseExpression(tokens)
		if err != nil {
			return tokens, err
		}
	}
	// symbol (;)
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

var op = []string{"+", "-", "*", "/", "&amp;", "&lt;", "&gt;", "|", "="}

func (pt *ParseTree) parseExpression(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("expression", "")
	for {
		if contains([]string{"]", ")", ";", ","}, tokens.Token[tokens.Index].Value) {
			break
		} else if contains(op, tokens.Token[tokens.Index].Value) && len(tree.Children) != 0 {
			tokens, err = tree.addChild(tokens)
			if err != nil {
				return tokens, err
			}
		} else {
			tokens, err = tree.parseTerm(tokens)
			if err != nil {
				return tokens, err
			}
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

var unaryOp = []string{"-", "~"}

func (pt *ParseTree) parseTerm(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("term", "")
	if tokens.Token[tokens.Index].Value == "(" {
		// symbol ("(")
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
		// expression
		tokens, err = tree.parseExpression(tokens)
		if err != nil {
			return tokens, err
		}
		// symbol (")")
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	} else if contains(unaryOp, tokens.Token[tokens.Index].Value) {
		// symbol (unaryOp)
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
		// term
		tokens, err = tree.parseTerm(tokens)
		if err != nil {
			return tokens, err
		}
	} else {
		// identifier
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	}
	if tokens.Token[tokens.Index].Value == "[" {
		// symbol ("[")
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
		// expression
		tokens, err = tree.parseExpression(tokens)
		if err != nil {
			return tokens, err
		}
		// symbol ("]")
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	} else if tokens.Token[tokens.Index].Value == "(" || tokens.Token[tokens.Index].Value == "." {
		// subroutineCall
		tokens, err = tree.parseSubroutineCall(tokens)
		if err != nil {
			return tokens, err
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseExpressionList(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("expressionList", "")
	for {
		if tokens.Token[tokens.Index].Value == ")" {
			break
		} else if contains([]string{","}, tokens.Token[tokens.Index].Value) {
			tokens, err = tree.addChild(tokens)
			if err != nil {
				return tokens, err
			}
		}
		tokens, err = tree.parseExpression(tokens)
		if err != nil {
			return tokens, err
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseIf(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("ifStatement", "")
	// if
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol ("(")
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// expression
	tokens, err = tree.parseExpression(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol (")")
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol ("{")
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// statements
	tokens, err = tree.parseStatements(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol ("}")
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	if tokens.Token[tokens.Index].Value == "else" {
		// else
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
		// symbol ("{")
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
		// statements
		tokens, err = tree.parseStatements(tokens)
		if err != nil {
			return tokens, err
		}
		// symbol ("}")
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseWhile(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("whileStatement", "")
	// while
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol ("(")
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// expression
	tokens, err = tree.parseExpression(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol (")")
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol ("{")
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	// statements
	tokens, err = tree.parseStatements(tokens)
	if err != nil {
		return tokens, err
	}
	// symbol ("}")
	tokens, err = tree.addChild(tokens)
	if err != nil {
		return tokens, err
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func (pt *ParseTree) parseParameterList(tokens Tokens) (Tokens, error) {
	var err error
	tree := newParseTree("parameterList", "")
	for {
		if tokens.Token[tokens.Index].Value == ")" {
			break
		}
		tokens, err = tree.addChild(tokens)
		if err != nil {
			return tokens, err
		}
	}
	pt.Children = append(pt.Children, tree)
	return tokens, nil
}

func newParseTree(t string, v string) *ParseTree {
	var c []*ParseTree
	return &ParseTree{t, v, c}
}

func (pt *ParseTree) addChild(tokens Tokens) (Tokens, error) {
	tk, err := tokens.Lookup(0)
	if err != nil {
		return tokens, err
	}
	pt.Children = append(pt.Children, newParseTree(tk.Type, tk.Value))
	tokens.Next()
	return tokens, nil
}
