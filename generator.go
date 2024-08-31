package main

import (
	"math/rand/v2"
	"strings"
	"time"
)

func generateUsers(seed uint64, count int) []User {
	users := make([]User, count)
	usersID := generateIDs(seed, count)
	for i, id := range usersID {
		users[i] = generateUser(id)
	}
	return users
}

func generateIDs(seed uint64, count int) []int64 {
	rnd := rand.New(rand.NewPCG(seed, seed))
	res := make([]int64, count)
	for i := range res {
		res[i] = rnd.Int64()
	}
	return res
}

func generateUser(id int64) (res User) {
	pcg := rand.NewPCG(uint64(id), uint64(id))
	rnd := rand.New(pcg)

	res.ID = id
	res.Token = genString(rnd, 30, 30)
	res.Referal = rnd.Int64N(1000)
	res.Rk = rnd.Int64N(1000)
	res.Avatar = genString(rnd, 20, 30)
	res.FirstLogin = time.Date(2000+rnd.IntN(23), time.Month(rnd.Int64N(12)), rnd.IntN(20), rnd.IntN(12), rnd.IntN(60), rnd.IntN(60), 0, time.UTC)
	res.LastLogin = time.Date(2000+rnd.IntN(23), time.Month(rnd.Int64N(12)), rnd.IntN(20), rnd.IntN(12), rnd.IntN(60), rnd.IntN(60), 0, time.UTC)
	res.LastLeave = time.Date(2000+rnd.IntN(23), time.Month(rnd.Int64N(12)), rnd.IntN(20), rnd.IntN(12), rnd.IntN(60), rnd.IntN(60), 0, time.UTC)
	res.InvitedReferals = rnd.Int64N(100)
	res.RaffleRules = rnd.Int64N(100)
	res.InviteCopy = rnd.Int64N(1000)
	return res
}

func generateTasks(seed uint64, count int) []Task {
	tasks := make([]Task, count)
	ids := generateIDs(seed, count)
	for i, id := range ids {
		tasks[i] = generateTask(id)
	}
	return tasks
}

func generateTask(id int64) (res Task) {
	pcg := rand.NewPCG(uint64(id), uint64(id))
	rnd := rand.New(pcg)

	res.ID = id
	res.Name = genString(rnd, 10, 20)
	return res
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFJHIJKLMNOPQRSTUVWXYZ0123456789"

func genString(rnd *rand.Rand, minLen, maxLen int) string {
	lenDiff := maxLen - minLen
	dstLen := rnd.IntN(lenDiff+1) + minLen
	sb := strings.Builder{}
	sb.Grow(dstLen)
	for range dstLen {
		sb.WriteByte(chars[rnd.IntN(len(chars))])
	}
	return sb.String()
}
