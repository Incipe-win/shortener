package mq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// KafkaProducer 封装 Kafka 消息生产者
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer 创建 Kafka Producer
// brokers: Kafka 集群地址列表
func NewKafkaProducer(brokers []string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond, // 低延迟场景，快速刷盘
		RequiredAcks: kafka.RequireOne,       // leader 确认即可
		Async:        false,                  // 同步写入，确保消息可靠
		MaxAttempts:  3,                      // 最大重试次数
	}

	logx.Infof("[Kafka Producer] initialized with brokers: %v", brokers)
	return &KafkaProducer{writer: w}
}

// Send 发送消息到指定 topic
func (p *KafkaProducer) Send(ctx context.Context, topic, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		logx.Errorw("[Kafka Producer] failed to marshal message",
			logx.LogField{Key: "topic", Value: topic},
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	}

	start := time.Now()
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		logx.Errorw("[Kafka Producer] failed to send message",
			logx.LogField{Key: "topic", Value: topic},
			logx.LogField{Key: "key", Value: key},
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}

	logx.Debugf("[Kafka Producer] message sent to topic=%s key=%s latency=%v",
		topic, key, time.Since(start))
	return nil
}

// Close 关闭 Producer
func (p *KafkaProducer) Close() error {
	logx.Info("[Kafka Producer] closing...")
	return p.writer.Close()
}
