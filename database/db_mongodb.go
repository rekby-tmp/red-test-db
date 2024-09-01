package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"red-db-test/model"
	"strconv"
	"time"
)

type MongoDB struct {
	db              *mongo.Database
	usersCollection *mongo.Collection
	taskCollection  *mongo.Collection
}

const (
	mongoUsersCollection  = "users"
	mongoTasksCollectoins = "tasks"
)

func NewMongo(endpoint, authSource, login, password, database string) (*MongoDB, error) {
	opt := options.Client()
	if login != "" {
		cred := options.Credential{
			AuthSource: authSource,
			Username:   login,
			Password:   password,
		}
		opt.SetAuth(cred)
	}
	opt.ApplyURI(endpoint)
	client, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	db := client.Database(database)
	return &MongoDB{
		db:              db,
		usersCollection: db.Collection(mongoUsersCollection),
		taskCollection:  db.Collection(mongoTasksCollectoins),
	}, nil
}

func (m *MongoDB) CreateTables() error {
	if err := m.db.Collection(mongoUsersCollection).Drop(context.Background()); err != nil {
		return err
	}
	if err := m.db.Collection(mongoTasksCollectoins).Drop(context.Background()); err != nil {
		return err
	}
	m.db.Collection(mongoUsersCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"completedtasks", 1}},
		Options: nil,
	})
	return nil
}

func (m *MongoDB) UploadUsers(users []model.User) error {
	mongoUsers := make([]any, len(users))
	for i := range users {
		mongoUsers[i] = mongoUserFromUser(&users[i])
	}
	_, err := m.usersCollection.InsertMany(context.Background(), mongoUsers)
	return err
}

func (m *MongoDB) CreateUser(user model.User) error {
	_, err := m.usersCollection.InsertOne(context.Background(), mongoUserFromUser(&user))
	return err
}

func (m *MongoDB) CreateTask(task model.Task) error {
	_, err := m.taskCollection.InsertOne(context.Background(), mongoTaskFromTask(task))
	return err
}

func (m *MongoDB) Login(userID int64, token string) error {
	res, err := m.usersCollection.UpdateOne(
		context.Background(),
		bson.D{
			{"_id", strconv.FormatInt(userID, 10)},
			{"token", token},
		},
		bson.D{{"$set", bson.D{{"lastlogin", time.Now().UnixMicro()}}}},
	)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("login failed - unknow id or token")
	}
	return nil
}

func (m *MongoDB) ClickInviteFeferals(userID int64) error {
	update := bson.D{{"$inc", bson.D{{"invitedreferals", 1}}}}
	res, err := m.usersCollection.UpdateByID(
		context.Background(),
		strconv.FormatInt(userID, 10),
		update,
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("id not found")
	}
	return nil
}

func (m *MongoDB) CompleteTask(userID, taskID int64) error {
	update := bson.D{{"$addToSet", bson.D{{"completedtasks", strconv.FormatInt(taskID, 10)}}}}
	_, err := m.usersCollection.UpdateByID(context.Background(), strconv.FormatInt(userID, 10), update)
	return err
}

type mongoUser struct {
	ID      string `bson:"_id"`
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

	CompletedTasks []string `bson:"completedtasks,omitempty"`
}

func mongoUserFromUser(user *model.User) mongoUser {
	return mongoUser{
		ID:              strconv.FormatInt(user.ID, 10),
		Token:           user.Token,
		Referal:         user.Referal,
		Rk:              user.Rk,
		Avatar:          user.Avatar,
		FirstLogin:      user.FirstLogin,
		LastLogin:       user.LastLogin,
		LastLeave:       user.LastLeave,
		InvitedReferals: user.InvitedReferals,
		RaffleRules:     user.RaffleRules,
		InviteCopy:      user.InviteCopy,
	}
}

type mongoTask struct {
	ID   string `bson:"_id"`
	Name string
}

func mongoTaskFromTask(task model.Task) mongoTask {
	return mongoTask{
		ID:   strconv.FormatInt(task.ID, 10),
		Name: task.Name,
	}
}
