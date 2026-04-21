package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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
		FindList(ctx context.Context, userID uint64, page, pageSize int, search string) ([]*ShortUrlMap, error)
		Count(ctx context.Context, userID uint64, search string) (int64, error)
		IncrementClickCount(ctx context.Context, surl string) error
		GetStats(ctx context.Context, userID uint64) (map[string]int64, error)
		CountUnregistered(ctx context.Context) (int64, error)
	}

	// ShortUrlMapWithStats extends ShortUrlMap with click_count for queries
	ShortUrlMapWithStats struct {
		Id         uint64
		CreateAt   time.Time
		CreateBy   string
		UserID     sql.NullInt64
		IsDel      uint64
		Lurl       sql.NullString
		Md5        sql.NullString
		Surl       sql.NullString
		AiSummary  sql.NullString
		AiKeywords sql.NullString
		AiSlug     sql.NullString
		RiskLevel  sql.NullString
		RiskReason sql.NullString
		ClickCount uint64
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

// FindList 分页查询链接列表（包含 click_count，按用户过滤）
func (s *customShortUrlMapModel) FindList(ctx context.Context, userID uint64, page, pageSize int, search string) ([]*ShortUrlMap, error) {
	offset := (page - 1) * pageSize

	// 构建用户过滤条件
	userWhere := ""
	var userArgs []interface{}
	if userID > 0 {
		userWhere = " AND `user_id` = ?"
		userArgs = append(userArgs, userID)
	}

	var items []*ShortUrlMapWithStats
	var err error

	if search != "" {
		query := fmt.Sprintf("select `id`,`create_at`,`create_by`,`user_id`,`is_del`,`lurl`,`md5`,`surl`,`ai_summary`,`ai_keywords`,`ai_slug`,`risk_level`,`risk_reason`,`click_count` from %s where `is_del` = 0"+userWhere+" and (`surl` like ? or `lurl` like ?) order by `id` desc limit ?, ?", s.table)
		args := append(userArgs, "%"+search+"%", "%"+search+"%", offset, pageSize)
		err = s.CachedConn.QueryRowsNoCacheCtx(ctx, &items, query, args...)
	} else {
		query := fmt.Sprintf("select `id`,`create_at`,`create_by`,`user_id`,`is_del`,`lurl`,`md5`,`surl`,`ai_summary`,`ai_keywords`,`ai_slug`,`risk_level`,`risk_reason`,`click_count` from %s where `is_del` = 0"+userWhere+" order by `id` desc limit ?, ?", s.table)
		args := append(userArgs, offset, pageSize)
		err = s.CachedConn.QueryRowsNoCacheCtx(ctx, &items, query, args...)
	}
	if err != nil {
		return nil, err
	}

	// 转换为 ShortUrlMap 以兼容现有接口
	resp := make([]*ShortUrlMap, 0, len(items))
	for _, item := range items {
		resp = append(resp, &ShortUrlMap{
			Id:         item.Id,
			CreateAt:   item.CreateAt,
			CreateBy:   item.CreateBy,
			UserID:     item.UserID,
			IsDel:      item.IsDel,
			Lurl:       item.Lurl,
			Md5:        item.Md5,
			Surl:       item.Surl,
			AiSummary:  item.AiSummary,
			AiKeywords: item.AiKeywords,
			AiSlug:     item.AiSlug,
			RiskLevel:  item.RiskLevel,
			RiskReason: item.RiskReason,
			ClickCount: item.ClickCount,
		})
	}
	return resp, nil
}

// Count 统计链接总数（按用户过滤）
func (s *customShortUrlMapModel) Count(ctx context.Context, userID uint64, search string) (int64, error) {
	userWhere := ""
	var userArgs []interface{}
	if userID > 0 {
		userWhere = " AND `user_id` = ?"
		userArgs = append(userArgs, userID)
	}

	var count int64
	var err error

	if search != "" {
		query := fmt.Sprintf("select count(*) from %s where `is_del` = 0"+userWhere+" and (`surl` like ? or `lurl` like ?)", s.table)
		like := "%" + search + "%"
		args := append(userArgs, like, like)
		err = s.CachedConn.QueryRowNoCacheCtx(ctx, &count, query, args...)
	} else {
		query := fmt.Sprintf("select count(*) from %s where `is_del` = 0"+userWhere, s.table)
		err = s.CachedConn.QueryRowNoCacheCtx(ctx, &count, query, userArgs...)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// IncrementClickCount 增加短链接点击次数
func (s *customShortUrlMapModel) IncrementClickCount(ctx context.Context, surl string) error {
	// 清除缓存，确保拿到最新数据
	sqlTestShortUrlMapSurlKey := fmt.Sprintf("%s%v", cacheSqlTestShortUrlMapSurlPrefix,
		sql.NullString{String: surl, Valid: true})

	_, err := s.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `click_count` = `click_count` + 1 where `surl` = ?", s.table)
		return conn.ExecCtx(ctx, query, surl)
	}, sqlTestShortUrlMapSurlKey)
	return err
}

// GetStats 获取仪表盘统计数据（按用户过滤）
func (s *customShortUrlMapModel) GetStats(ctx context.Context, userID uint64) (map[string]int64, error) {
	stats := make(map[string]int64)

	userWhere := ""
	var userArgs []interface{}
	if userID > 0 {
		userWhere = " AND `user_id` = ?"
		userArgs = append(userArgs, userID)
	}

	// 总链接数
	var totalLinks int64
	if err := s.CachedConn.QueryRowNoCacheCtx(ctx, &totalLinks,
		fmt.Sprintf("select count(*) from %s where `is_del` = 0"+userWhere, s.table), userArgs...); err != nil {
		return nil, err
	}
	stats["total_links"] = totalLinks

	// 总点击数
	var totalClicks int64
	if err := s.CachedConn.QueryRowNoCacheCtx(ctx, &totalClicks,
		fmt.Sprintf("select coalesce(sum(`click_count`), 0) from %s where `is_del` = 0"+userWhere, s.table), userArgs...); err != nil {
		return nil, err
	}
	stats["total_clicks"] = totalClicks

	// 今日新增链接
	var todayLinks int64
	if err := s.CachedConn.QueryRowNoCacheCtx(ctx, &todayLinks,
		fmt.Sprintf("select count(*) from %s where `is_del` = 0 and date(`create_at`) = curdate()"+userWhere, s.table), userArgs...); err != nil {
		return nil, err
	}
	stats["today_links"] = todayLinks

	// 今日点击数
	var todayClicks int64
	if err := s.CachedConn.QueryRowNoCacheCtx(ctx, &todayClicks,
		fmt.Sprintf("select coalesce(sum(`click_count`), 0) from %s where `is_del` = 0 and date(`create_at`) = curdate()"+userWhere, s.table), userArgs...); err != nil {
		return nil, err
	}
	stats["today_clicks"] = todayClicks

	// 安全拦截（danger 级别）
	var blocked int64
	if err := s.CachedConn.QueryRowNoCacheCtx(ctx, &blocked,
		fmt.Sprintf("select count(*) from %s where `risk_level` = 'danger'"+userWhere, s.table), userArgs...); err != nil {
		return nil, err
	}
	stats["blocked_count"] = blocked

	return stats, nil
}

// CountUnregistered 统计未注册用户（user_id IS NULL）创建的链接数
func (s *customShortUrlMapModel) CountUnregistered(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where `is_del` = 0 and `user_id` is null", s.table)
	if err := s.CachedConn.QueryRowNoCacheCtx(ctx, &count, query); err != nil {
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
