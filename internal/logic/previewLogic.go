package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/pkg/otel"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewLogic {
	return &PreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Preview 返回短链的 AI 分析预览信息（摘要、关键词、风险等级）
func (l *PreviewLogic) Preview(req *types.PreviewRequest) (resp *types.PreviewResponse, err error) {
	_, span := otel.Tracer().Start(l.ctx, "PreviewLogic.Preview")
	defer span.End()

	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("short URL not found")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed",
			logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	// 解析关键词 JSON 数组
	keywords := parseKeywords(u.AiKeywords.String)

	return &types.PreviewResponse{
		ShortUrl:  req.ShortUrl,
		LongUrl:   u.Lurl.String,
		Summary:   u.AiSummary.String,
		Keywords:  keywords,
		RiskLevel: u.RiskLevel.String,
	}, nil
}

// parseKeywords 将 JSON 格式的关键词字符串解析为 []string
func parseKeywords(raw string) []string {
	if raw == "" {
		return []string{}
	}
	var result []string
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return []string{}
	}
	return result
}
