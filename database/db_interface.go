package database

import "time"

type DB interface {
	CreateTables() error

	UploadUsers(users []User) error
	CreateUser(user User) error
	CreateTask(task Task) error
	Login(userID int64, token string) error
	ClickInviteFeferals(userID int64) error
	CompleteTask(userID, taskID int64) error
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

type Task struct {
	ID   int64
	Name string
}
