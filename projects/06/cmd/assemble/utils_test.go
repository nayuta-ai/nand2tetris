package main

import (
	"bufio"
	"os"
	"testing"
)

// This function test argParser function
func TestArgParser(t *testing.T) {
	t.Run("Test ArgParser", func(t *testing.T) {
		if got, _ := parseArgs("--filename", "add/Add"); got != "test/add/Add" {
			t.Errorf("want %s. but %s", "test/add/Add", got)
		}
	})
}

// This function test scanline function
func TestScanline(t *testing.T) {
	fp, _ := os.Open("test/add/Add.asm")
	want := make([]string, 0)
	want = append(want, "0000000000000010\n")
	want = append(want, "1110110000010000\n")
	want = append(want, "0000000000000011\n")
	want = append(want, "1110000010010000\n")
	want = append(want, "0000000000000000\n")
	want = append(want, "1110001100001000\n")
	defer fp.Close()
	t.Run("Test scanline", func(t *testing.T) {
		got, _ := scanline(fp)
		for i, num := range got {
			if num != want[i] {
				t.Errorf("want %s. but %s", want[i], num)
			}
		}
	})
}

// This function test createHack function
func TestCreateHack(t *testing.T) {
	want := make([]string, 0)
	want = append(want, "0000000000000010\n")
	want = append(want, "1110110000010000\n")
	want = append(want, "0000000000000011\n")
	want = append(want, "1110000010010000\n")
	want = append(want, "0000000000000000\n")
	want = append(want, "1110001100001000\n")
	_ = createHack(want, "test/add/Add")
	fp, _ := os.Open("test/add/Add.hack")
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	t.Run("Test scanline", func(t *testing.T) {
		count := 0
		for scanner.Scan() {
			line := scanner.Text()
			if line+"\n" != want[count] {
				t.Errorf("want %s. but %s", want[count], line)
			}
			count += 1
		}
	})
}
