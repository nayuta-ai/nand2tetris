package engine

// identifier type constants
const (
	IdentifierTypeStatic = "static"
	IdentifierTypeField  = "field"
	IdentifierTypeArg    = "argument"
	IdentifierTypeVar    = "local"
)

type TableStruct struct {
	Name   string
	Type   string
	Kind   string
	Number int
}

type Table map[string]*TableStruct

type SymbolTableInterface interface {
	Define(string, string, string)
	StartSubroutine()
	VarCount(string) int
	KindOf(string) string
	TypeOf(string) string
	IndexOf(string) int
}

type SymbolTable struct {
	SymbolTableInterface
	classTable      Table
	subroutineTable Table
}

func New() *SymbolTable {
	return &SymbolTable{
		classTable:      make(Table),
		subroutineTable: make(Table),
	}
}

func (st *SymbolTable) StartSubroutine() {
	st.subroutineTable = make(Table)
}

func (st *SymbolTable) Define(name, varType, kind string) {
	entry := &TableStruct{name, varType, kind, st.VarCount(kind)}
	if kind == IdentifierTypeVar || kind == IdentifierTypeArg {
		st.subroutineTable[name] = entry
	} else {
		st.classTable[name] = entry
	}
}

func (st *SymbolTable) KindOf(name string) string {
	kind := ""
	if val, ok := st.classTable[name]; ok {
		kind = val.Kind
	}
	if val, ok := st.subroutineTable[name]; ok {
		kind = val.Kind
	}
	if kind == IdentifierTypeField {
		kind = "this"
	}
	return kind
}

func (st *SymbolTable) TypeOf(name string) string {
	if val, ok := st.classTable[name]; ok {
		return val.Type
	}
	if val, ok := st.subroutineTable[name]; ok {
		return val.Type
	}
	return ""
}

func (st *SymbolTable) IndexOf(name string) int {
	if val, ok := st.classTable[name]; ok {
		return val.Number
	}
	if val, ok := st.subroutineTable[name]; ok {
		return val.Number
	}
	return -1
}

func (st *SymbolTable) VarCount(kind string) int {
	if kind == IdentifierTypeVar || kind == IdentifierTypeArg {
		return varCount(st.subroutineTable, kind)
	}
	return varCount(st.classTable, kind)
}

func varCount(t Table, kind string) int {
	var number int
	for _, val := range t {
		if val.Kind == kind {
			number++
		}
	}
	return number
}
