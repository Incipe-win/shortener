package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shortener/internal/ctxdata"
	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/pkg/metrics"
	"shortener/pkg/mq"
	"shortener/pkg/otel"

	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 创建 tracing span
	ctx, span := otel.Tracer().Start(l.ctx, "ShowLogic.Show")
	defer span.End()

	// Bloom Filter 前置检查
	exist, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
	if err != nil {
		logx.Errorw("bloom filter check failed", logx.LogField{Value: err.Error(), Key: "err"})
		metrics.ShowTotal.WithLabelValues("error").Inc()
		return nil, err
	}
	if !exist {
		metrics.BloomFilterHits.WithLabelValues("miss").Inc()
		metrics.ShowTotal.WithLabelValues("not_found").Inc()
		return nil, errors.New("404")
	}
	metrics.BloomFilterHits.WithLabelValues("hit").Inc()

	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			metrics.ShowTotal.WithLabelValues("not_found").Inc()
			return nil, errors.New("short URL not found")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.LogField{Value: err.Error(), Key: "err"})
		metrics.ShowTotal.WithLabelValues("error").Inc()
		return nil, err
	}

	// 【新增】安全等级检查
	span.SetAttributes(attribute.String("risk_level", u.RiskLevel.String))
	if u.RiskLevel.Valid && u.RiskLevel.String == "danger" {
		metrics.SafetyBlocked.Inc()
		metrics.ShowTotal.WithLabelValues("blocked").Inc()
		logx.Infow("redirect blocked due to danger risk level",
			logx.LogField{Key: "surl", Value: req.ShortUrl},
			logx.LogField{Key: "risk_reason", Value: u.RiskReason.String})
		return nil, errors.New("this link has been flagged as potentially unsafe and cannot be accessed")
	}

	metrics.ShowTotal.WithLabelValues("success").Inc()

	// 【新增】发送点击事件到 Kafka
	l.sendClickEvent(req.ShortUrl)

	// 构建响应，附带安全警告信息（供 handler 层使用）
	result := &types.ShowResponse{
		LongUrl: u.Lurl.String,
	}

	// 如果是 warning 级别，通过 RiskWarning 字段传递警告信息
	if u.RiskLevel.Valid && u.RiskLevel.String == "warning" {
		result.RiskWarning = u.RiskReason.String
	}

	return result, nil
}

// sendClickEvent 向 Kafka 发送点击事件消息
func (l *ShowLogic) sendClickEvent(surl string) {
	if l.svcCtx.KafkaProducer == nil {
		return
	}

	// 从 context 中提取 HTTP 请求信息
	clientIP, _ := l.ctx.Value(ctxdata.KeyClientIP).(string)
	userAgent, _ := l.ctx.Value(ctxdata.KeyUserAgent).(string)
	referer, _ := l.ctx.Value(ctxdata.KeyReferer).(string)

	topic := l.svcCtx.Config.Kafka.Topics.ClickEvent
	msg := mq.ClickEventMessage{
		Surl:      surl,
		ClientIP:  clientIP,
		UserAgent: userAgent,
		Referer:   referer,
		Timestamp: time.Now().Unix(),
	}

	if err := l.svcCtx.KafkaProducer.Send(l.ctx, topic, surl, msg); err != nil {
		logx.Errorw("failed to send click event to Kafka",
			logx.LogField{Key: "surl", Value: surl},
			logx.LogField{Key: "err", Value: err.Error()})
		metrics.KafkaProduceTotal.WithLabelValues(topic, "error").Inc()
		return
	}
	metrics.KafkaProduceTotal.WithLabelValues(topic, "success").Inc()
}
