package main

import (
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkFind(b *testing.B) {
	token := genString(rand.New(rand.NewPCG(333, 524)), 30, 30)
	users := generateUsers(123, 1000000)
	b.ResetTimer()
	var res atomic.Int64
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for range b.N {
				res.Add(int64(slices.IndexFunc(users, func(user User) bool {
					return user.Token == token
				})))
			}
		}()
	}
	wg.Wait()
	b.Log(res.Load())
}
