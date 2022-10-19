package main

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		{
			"test_1",
			"function Class1.set 0",
			"//function Class1.set 0\n(Class1.set)\n",
		},
		{
			"test_2",
			"push argument 0",
			"//push argument 0\n@ARG\nD=M\n@0\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n",
		},
		{
			"test_3",
			"pop static 0",
			"//pop static 0\n@Class1.0\nD=A\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n",
		},
		{
			"test_4",
			"push constant 0",
			"//push constant 0\n@0\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n",
		},
		{
			"test_5",
			"return",
			"//return\n@LCL\nD=M\n@R13\nM=D\n@R13\nD=M\n@5\nD=D-A\nA=D\nD=M\n@R14\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nM=D\n@ARG\nD=M\n@SP\nM=D+1\n@R13\nD=M\n@1\nD=D-A\nA=D\nD=M\n@THAT\nM=D\n@R13\nD=M\n@2\nD=D-A\nA=D\nD=M\n@THIS\nM=D\n@R13\nD=M\n@3\nD=D-A\nA=D\nD=M\n@ARG\nM=D\n@R13\nD=M\n@4\nD=D-A\nA=D\nD=M\n@LCL\nM=D\n@R14\nA=M\n0;JMP\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser(tt.command)
			div_got := strings.Split(got, "\n")
			div_want := strings.Split(tt.want, "\n")
			if len(div_got) != len(div_want) {
				t.Errorf("want %d. but %d", len(div_got), len(div_want))
			}
			for i := 0; i < len(div_got); i++ {
				if div_got[i] != div_want[i] {
					t.Errorf("index %d: want %s. but %s", i, div_want[i], div_got[i])
				}
			}
		})
	}
}
