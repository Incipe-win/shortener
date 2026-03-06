package safety

import (
	"context"
	"net"
	"net/url"
	"strings"

	"shortener/pkg/llm"

	"github.com/zeromicro/go-zero/core/logx"
)

// SafetyResult 表示安全检查的结果
type SafetyResult struct {
	IsSafe    bool    `json:"is_safe"`
	RiskLevel string  `json:"risk_level"` // safe / warning / danger
	RiskScore float64 `json:"risk_score"`
	Reason    string  `json:"reason"`
}

// CheckURL 对目标 URL 进行多层安全检查
// Step 1: 黑名单域名匹配
// Step 2: URL 特征规则检查（IP 直连、非标端口、可疑特征）
// Step 3: LLM 深度风险评估（可选）
func CheckURL(ctx context.Context, targetURL string, blackListDomains map[string]struct{}, llmClient *llm.Client, pageContent string) *SafetyResult {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return &SafetyResult{
			IsSafe:    false,
			RiskLevel: "danger",
			RiskScore: 1.0,
			Reason:    "URL 解析失败: " + err.Error(),
		}
	}

	host := strings.ToLower(parsedURL.Hostname())

	// Step 1: 黑名单域名匹配
	if _, ok := blackListDomains[host]; ok {
		return &SafetyResult{
			IsSafe:    false,
			RiskLevel: "danger",
			RiskScore: 1.0,
			Reason:    "域名命中黑名单: " + host,
		}
	}

	// Step 2: URL 特征规则检查
	if result := checkURLFeatures(parsedURL, host); result != nil {
		return result
	}

	// Step 3: LLM 风险评估（仅在有 LLM client 且有页面内容时执行）
	if llmClient != nil && pageContent != "" {
		analysis, err := llmClient.Summarize(ctx, pageContent)
		if err != nil {
			logx.Errorw("LLM safety analysis failed, defaulting to warning",
				logx.LogField{Key: "url", Value: targetURL},
				logx.LogField{Key: "err", Value: err.Error()})
			// LLM 分析失败不拦截，降级为 warning
			return &SafetyResult{
				IsSafe:    true,
				RiskLevel: "warning",
				RiskScore: 0.5,
				Reason:    "AI 安全分析暂不可用",
			}
		}

		riskLevel := scoreToLevel(analysis.RiskScore)
		return &SafetyResult{
			IsSafe:    riskLevel != "danger",
			RiskLevel: riskLevel,
			RiskScore: analysis.RiskScore,
			Reason:    analysis.RiskReason,
		}
	}

	return &SafetyResult{
		IsSafe:    true,
		RiskLevel: "safe",
		RiskScore: 0.0,
		Reason:    "",
	}
}

// checkURLFeatures 检查 URL 的可疑特征
func checkURLFeatures(parsedURL *url.URL, host string) *SafetyResult {
	// 检查 IP 直连（非域名访问）
	if ip := net.ParseIP(host); ip != nil {
		return &SafetyResult{
			IsSafe:    true,
			RiskLevel: "warning",
			RiskScore: 0.5,
			Reason:    "目标使用 IP 直连而非域名，存在一定风险",
		}
	}

	// 检查非标准端口
	port := parsedURL.Port()
	if port != "" && port != "80" && port != "443" {
		return &SafetyResult{
			IsSafe:    true,
			RiskLevel: "warning",
			RiskScore: 0.4,
			Reason:    "目标使用非标准端口: " + port,
		}
	}

	// 检查过长路径（可能隐藏恶意参数）
	if len(parsedURL.Path) > 256 {
		return &SafetyResult{
			IsSafe:    true,
			RiskLevel: "warning",
			RiskScore: 0.4,
			Reason:    "URL 路径异常过长",
		}
	}

	// 检查可疑参数编码（如大量%编码）
	rawQuery := parsedURL.RawQuery
	if len(rawQuery) > 0 {
		encodedCount := strings.Count(rawQuery, "%")
		if float64(encodedCount)/float64(len(rawQuery)) > 0.3 {
			return &SafetyResult{
				IsSafe:    true,
				RiskLevel: "warning",
				RiskScore: 0.5,
				Reason:    "URL 参数包含大量编码字符，可能隐藏恶意内容",
			}
		}
	}

	return nil
}

// scoreToLevel 将风险评分转为风险等级
func scoreToLevel(score float64) string {
	switch {
	case score >= 0.7:
		return "danger"
	case score >= 0.3:
		return "warning"
	default:
		return "safe"
	}
}
