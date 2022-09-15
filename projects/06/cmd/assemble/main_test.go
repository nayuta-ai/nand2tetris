package main

import "testing"

func BenchmarkMain(b *testing.B) {
	for i := 0; i < 2; i++ {
		main()
	}
}
