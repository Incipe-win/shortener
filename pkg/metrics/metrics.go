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
)

func init() {
	prometheus.MustRegister(
		ConvertTotal,
		ShowTotal,
		BloomFilterHits,
		LLMLatency,
		SafetyBlocked,
	)
}
