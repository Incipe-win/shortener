package consumer

import (
	"context"
	"encoding/json"

	"shortener/model"
	"shortener/pkg/metrics"
	"shortener/pkg/mq"

	"github.com/zeromicro/go-zero/core/logx"
)

// ClickEventHandler 返回点击事件消息的处理函数
func ClickEventHandler(shortUrlModel model.ShortUrlMapModel) mq.MessageHandler {
	return func(ctx context.Context, key, value []byte) error {
		var msg mq.ClickEventMessage
		if err := json.Unmarshal(value, &msg); err != nil {
			logx.Errorw("[Click Consumer] failed to unmarshal message",
				logx.LogField{Key: "err", Value: err.Error()})
			metrics.KafkaConsumeTotal.WithLabelValues(mq.TopicClickEvent, "error").Inc()
			return err
		}

		// 记录点击事件指标
		metrics.ClickEventTotal.Inc()
		metrics.KafkaConsumeTotal.WithLabelValues(mq.TopicClickEvent, "success").Inc()

		logx.Infof("[Click Consumer] click event: surl=%s ip=%s ua=%s referer=%s",
			msg.Surl, msg.ClientIP, msg.UserAgent, msg.Referer)

		// 持久化点击计数到数据库
		if err := shortUrlModel.IncrementClickCount(ctx, msg.Surl); err != nil {
			logx.Errorw("[Click Consumer] failed to increment click count",
				logx.LogField{Key: "surl", Value: msg.Surl},
				logx.LogField{Key: "err", Value: err.Error()})
		}

		return nil
	}
}
