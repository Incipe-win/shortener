package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	User struct {
		Id           uint64    `db:"id"`
		CreateAt     time.Time `db:"create_at"`
		Username     string    `db:"username"`
		PasswordHash string    `db:"password_hash"`
		IsDel        uint64    `db:"is_del"`
	}

	UserModel interface {
		Insert(ctx context.Context, data *User) (sql.Result, error)
		FindOneByUsername(ctx context.Context, username string) (*User, error)
	}

	customUserModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

func newUserModel(conn sqlx.SqlConn) *customUserModel {
	return &customUserModel{
		conn:  conn,
		table: "`user`",
	}
}

func (m *customUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`username`,`password_hash`) values (?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.Username, data.PasswordHash)
}

func (m *customUserModel) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	var resp User
	query := fmt.Sprintf("select `id`,`create_at`,`username`,`password_hash`,`is_del` from %s where `username` = ? limit 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &resp, nil
}

// NewUserModel returns a model for the user table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		conn:  conn,
		table: "`user`",
	}
}
