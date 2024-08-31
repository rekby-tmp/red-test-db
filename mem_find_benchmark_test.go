package main

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func BenchmarkFind(b *testing.B) {
	token := genString(rand.New(rand.NewPCG(333, 524)), 30, 30)
	users := generateUsers(123, 1000000)
	b.ResetTimer()
	var res int
	for range b.N {
		res += slices.IndexFunc(users, func(user User) bool {
			return user.Token == token
		})
	}
	b.Log(res)
}
