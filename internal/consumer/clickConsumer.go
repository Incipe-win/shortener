package consumer

import (
	"context"
	"encoding/json"

	"shortener/pkg/metrics"
	"shortener/pkg/mq"

	"github.com/zeromicro/go-zero/core/logx"
)

// ClickEventHandler 返回点击事件消息的处理函数
// 当前版本将点击事件记录到日志和 Prometheus 指标
// 后续可扩展为写入 MySQL click_count 或 Redis 计数器
func ClickEventHandler() mq.MessageHandler {
	return func(ctx context.Context, key, value []byte) error {
		var msg mq.ClickEventMessage
		if err := json.Unmarshal(value, &msg); err != nil {
			logx.Errorw("[Click Consumer] failed to unmarshal message",
				logx.LogField{Key: "err", Value: err.Error()})
			return err
		}

		// 记录点击事件指标
		metrics.ClickEventTotal.Inc()
		metrics.KafkaConsumeTotal.WithLabelValues(mq.TopicClickEvent, "success").Inc()

		logx.Infof("[Click Consumer] click event: surl=%s ip=%s ua=%s referer=%s",
			msg.Surl, msg.ClientIP, msg.UserAgent, msg.Referer)

		return nil
	}
}
