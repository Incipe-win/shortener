package mq

// Topic 常量
const (
	TopicAIAnalysis  = "ai-analysis"
	TopicClickEvent  = "click-events"
	TopicSafetyAlert = "safety-alerts"
)

// AIAnalysisMessage AI 分析任务消息
type AIAnalysisMessage struct {
	Surl    string `json:"surl"`
	LongUrl string `json:"long_url"`
}

// ClickEventMessage 链接点击事件消息
type ClickEventMessage struct {
	Surl      string `json:"surl"`
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	Timestamp int64  `json:"timestamp"`
}

// SafetyAlertMessage 安全告警消息
type SafetyAlertMessage struct {
	Surl      string `json:"surl"`
	LongUrl   string `json:"long_url"`
	RiskLevel string `json:"risk_level"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}
