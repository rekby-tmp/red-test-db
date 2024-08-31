package main

import "fmt"

func main() {
	db := must(NewRediDB("localhost", 5001, "root", "root", "anton"))
	user := generateUser(123)
	must0(db.CreateUser(user))
	fmt.Println(db)
}

func must0(err error) {
	if err != nil {
		panic(err)
	}
}

func must[R any](res R, err error) R {
	must0(err)
	return res
}
