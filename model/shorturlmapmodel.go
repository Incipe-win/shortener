package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ShortUrlMapModel = (*customShortUrlMapModel)(nil)

type (
	// ShortUrlMapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShortUrlMapModel.
	ShortUrlMapModel interface {
		shortUrlMapModel
		FindAll() ([]string, error)
	}

	customShortUrlMapModel struct {
		*defaultShortUrlMapModel
	}
)

func (s *customShortUrlMapModel) FindAll() ([]string, error) {
	query := fmt.Sprintf("select `surl` from %s where `is_del` = 0 and `surl` is not null", s.table)
	type SurlResult struct {
		Surl string `db:"surl"`
	}
	var tempResp []*SurlResult

	err := s.CachedConn.QueryRowsNoCacheCtx(context.Background(), &tempResp, query)
	if err != nil {
		return nil, err
	}

	surls := make([]string, 0, len(tempResp))
	for _, item := range tempResp {
		surls = append(surls, item.Surl)
	}

	return surls, nil
}

// NewShortUrlMapModel returns a model for the database table.
func NewShortUrlMapModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ShortUrlMapModel {
	return &customShortUrlMapModel{
		defaultShortUrlMapModel: newShortUrlMapModel(conn, c, opts...),
	}
}
