package main

import (
	"testing"
)

// This function test aBinary function
func TestABinary(t *testing.T) {
	preDict()
	tests := []struct {
		name    string
		command aCommand
		want    string
	}{
		{"test_1", aCommand{"R0"}, "0000000000000000\n"},
		{"test_2", aCommand{"sum"}, "0000000000010000\n"},
		{"test_3", aCommand{"256"}, "0000000100000000\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := aBinary(tt.command); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}

// This function test compToBinary function
func TestCompToBinary(t *testing.T) {
	tests := []struct {
		name    string
		command cCommand
		want    string
	}{
		{"test_1", cCommand{"A", "D", ""}, "0110000"},
		{"test_2", cCommand{"0", "", "JMP"}, "0101010"},
		{"test_3", cCommand{"D+M", "AMD", "JGT"}, "1000010"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compToBinary(tt.command); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}

// This function test destToBinary function
func TestDestToBinary(t *testing.T) {
	tests := []struct {
		name    string
		command cCommand
		want    string
	}{
		{"test_1", cCommand{"A", "D", ""}, "010"},
		{"test_2", cCommand{"0", "", "JMP"}, "000"},
		{"test_3", cCommand{"D+M", "AMD", "JGT"}, "111"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := destToBinary(tt.command); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}

// This function test jumpToBinary function
func TestJumpToBinary(t *testing.T) {
	tests := []struct {
		name    string
		command cCommand
		want    string
	}{
		{"test_1", cCommand{"A", "D", ""}, "000"},
		{"test_2", cCommand{"0", "", "JMP"}, "111"},
		{"test_3", cCommand{"D+M", "AMD", "JGT"}, "001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jumpToBinary(tt.command); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}

// This function test cBinary function
func TestCBinary(t *testing.T) {
	tests := []struct {
		name    string
		command cCommand
		want    string
	}{
		{"test_1", cCommand{"A", "D", ""}, "1110110000010000\n"},
		{"test_2", cCommand{"0", "", "JMP"}, "1110101010000111\n"},
		{"test_3", cCommand{"D", "", "JGT"}, "1110001100000001\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := cBinary(tt.command); got != tt.want {
				t.Errorf("want %s. but %s", tt.want, got)
			}
		})
	}
}
