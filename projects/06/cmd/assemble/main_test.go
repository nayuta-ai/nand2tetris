package main

import (
	"os"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	for i := 0; i < 2; i++ {
		filepath := "test/add/Add"
		fp, _ := os.Open(filepath + ".asm")
		defer fp.Close()
		byteInstraction, _ := scanline(fp)
		_ = createHack(byteInstraction, filepath)
	}
}
