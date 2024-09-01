package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"red-db-test/model"
)

type Postgres struct {
	db *pgx.Conn
}

func NewPostgres(connString string) (*Postgres, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &Postgres{db: conn}, nil
}

func (p *Postgres) CreateTables() error {
	if _, err := p.db.Exec(context.Background(), `
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS user_task;

CREATE TABLE users (
    id int8 PRIMARY KEY ,
    token Text,
    referal int8,
    rk int8,
    avatar text,
    first_login timestamp,
	last_login timestamp,
	last_leave timestamp,
	invite_feferals int8,
	raff_rules int8,
	invite_copy int8
);

CREATE TABLE tasks (
    id int8 PRIMARY KEY,
    name text
);

CREATE TABLE user_task (
    user_id int8,
    task_id int8,
    PRIMARY KEY (user_id, task_id)
);

CREATE INDEX user_task_task_user ON user_task (task_id, user_id);

`); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UploadUsers(users []model.User) error {
	//TODO implement me
	panic("implement me")
}

func (p *Postgres) CreateUser(user model.User) error {
	//TODO implement me
	panic("implement me")
}

func (p *Postgres) CreateTask(task model.Task) error {
	//TODO implement me
	panic("implement me")
}

func (p *Postgres) Login(userID int64, token string) error {
	//TODO implement me
	panic("implement me")
}

func (p *Postgres) ClickInviteFeferals(userID int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *Postgres) CompleteTask(userID, taskID int64) error {
	//TODO implement me
	panic("implement me")
}
