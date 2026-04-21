package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// ConvertTotal 转链请求总数
	ConvertTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shortener_convert_total",
			Help: "Total number of convert (long to short URL) requests",
		},
		[]string{"status"}, // success / error
	)

	// ShowTotal 跳转请求总数
	ShowTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shortener_show_total",
			Help: "Total number of show (redirect) requests",
		},
		[]string{"status"}, // success / not_found / blocked / error
	)

	// BloomFilterHits Bloom Filter 命中统计
	BloomFilterHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shortener_bloom_filter_hits",
			Help: "Bloom filter lookup results",
		},
		[]string{"result"}, // hit / miss
	)

	// LLMLatency LLM 调用延迟
	LLMLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shortener_llm_latency_seconds",
			Help:    "Latency of LLM API calls in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"}, // summarize / generate_slug
	)

	// SafetyBlocked 安全拦截计数
	SafetyBlocked = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "shortener_safety_blocked_total",
			Help: "Total number of requests blocked by safety checks",
		},
	)

	// KafkaProduceTotal Kafka 消息生产总数
	KafkaProduceTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shortener_kafka_produce_total",
			Help: "Total number of Kafka messages produced",
		},
		[]string{"topic", "status"}, // success / error
	)

	// KafkaConsumeTotal Kafka 消息消费总数
	KafkaConsumeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shortener_kafka_consume_total",
			Help: "Total number of Kafka messages consumed",
		},
		[]string{"topic", "status"}, // success / error
	)

	// KafkaProduceLatency Kafka 消息生产延迟
	KafkaProduceLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shortener_kafka_produce_latency_seconds",
			Help:    "Latency of Kafka produce operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)

	// ClickEventTotal 点击事件总数
	ClickEventTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "shortener_click_event_total",
			Help: "Total number of click events received",
		},
	)
)

func init() {
	prometheus.MustRegister(
		ConvertTotal,
		ShowTotal,
		BloomFilterHits,
		LLMLatency,
		SafetyBlocked,
		KafkaProduceTotal,
		KafkaConsumeTotal,
		KafkaProduceLatency,
		ClickEventTotal,
	)
}

