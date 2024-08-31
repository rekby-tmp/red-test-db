package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	start := time.Now()
	users := generateUsers(123, 1000000)
	duration := time.Since(start)
	fmt.Println(duration)
	runtime.KeepAlive(users)
}
