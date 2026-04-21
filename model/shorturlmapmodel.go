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
		FindList(ctx context.Context, page, pageSize int, search string) ([]*ShortUrlMap, error)
		Count(ctx context.Context, search string) (int64, error)
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

// FindList 分页查询链接列表
func (s *customShortUrlMapModel) FindList(ctx context.Context, page, pageSize int, search string) ([]*ShortUrlMap, error) {
	offset := (page - 1) * pageSize
	var resp []*ShortUrlMap
	var err error

	if search != "" {
		query := fmt.Sprintf("select %s from %s where `is_del` = 0 and (`surl` like ? or `lurl` like ?) order by `id` desc limit ?, ?", shortUrlMapRows, s.table)
		like := "%" + search + "%"
		err = s.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, like, like, offset, pageSize)
	} else {
		query := fmt.Sprintf("select %s from %s where `is_del` = 0 order by `id` desc limit ?, ?", shortUrlMapRows, s.table)
		err = s.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, offset, pageSize)
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Count 统计链接总数
func (s *customShortUrlMapModel) Count(ctx context.Context, search string) (int64, error) {
	var count int64
	var err error

	if search != "" {
		query := fmt.Sprintf("select count(*) from %s where `is_del` = 0 and (`surl` like ? or `lurl` like ?)", s.table)
		like := "%" + search + "%"
		err = s.CachedConn.QueryRowNoCacheCtx(ctx, &count, query, like, like)
	} else {
		query := fmt.Sprintf("select count(*) from %s where `is_del` = 0", s.table)
		err = s.CachedConn.QueryRowNoCacheCtx(ctx, &count, query)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// NewShortUrlMapModel returns a model for the database table.
func NewShortUrlMapModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ShortUrlMapModel {
	return &customShortUrlMapModel{
		defaultShortUrlMapModel: newShortUrlMapModel(conn, c, opts...),
	}
}
