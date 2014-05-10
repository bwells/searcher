package main

import "testing"

func BenchmarkRead(b *testing.B) {

	b.StopTimer()
	t := buildTrieFromDB()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		multiMatch(t, "annette")
	}
}
