package mq

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// MessageHandler 消息处理函数签名
type MessageHandler func(ctx context.Context, key, value []byte) error

// KafkaConsumer 封装 Kafka 消息消费者
// 实现 go-zero 的 Starter 接口，可注册到 ServiceGroup
type KafkaConsumer struct {
	reader  *kafka.Reader
	handler MessageHandler
	topic   string
}

// NewKafkaConsumer 创建 Kafka Consumer
func NewKafkaConsumer(brokers []string, topic, groupID string, handler MessageHandler) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,          // 最小拉取字节数
		MaxBytes: 10e6,       // 最大 10MB
		MaxWait:  250e6,      // 最多等 250ms
	})

	logx.Infof("[Kafka Consumer] initialized for topic=%s group=%s brokers=%v", topic, groupID, brokers)
	return &KafkaConsumer{
		reader:  r,
		handler: handler,
		topic:   topic,
	}
}

// Start 开始消费消息（阻塞，直到 ctx 取消）
// 实现 go-zero ServiceGroup 的 Starter 接口
func (c *KafkaConsumer) Start() {
	logx.Infof("[Kafka Consumer] starting consumer for topic=%s", c.topic)

	ctx := context.Background()
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			// reader 被关闭时会返回错误，正常退出
			logx.Infof("[Kafka Consumer] topic=%s fetch stopped: %v", c.topic, err)
			return
		}

		if err := c.handler(ctx, msg.Key, msg.Value); err != nil {
			logx.Errorw("[Kafka Consumer] handler failed",
				logx.LogField{Key: "topic", Value: c.topic},
				logx.LogField{Key: "partition", Value: msg.Partition},
				logx.LogField{Key: "offset", Value: msg.Offset},
				logx.LogField{Key: "err", Value: err.Error()})
			// 处理失败不提交 offset，下次重试
			continue
		}

		// 处理成功，提交 offset
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			logx.Errorw("[Kafka Consumer] failed to commit offset",
				logx.LogField{Key: "topic", Value: c.topic},
				logx.LogField{Key: "err", Value: err.Error()})
		}
	}
}

// Stop 停止消费（关闭 reader）
func (c *KafkaConsumer) Stop() {
	logx.Infof("[Kafka Consumer] stopping consumer for topic=%s", c.topic)
	if err := c.reader.Close(); err != nil {
		logx.Errorw("[Kafka Consumer] failed to close reader",
			logx.LogField{Key: "topic", Value: c.topic},
			logx.LogField{Key: "err", Value: err.Error()})
	}
}
