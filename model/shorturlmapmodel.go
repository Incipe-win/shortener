package model

import (
	"context"
	"database/sql"
	"encoding/json"
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
		UpdateAIFields(ctx context.Context, surl string, summary string, keywords []string, slug string, riskLevel string, riskReason string) error
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

// UpdateAIFields 更新 AI 分析相关字段
func (s *customShortUrlMapModel) UpdateAIFields(ctx context.Context, surl string, summary string, keywords []string, slug string, riskLevel string, riskReason string) error {
	// 将 keywords 序列化为 JSON
	keywordsJSON, err := json.Marshal(keywords)
	if err != nil {
		return fmt.Errorf("failed to marshal keywords: %w", err)
	}

	// 清除该 surl 的缓存
	sqlTestShortUrlMapSurlKey := fmt.Sprintf("%s%v", cacheSqlTestShortUrlMapSurlPrefix,
		sql.NullString{String: surl, Valid: true})

	_, err = s.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `ai_summary` = ?, `ai_keywords` = ?, `ai_slug` = ?, `risk_level` = ?, `risk_reason` = ? where `surl` = ?", s.table)
		return conn.ExecCtx(ctx, query, summary, string(keywordsJSON), slug, riskLevel, riskReason, surl)
	}, sqlTestShortUrlMapSurlKey)

	return err
}

// NewShortUrlMapModel returns a model for the database table.
func NewShortUrlMapModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ShortUrlMapModel {
	return &customShortUrlMapModel{
		defaultShortUrlMapModel: newShortUrlMapModel(conn, c, opts...),
	}
}
