package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlDB struct {
		DSN string
	}

	Sequence struct {
		DSN string
	}

	CacheRedis cache.CacheConf

	BaseString string // bas62指定基础字符串

	ShortUrlBlackList []string
	ShortDoamin       string

	// LLM 大模型配置
	LLM struct {
		Enabled bool
		BaseURL string // API 地址，如 https://api.openai.com/v1
		APIKey  string
		Model   string // 如 gpt-4o-mini
	}

	// Safety 安全巡检配置
	Safety struct {
		Enabled          bool
		BlackListDomains []string // 黑名单域名列表
	}

	// Otel 可观测性配置（避免与 RestConf 内置 Telemetry 冲突）
	Otel struct {
		Name     string
		Endpoint string
		Sampler  float64
	}

	// Kafka 消息队列配置
	Kafka struct {
		Enabled bool
		Brokers []string
		Topics  struct {
			AIAnalysis  string `json:",default=ai-analysis"`
			ClickEvent  string `json:",default=click-events"`
			SafetyAlert string `json:",default=safety-alerts"`
		}
	}
}
