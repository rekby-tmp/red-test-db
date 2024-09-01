package database

import (
	"context"
	"fmt"
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
	invite_referals int8,
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
	query := `INSERT INTO users (
                   id, 
                   token,
                   referal,
                   rk, 
                   avatar,
                   first_login,
                   last_login,
                   last_leave,
                   invite_referals,
                   raff_rules,
                   invite_copy
                   )
                VALUES (
                   @id, 
                   @token,
                   @referal,
                   @rk, 
                   @avatar,
                   @first_login,
                   @last_login,
                   @last_leave,
                   @invite_referals,
                   @raff_rules,
                   @invite_copy
                )
                   `

	batch := &pgx.Batch{}
	for i := range users {
		user := &users[i]
		args := pgx.NamedArgs{
			"id":              user.ID,
			"token":           user.Token,
			"referal":         user.Referal,
			"rk":              user.Rk,
			"avatar":          user.Avatar,
			"first_login":     user.FirstLogin,
			"last_login":      user.LastLogin,
			"last_leave":      user.LastLeave,
			"invite_referals": user.InvitedReferals,
			"raff_rules":      user.RaffleRules,
			"invite_copy":     user.InviteCopy,
		}
		batch.Queue(query, args)
	}

	results := p.db.SendBatch(context.Background(), batch)
	for i := range users {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("failed to batch upsert on user %v: %w", users[i].ID, err)
		}
	}
	if err := results.Close(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) CreateUser(user model.User) error {
	return p.UploadUsers([]model.User{user})
}

func (p *Postgres) CreateTask(task model.Task) error {
	_, err := p.db.Exec(context.Background(), `
INSERT INTO tasks (id, name) VALUES (@id, @name)
`, pgx.NamedArgs{"id": task.ID, "name": task.Name})

	return err
}

func (p *Postgres) Login(userID int64, token string) error {
	row := p.db.QueryRow(
		context.Background(),
		`SELECT 1 FROM users WHERE id=@id AND token=@token`,
		pgx.NamedArgs{"id": userID, "token": token},
	)
	var res int
	return row.Scan(&res)
}

func (p *Postgres) ClickInviteFeferals(userID int64) error {
	_, err := p.db.Exec(context.Background(), `UPDATE users SET invite_referals=invite_referals+1 WHERE id=@id`)
	return err
}

func (p *Postgres) CompleteTask(userID, taskID int64) error {
	_, err := p.db.Exec(context.Background(), `
INSERT INTO user_task (user_id, task_id) VALUES (@user_id, @task_id)
ON CONFLICT DO NOTHING
`, pgx.NamedArgs{"user_id": userID, "task_id": taskID})
	return err
}
