package engine

import (
	"os"
	"strconv"
)

const (
	SegmentARG    = "argument"
	SegmentLOCAL  = "local"
	SegmentSTATIC = "static"
	SegmentTHIS   = "this"
	SegmentTHAT   = "that"
	SegmentPOINT  = "pointer"
	SegmentTEMP   = "temp"
	SegmentCONST  = "constant"
)

// arithmetic commands
const (
	ArithmCmdADD = "add"
	ArithmCmdSUB = "sub"
	ArithmCmdNEG = "neg"
	ArithmCmdEQ  = "eq"
	ArithmCmdGT  = "gt"
	ArithmCmdLT  = "lt"
	ArithmCmdAND = "and"
	ArithmCmdOR  = "or"
	ArithmCmdNOT = "not"
)

var ArithmSymbols = map[string]string{
	"+":     ArithmCmdADD,
	"-":     ArithmCmdSUB,
	"*":     "call Math.multiply 2",
	"/":     "call Math.divide 2",
	"&amp;": ArithmCmdAND,
	"|":     ArithmCmdOR,
	"&lt;":  ArithmCmdLT,
	"&gt;":  ArithmCmdGT,
	"=":     ArithmCmdEQ,
}

func writeComment(f *os.File, comment string) error {
	_, err := f.WriteString("\n// " + comment + "\n")
	return err
}

func writeLine(f *os.File, line string) error {
	_, err := f.WriteString(line + "\n")
	return err
}

func writePush(f *os.File, segment string, idx int) error {
	return writeLine(f, "push "+segment+" "+strconv.Itoa(idx))
}

func writePop(f *os.File, segment string, idx int) error {
	return writeLine(f, "pop "+segment+" "+strconv.Itoa(idx))
}

func writeArithmetic(f *os.File, command string) error {
	return writeLine(f, command)
}

func writeLabel(f *os.File, label string) error {
	return writeLine(f, "label "+label)
}

func writeGoto(f *os.File, label string) error {
	return writeLine(f, "goto "+label)
}

func writeIf(f *os.File, label string) error {
	return writeLine(f, "if-goto "+label)
}

func writeCall(f *os.File, name string, nArgs int) error {
	return writeLine(f, "call "+name+" "+strconv.Itoa(nArgs))
}

func writeFunction(f *os.File, name string, nLocals int) error {
	return writeLine(f, "function "+name+" "+strconv.Itoa(nLocals))
}

func writeReturn(f *os.File) error {
	return writeLine(f, "return")
}
