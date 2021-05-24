// Benchmark for program #8
//
// Author Viktor Seifert

package main

import "testing"

func BenchmarkEight21(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(21)
	}
}
