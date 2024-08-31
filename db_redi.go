package main

import (
	"errors"
	"github.com/redi-db/redi.db.go"
)

type RediDB struct {
	redi     redidb.DB
	database string
}

const (
	rediUsersCollection    = "users"
	rediTasksCollection    = "tasks"
	rediUserTaskCollection = "user_tasks"
)

func NewRediDB(host string, port int, login, password, database string) (*RediDB, error) {
	db := &RediDB{
		redi: redidb.DB{
			Login:    login,
			Password: password,
			Ip:       host,
			Port:     port,
		},
		database: database,
	}
	collection := db.redi.Database(database)
	collection = collection.Collection(rediTasksCollection)
	if _, err := collection.SearchOne(redidb.Filter{
		"id": -1,
	}); !errors.Is(err, redidb.NOT_FOUND) {
		return nil, err
	}
	return db, nil
}

func (r *RediDB) UploadUsers(users []User) error {
	collection := r.redi.Database(r.database)
	collection = collection.Collection(rediUsersCollection)

	data := make([]redidb.CreateData, len(users))
	for i := range users {
		data[i] = r.createUserData(&users[i])
	}
	_, err := collection.Create(data...)
	return err
}

func (r *RediDB) CreateUser(user User) error {
	collection := r.redi.Database(r.database)
	collection = collection.Collection(rediUsersCollection)

	_, err := collection.Create(r.createUserData(&user))
	return err
}

func (r *RediDB) createUserData(user *User) redidb.CreateData {
	return redidb.CreateData{
		"id":              user.ID,
		"token":           user.Token,
		"referal":         user.Referal,
		"rk":              user.Rk,
		"avatar":          user.Avatar,
		"first_login":     user.FirstLogin,
		"last_login":      user.LastLogin,
		"last_leave":      user.LastLeave,
		"invite_referals": user.InvitedReferals,
		"raffle_rules":    user.RaffleRules,
		"invite_copy":     user.InviteCopy,
	}
}

func (r *RediDB) CreateTask(task Task) error {
	collection := r.redi.Database(r.database)
	collection = collection.Collection(rediTasksCollection)

	createData := redidb.CreateData{
		"id":   task.ID,
		"name": task.Name,
	}
	_, err := collection.Create(createData)
	return err
}

func (r *RediDB) Login(userID int64, token string) error {
	collection := r.redi.Database(r.database)
	collection = collection.Collection(rediUsersCollection)

	_, err := collection.SearchOne(redidb.Filter{
		"id":    userID,
		"token": token,
	})

	if err != nil {
		return err
	}

	// TODO: update last login
	return nil
}

func (r *RediDB) ClickInviteFeferals(userID int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RediDB) CompleteTask(userID, taskID int64) error {
	collection := r.redi.Database(r.database)
	collection = collection.Collection(rediUserTaskCollection)

	_, err := collection.SearchOrCreate(redidb.Filter{
		"user_id": userID,
		"task_id": taskID,
	}, redidb.CreateData{
		"user_id": userID,
		"task_id": taskID,
	})
	return err
}
