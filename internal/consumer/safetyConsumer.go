package consumer

import (
	"context"
	"encoding/json"

	"shortener/pkg/metrics"
	"shortener/pkg/mq"

	"github.com/zeromicro/go-zero/core/logx"
)

// SafetyAlertHandler 返回安全告警消息的处理函数
// 当前版本将告警记录到日志和 Prometheus 指标
// 后续可扩展为钉钉/邮件/Webhook 通知
func SafetyAlertHandler() mq.MessageHandler {
	return func(ctx context.Context, key, value []byte) error {
		var msg mq.SafetyAlertMessage
		if err := json.Unmarshal(value, &msg); err != nil {
			logx.Errorw("[Safety Consumer] failed to unmarshal message",
				logx.LogField{Key: "err", Value: err.Error()})
			return err
		}

		metrics.KafkaConsumeTotal.WithLabelValues(mq.TopicSafetyAlert, "success").Inc()

		// 根据风险等级使用不同日志级别
		switch msg.RiskLevel {
		case "danger":
			logx.Errorw("[Safety Consumer] 🚨 DANGER alert",
				logx.LogField{Key: "surl", Value: msg.Surl},
				logx.LogField{Key: "long_url", Value: msg.LongUrl},
				logx.LogField{Key: "reason", Value: msg.Reason})
		case "warning":
			logx.Sloww("[Safety Consumer] ⚠️ WARNING alert",
				logx.LogField{Key: "surl", Value: msg.Surl},
				logx.LogField{Key: "long_url", Value: msg.LongUrl},
				logx.LogField{Key: "reason", Value: msg.Reason})
		default:
			logx.Infof("[Safety Consumer] safety alert: surl=%s risk=%s reason=%s",
				msg.Surl, msg.RiskLevel, msg.Reason)
		}

		return nil
	}
}
