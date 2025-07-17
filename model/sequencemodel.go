package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SequenceModel = (*customSequenceModel)(nil)

type (
	// SequenceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSequenceModel.
	SequenceModel interface {
		sequenceModel
		ReplaceInto(ctx context.Context, data *Sequence) (sql.Result, error)
	}

	customSequenceModel struct {
		*defaultSequenceModel
	}
)

func (c *customSequenceModel) ReplaceInto(ctx context.Context, data *Sequence) (sql.Result, error) {
	sqlTestSequenceIdKey := fmt.Sprintf("%s%v", cacheSqlTestSequenceIdPrefix, data.Id)
	sqlTestSequenceStubKey := fmt.Sprintf("%s%v", cacheSqlTestSequenceStubPrefix, data.Stub)
	ret, err := c.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("replace into %s (`stub`) values (?)", c.table)
		return conn.ExecCtx(ctx, query, data.Stub)
	}, sqlTestSequenceIdKey, sqlTestSequenceStubKey)
	return ret, err
}

// NewSequenceModel returns a model for the database table.
func NewSequenceModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SequenceModel {
	return &customSequenceModel{
		defaultSequenceModel: newSequenceModel(conn, c, opts...),
	}
}
