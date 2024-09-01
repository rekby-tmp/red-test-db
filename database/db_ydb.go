package database

import (
	"context"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	yc "github.com/ydb-platform/ydb-go-yc-metadata"
	"path"
	"red-db-test/model"
	"time"
)

type YDB struct {
	db *ydb.Driver
}

func (y *YDB) CreateTables() error {
	if err := y.db.Query().Exec(context.Background(), `DROP TABLE IF EXISTS users`); err != nil {
		return err
	}
	if err := y.db.Query().Exec(context.Background(), `DROP TABLE IF EXISTS tasks`); err != nil {
		return err
	}
	if err := y.db.Query().Exec(context.Background(), `DROP TABLE IF EXISTS user_task`); err != nil {
		return err
	}

	if err := y.db.Query().Exec(context.Background(), `
CREATE TABLE users (
    id Int64 NOT NULL,
    token Text,
    referal Int64,
    rk Int64,
    avatar Text,
    first_login Timestamp,
	last_login Timestamp,
	last_leave Timestamp,
	invite_referals Int64,
	raff_rules Int64,
	invite_copy Int64,
	PRIMARY KEY (id)
)
`); err != nil {
		return err
	}

	if err := y.db.Query().Exec(context.Background(), `
CREATE TABLE tasks (
    id Int64 NOT NULL,
    name Text,
	PRIMARY KEY (id)
)
`); err != nil {
		return err
	}

	if err := y.db.Query().Exec(context.Background(), `
CREATE TABLE user_task (
    user_id Int64 NOT NULL,
    task_id Int64 NOT NULL,
    PRIMARY KEY (user_id, task_id),
    INDEX user_task__task_user GLOBAL ON (task_id, user_id)
);
`); err != nil {
		return err
	}

	return nil
}

func (y *YDB) UploadUsers(users []model.User) error {
	if len(users) == 0 {
		return nil
	}

	rows := make([]types.Value, len(users))
	for i := range users {
		user := &users[i]
		rows[i] = types.StructValue(
			types.StructFieldValue("id", types.Int64Value(user.ID)),
			types.StructFieldValue("token", types.TextValue(user.Token)),
			types.StructFieldValue("referal", types.Int64Value(user.Referal)),
			types.StructFieldValue("rk", types.Int64Value(user.Rk)),
			types.StructFieldValue("avatar", types.TextValue(user.Avatar)),
			types.StructFieldValue("first_login", types.TimestampValueFromTime(user.FirstLogin)),
			types.StructFieldValue("last_login", types.TimestampValueFromTime(user.LastLogin)),
			types.StructFieldValue("last_leave", types.TimestampValueFromTime(user.LastLeave)),
			types.StructFieldValue("invite_referals", types.Int64Value(user.InvitedReferals)),
			types.StructFieldValue("raff_rules", types.Int64Value(user.RaffleRules)),
			types.StructFieldValue("invite_copy", types.Int64Value(user.InviteCopy)),
		)
	}
	tablePath := path.Join(y.db.Name(), "users")
	return y.db.Table().Do(context.Background(), func(ctx context.Context, s table.Session) error {
		return s.BulkUpsert(ctx, tablePath, types.ListValue(rows...))
	}, table.WithIdempotent())
}

func (y *YDB) CreateUser(user model.User) error {
	return y.db.Query().Exec(context.Background(), `
DECLARE $id AS Int64 NOT NULL;
DECLARE $token AS Text NOT NULL;
DECLARE $referal AS Int64 NOT NULL;
DECLARE $rk AS Int64 NOT NULL;
DECLARE $avatar AS Text NOT NULL;
DECLARE $first_login AS Timestamp NOT NULL;
DECLARE $last_login AS Timestamp NOT NULL;
DECLARE $last_leave AS Timestamp NOT NULL;
DECLARE $invite_referals AS Int64 NOT NULL;
DECLARE $raff_rules AS Int64 NOT NULL;
DECLARE $invite_copy AS Int64 NOT NULL;

INSERT INTO users (
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
	invite_copy,
) VALUES (
    $id,
    $token,
    $referal,
    $rk,
    $avatar,
    $first_login,
	$last_login,
	$last_leave,
	$invite_referals,
	$raff_rules,
	$invite_copy,
);
`,
		query.WithParameters(
			ydb.ParamsBuilder().
				Param("$id").Int64(user.ID).
				Param("$token").Text(user.Token).
				Param("$referal").Int64(user.Referal).
				Param("$rk").Int64(user.Rk).
				Param("$avatar").Text(user.Avatar).
				Param("$first_login").Timestamp(user.FirstLogin).
				Param("$last_login").Timestamp(user.LastLogin).
				Param("$last_leave").Timestamp(user.LastLeave).
				Param("$invite_referals").Int64(user.InvitedReferals).
				Param("$raff_rules").Int64(user.RaffleRules).
				Param("$invite_copy").Int64(user.InviteCopy).
				Build(),
		),
	)
}

func (y *YDB) CreateTask(task model.Task) error {
	return y.db.Query().Exec(context.Background(), `
DECLARE $id AS Int64;
DECLARE $name AS Text;

UPSERT INTO tasks (id, name) VALUES ($id, $name);
`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$id").Int64(task.ID).
			Param("$name").Text(task.Name).
			Build(),
		),
		query.WithIdempotent(),
	)
}

func (y *YDB) Login(userID int64, token string) error {
	q := `
DECLARE $user_id AS Int64;
DECLARE $token AS Text;
DECLARE $now AS Timestamp;

$users = SELECT id FROM users WHERE id=$user_id AND token=$token;

SELECT COUNT(*) FROM $users;

UPSERT INTO users
SELECT id, $now AS last_login FROM $users;
`
	row, err := y.db.Query().QueryRow(context.Background(), q, query.WithIdempotent(),
		query.WithParameters(
			ydb.ParamsBuilder().
				Param("$user_id").Int64(userID).
				Param("$token").Text(token).
				Param("$now").Timestamp(time.Now()).
				Build(),
		))
	if err != nil {
		return err
	}
	var count uint64
	if err = row.Scan(&count); err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("auth failed")
	}
	return nil
}

func (y *YDB) ClickInviteFeferals(userID int64) error {
	q := `
DECLARE $user_id AS Int64;

UPDATE users SET invite_referals=invite_referals+1 WHERE id=$user_id;
`
	return y.db.Query().Exec(context.Background(), q, query.WithParameters(
		ydb.ParamsBuilder().Param("$user_id").Int64(userID).Build(),
	))
}

func (y *YDB) CompleteTask(userID, taskID int64) error {
	q := `
DECLARE $user_id AS Int64;
DECLARE $task_id AS Int64;

UPSERT INTO user_task (user_id, task_id) VALUES ($user_id, $task_id);
`
	return y.db.Query().Exec(context.Background(), q, query.WithIdempotent(),
		query.WithParameters(
			ydb.ParamsBuilder().
				Param("$user_id").Int64(userID).
				Param("$task_id").Int64(taskID).
				Build(),
		),
	)
}

func NewYDB(connectionString, token string, metadataToken bool) (*YDB, error) {
	var authOption ydb.Option
	if metadataToken {
		authOption = yc.WithCredentials()
	} else {
		authOption = ydb.WithAccessTokenCredentials(token)
	}
	db, err := ydb.Open(context.Background(), connectionString, authOption, yc.WithInternalCA())
	if err != nil {
		return nil, err
	}
	return &YDB{db: db}, nil
}
