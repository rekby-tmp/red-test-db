package main

import "time"

type DB interface {
	CreateUsers(users []User) error
	Login(userID int64)
	ClickInviteFeferals(userID int64) error
}

type User struct {
	ID      int64
	Token   string
	Referal int64
	Rk      int64
	Avatar  string

	FirstLogin      time.Time
	LastLogin       time.Time
	LastLeave       time.Time
	InvitedReferals int64
	RaffleRules     int64
	InviteCopy      int64
}
