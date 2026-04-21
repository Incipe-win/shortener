package consumer

import (
	"context"
	"encoding/json"
	"time"

	"shortener/model"
	"shortener/pkg/llm"
	"shortener/pkg/metrics"
	"shortener/pkg/mq"
	"shortener/pkg/scraper"

	"github.com/zeromicro/go-zero/core/logx"
)

// AIAnalysisHandler 返回 AI 分析消息的处理函数
func AIAnalysisHandler(llmClient *llm.Client, shortUrlModel model.ShortUrlMapModel) mq.MessageHandler {
	return func(ctx context.Context, key, value []byte) error {
		var msg mq.AIAnalysisMessage
		if err := json.Unmarshal(value, &msg); err != nil {
			logx.Errorw("[AI Consumer] failed to unmarshal message",
				logx.LogField{Key: "err", Value: err.Error()})
			return err
		}

		logx.Infof("[AI Consumer] processing surl=%s lurl=%s", msg.Surl, msg.LongUrl)

		// 使用独立超时 context
		analysisCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		start := time.Now()
		defer func() {
			metrics.LLMLatency.WithLabelValues("summarize").Observe(time.Since(start).Seconds())
		}()

		// Step 1: 抓取页面内容
		content, err := scraper.FetchPageContent(analysisCtx, msg.LongUrl)
		if err != nil {
			logx.Errorw("[AI Consumer] failed to fetch page content",
				logx.LogField{Key: "surl", Value: msg.Surl},
				logx.LogField{Key: "err", Value: err.Error()})
			return err
		}

		// Step 2: LLM 分析
		analysis, err := llmClient.Summarize(analysisCtx, content)
		if err != nil {
			logx.Errorw("[AI Consumer] LLM summarize failed",
				logx.LogField{Key: "surl", Value: msg.Surl},
				logx.LogField{Key: "err", Value: err.Error()})
			return err
		}

		// Step 3: 判定风险等级
		riskLevel := "safe"
		if analysis.RiskScore >= 0.7 {
			riskLevel = "danger"
		} else if analysis.RiskScore >= 0.3 {
			riskLevel = "warning"
		}

		// Step 4: 更新数据库 AI 字段
		err = shortUrlModel.UpdateAIFields(analysisCtx, msg.Surl, analysis.Summary,
			analysis.Keywords, analysis.Slug, riskLevel, analysis.RiskReason)
		if err != nil {
			logx.Errorw("[AI Consumer] failed to update AI fields",
				logx.LogField{Key: "surl", Value: msg.Surl},
				logx.LogField{Key: "err", Value: err.Error()})
			return err
		}

		metrics.KafkaConsumeTotal.WithLabelValues(mq.TopicAIAnalysis, "success").Inc()
		logx.Infof("[AI Consumer] analysis completed for surl=%s risk_level=%s", msg.Surl, riskLevel)
		return nil
	}
}
