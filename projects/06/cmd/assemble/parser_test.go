package main

import (
	"testing"
)

func TestRemoveComment(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{"test_1", "// create", ""},
		{"test_2", "@13", "@13"},
		{"test_3", "D=M // D = second number", "D=M"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeComment(tt.line); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}

func TestLParser(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{"test_1", "@13", ""},
		{"test_2", "(LOOP)", "LOOP"},
		{"test_3", "D=M", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lParse(tt.line); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}

func TestAParse(t *testing.T) {
	tests := []struct {
		name string
		line string
		want aCommand
	}{
		{"test_1", "@13", aCommand{"13"}},
		{"test_2", "@R0", aCommand{"R0"}},
		{"test_3", "@LOOP", aCommand{"LOOP"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aParse(tt.line); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}

func TestCParse(t *testing.T) {
	tests := []struct {
		name string
		line string
		want cCommand
	}{
		{"test_1", "D=A", cCommand{"A", "D", ""}},
		{"test_2", "0;JMP", cCommand{"0", "", "JMP"}},
		{"test_3", "D;JGT", cCommand{"D", "", "JGT"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cParse(tt.line); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}

func TestCommandType(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{"test_1", "D=A", "1110110000010000\n"},
		{"test_2", "@256", "0000000100000000\n"},
		{"test_3", "(LOOP)", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := commandType(tt.line); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}
