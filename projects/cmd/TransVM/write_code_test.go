package main

import (
	"strings"
	"testing"
)

// This function test commandFunction function
func TestCommandFunction(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		{"test_1", "function SimpleFunction.test 2", "//function SimpleFunction.test 2\n(SimpleFunction.test)\nD=0\n@SP\nA=M\nM=D\n@SP\nM=M+1\nD=0\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{"test_2", "function f 3", "//function f 3\n(f)\nD=0\n@SP\nA=M\nM=D\n@SP\nM=M+1\nD=0\n@SP\nA=M\nM=D\n@SP\nM=M+1\nD=0\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			div_command := divideCommand(tt.command)
			got, _ := commandFunction(tt.command, div_command)
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

// This function test commandReturn function
func TestCommandReturn(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		{"test_1", "return", "//return\n@LCL\nD=M\n@R13\nM=D\n@R13\nD=M\n@5\nD=D-A\nA=D\nD=M\n@R14\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nM=D\n@ARG\nD=M\n@SP\nM=D+1\n@R13\nD=M\n@1\nD=D-A\nA=D\nD=M\n@THAT\nM=D\n@R13\nD=M\n@2\nD=D-A\nA=D\nD=M\n@THIS\nM=D\n@R13\nD=M\n@3\nD=D-A\nA=D\nD=M\n@ARG\nM=D\n@R13\nD=M\n@4\nD=D-A\nA=D\nD=M\n@LCL\nM=D\n@R14\nA=M\n0;JMP\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commandReturn(tt.command)
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
