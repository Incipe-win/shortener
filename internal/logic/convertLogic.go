package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/llm"
	"shortener/pkg/md5"
	"shortener/pkg/metrics"
	"shortener/pkg/mq"
	"shortener/pkg/otel"
	"shortener/pkg/safety"
	"shortener/pkg/scraper"
	"shortener/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.opentelemetry.io/otel/attribute"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 转链：输一个长链接 --> 转为短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 创建 tracing span
	ctx, span := otel.Tracer().Start(l.ctx, "ConvertLogic.Convert")
	defer span.End()

	// 1. 校验数据（使用validator）
	// 1.1 数据不能空
	// 1.2 输入的长链接必须是一能请求通的网址
	if ok := connect.Get(req.LongUrl); !ok {
		metrics.ConvertTotal.WithLabelValues("error").Inc()
		return nil, errors.New("invalid long URL")
	}

	// 1.2.5 【新增】域名层面防循环检测
	if urltool.IsCircularURL(req.LongUrl, l.svcCtx.Config.ShortDomain) {
		metrics.ConvertTotal.WithLabelValues("error").Inc()
		return nil, errors.New("cannot convert a URL from this service (circular reference)")
	}

	// 1.2.6 【新增】安全巡检
	if l.svcCtx.Config.Safety.Enabled {
		safetyResult := safety.CheckURL(ctx, req.LongUrl, l.svcCtx.BlackListDomains, nil, "")
		if safetyResult.RiskLevel == "danger" {
			metrics.SafetyBlocked.Inc()
			metrics.ConvertTotal.WithLabelValues("blocked").Inc()
			span.SetAttributes(attribute.String("safety.risk_level", "danger"))

			// 发送安全告警到 Kafka
			l.sendSafetyAlert(req.LongUrl, "", safetyResult.RiskLevel, safetyResult.Reason)

			return nil, fmt.Errorf("URL blocked by safety check: %s", safetyResult.Reason)
		}
		span.SetAttributes(attribute.String("safety.risk_level", safetyResult.RiskLevel))

		// warning 级别也发送告警
		if safetyResult.RiskLevel == "warning" {
			l.sendSafetyAlert(req.LongUrl, "", safetyResult.RiskLevel, safetyResult.Reason)
		}
	}

	// 1.3 判断之前是否已经转链过（数据库中是否已经存在该长链接）
	// 1.3.1 给长链接生成 md5
	md5Hash := md5.Sum([]byte(req.LongUrl))
	// 1.3.2 拿 md5 去数据库中查是否存在
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(ctx, sql.NullString{String: md5Hash, Valid: true})
	if err == nil {
		metrics.ConvertTotal.WithLabelValues("success").Inc()
		return &types.ConvertResponse{
			ShortUrl: l.svcCtx.Config.ShortDomain + "/" + u.Surl.String,
		}, nil
	}
	if err != sqlx.ErrNotFound {
		logx.Errorw("failed to find long URL by md5", logx.LogField{Key: "error", Value: err.Error()})
		metrics.ConvertTotal.WithLabelValues("error").Inc()
		return nil, err
	}
	// 1.4 输入不能是一个短链接（避免循环转链）
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool.GetBasePath failed", logx.LogField{Key: "lurl", Value: req.LongUrl}, logx.LogField{Key: "err", Value: err.Error()})
		metrics.ConvertTotal.WithLabelValues("error").Inc()
		return nil, err
	}
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			metrics.ConvertTotal.WithLabelValues("error").Inc()
			return nil, errors.New("url already is a short link")
		}
		logx.Errorw("failed to find short URL by base path", logx.LogField{Key: "error", Value: err.Error()})
		metrics.ConvertTotal.WithLabelValues("error").Inc()
		return nil, err
	}
	var shortUrl string
	for {
		// 2. 取号
		// 每来一个转链请求，使用 replace into 语句往 sequence 表中插入一条数据，并取出主键 id 作为号码
		res, err := l.svcCtx.SequenceModel.ReplaceInto(ctx, &model.Sequence{Stub: "a"})
		if err != nil {
			logx.Errorw("failed to replace into sequence", logx.LogField{Key: "error", Value: err.Error()})
			metrics.ConvertTotal.WithLabelValues("error").Inc()
			return nil, err
		}
		seq, err := res.LastInsertId()
		if err != nil {
			logx.Errorw("failed to get last insert id from sequence", logx.LogField{Key: "error", Value: err.Error()})
			metrics.ConvertTotal.WithLabelValues("error").Inc()
			return nil, err
		}
		// 3. 号码转短链
		// 3.1 安全性，打乱 basestring 顺序
		shortUrl = base62.Int2String(uint64(seq))
		// 3.2 短域名黑名单，避免某些特殊词比如 api、fuck等
		if _, ok := l.svcCtx.ShortUrlBlackList[shortUrl]; !ok {
			break
		}
	}
	// 4. 存储长链接短链接映射关系
	logx.Debugf("short URL generated: %s", shortUrl)
	if _, err := l.svcCtx.ShortUrlModel.Insert(ctx, &model.ShortUrlMap{
		Lurl: sql.NullString{String: req.LongUrl, Valid: true},
		Md5:  sql.NullString{String: md5Hash, Valid: true},
		Surl: sql.NullString{String: shortUrl, Valid: true},
	}); err != nil {
		logx.Errorw("failed to insert short URL map", logx.LogField{Key: "error", Value: err.Error()})
	}

	// 4.1 将短链接存入 bloom filter
	if err := l.svcCtx.Filter.Add([]byte(shortUrl)); err != nil {
		logx.Errorw("failed to add short URL to bloom filter", logx.LogField{Key: "error", Value: err.Error()})
		metrics.ConvertTotal.WithLabelValues("error").Inc()
		return nil, err
	}

	// 4.2 异步 AI 分析：优先走 Kafka，降级走 goroutine
	if l.svcCtx.LLMClient != nil {
		if l.svcCtx.KafkaProducer != nil {
			// 优先通过 Kafka 发送 AI 分析任务
			l.sendAIAnalysisMessage(shortUrl, req.LongUrl)
		} else {
			// Kafka 未启用，降级走 goroutine
			surl := shortUrl
			longUrl := req.LongUrl
			llmClient := l.svcCtx.LLMClient
			shortUrlModel := l.svcCtx.ShortUrlModel
			go asyncAIAnalysis(surl, longUrl, llmClient, shortUrlModel)
		}
	}

	// 5. 返回响应
	// 5.1 返回的是 短域名+短链接  q1mi.cn/1En
	metrics.ConvertTotal.WithLabelValues("success").Inc()
	shortUrl = l.svcCtx.Config.ShortDomain + "/" + shortUrl
	return &types.ConvertResponse{
		ShortUrl: shortUrl,
	}, nil
}

// sendAIAnalysisMessage 向 Kafka 发送 AI 分析任务消息
func (l *ConvertLogic) sendAIAnalysisMessage(surl, longUrl string) {
	topic := l.svcCtx.Config.Kafka.Topics.AIAnalysis
	msg := mq.AIAnalysisMessage{
		Surl:    surl,
		LongUrl: longUrl,
	}

	start := time.Now()
	if err := l.svcCtx.KafkaProducer.Send(l.ctx, topic, surl, msg); err != nil {
		logx.Errorw("failed to send AI analysis message to Kafka",
			logx.LogField{Key: "surl", Value: surl},
			logx.LogField{Key: "err", Value: err.Error()})
		metrics.KafkaProduceTotal.WithLabelValues(topic, "error").Inc()

		// Kafka 发送失败，降级走 goroutine
		llmClient := l.svcCtx.LLMClient
		shortUrlModel := l.svcCtx.ShortUrlModel
		go asyncAIAnalysis(surl, longUrl, llmClient, shortUrlModel)
		return
	}
	metrics.KafkaProduceTotal.WithLabelValues(topic, "success").Inc()
	metrics.KafkaProduceLatency.WithLabelValues(topic).Observe(time.Since(start).Seconds())
}

// sendSafetyAlert 向 Kafka 发送安全告警消息
func (l *ConvertLogic) sendSafetyAlert(longUrl, surl, riskLevel, reason string) {
	if l.svcCtx.KafkaProducer == nil {
		return
	}

	topic := l.svcCtx.Config.Kafka.Topics.SafetyAlert
	msg := mq.SafetyAlertMessage{
		Surl:      surl,
		LongUrl:   longUrl,
		RiskLevel: riskLevel,
		Reason:    reason,
		Timestamp: time.Now().Unix(),
	}

	if err := l.svcCtx.KafkaProducer.Send(l.ctx, topic, longUrl, msg); err != nil {
		logx.Errorw("failed to send safety alert to Kafka",
			logx.LogField{Key: "err", Value: err.Error()})
		metrics.KafkaProduceTotal.WithLabelValues(topic, "error").Inc()
		return
	}
	metrics.KafkaProduceTotal.WithLabelValues(topic, "success").Inc()
}

// asyncAIAnalysis 异步执行 AI 页面分析（降级方案，Kafka 未启用时使用）
func asyncAIAnalysis(surl, longUrl string, llmClient *llm.Client, shortUrlModel model.ShortUrlMapModel) {
	// 使用独立 context，不受请求 context 生命周期影响
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	start := time.Now()
	defer func() {
		metrics.LLMLatency.WithLabelValues("summarize").Observe(time.Since(start).Seconds())
	}()

	// Step 1: 抓取页面内容
	content, err := scraper.FetchPageContent(ctx, longUrl)
	if err != nil {
		logx.Errorw("async AI: failed to fetch page content",
			logx.LogField{Key: "surl", Value: surl},
			logx.LogField{Key: "err", Value: err.Error()})
		return
	}

	// Step 2: LLM 分析
	analysis, err := llmClient.Summarize(ctx, content)
	if err != nil {
		logx.Errorw("async AI: LLM summarize failed",
			logx.LogField{Key: "surl", Value: surl},
			logx.LogField{Key: "err", Value: err.Error()})
		return
	}

	// Step 3: 更新数据库 AI 字段
	riskLevel := "safe"
	if analysis.RiskScore >= 0.7 {
		riskLevel = "danger"
	} else if analysis.RiskScore >= 0.3 {
		riskLevel = "warning"
	}

	err = shortUrlModel.UpdateAIFields(ctx, surl, analysis.Summary,
		analysis.Keywords, analysis.Slug, riskLevel, analysis.RiskReason)
	if err != nil {
		logx.Errorw("async AI: failed to update AI fields",
			logx.LogField{Key: "surl", Value: surl},
			logx.LogField{Key: "err", Value: err.Error()})
		return
	}

	logx.Infof("async AI analysis completed for surl=%s, risk_level=%s", surl, riskLevel)
}
