package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// extractJSON 从 LLM 返回的文本中提取纯 JSON
// LLM 经常在 JSON 外面包裹 ```json ... ``` 或其他文字
func extractJSON(raw string) string {
	s := strings.TrimSpace(raw)

	// 去除 markdown 代码块标记: ```json ... ``` 或 ``` ... ```
	if strings.HasPrefix(s, "```") {
		// 找到第一个换行，跳过 ```json 行
		if idx := strings.Index(s, "\n"); idx != -1 {
			s = s[idx+1:]
		}
		// 去除结尾的 ```
		if idx := strings.LastIndex(s, "```"); idx != -1 {
			s = s[:idx]
		}
		s = strings.TrimSpace(s)
	}

	// 尝试找到 JSON 对象的边界 { ... }
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start != -1 && end != -1 && end > start {
		s = s[start : end+1]
	}

	return s
}

// PageAnalysis 表示 LLM 对页面的分析结果
type PageAnalysis struct {
	Summary    string   `json:"summary"`     // 页面摘要（200字以内）
	Keywords   []string `json:"keywords"`    // 关键词列表
	RiskScore  float64  `json:"risk_score"`  // 风险评分 0-1
	RiskReason string   `json:"risk_reason"` // 风险原因
	Slug       string   `json:"slug"`        // 语义化短链建议
}

// Client 是 OpenAI 兼容 API 的客户端
type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClient 创建一个新的 LLM 客户端
func NewClient(baseURL, apiKey, model string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// chatMessage 表示 chat completions API 的消息
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest 表示 chat completions API 的请求体
type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// chatResponse 表示 chat completions API 的响应体
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Summarize 调用 LLM 分析页面内容，返回摘要、关键词、风险评分和语义化短链
func (c *Client) Summarize(ctx context.Context, content string) (*PageAnalysis, error) {
	systemPrompt := `你是一个专业的网页内容分析助手。请分析用户提供的网页内容，并以严格的 JSON 格式返回以下信息：
{
  "summary": "页面摘要，不超过200字",
  "keywords": ["关键词1", "关键词2", "关键词3"],
  "risk_score": 0.0,
  "risk_reason": "如果有风险，说明原因；无风险则为空字符串",
  "slug": "基于内容语义生成的英文短链建议，使用小写字母和连字符，不超过30字符"
}

风险评分规则（0-1）：
- 0.0-0.3: 安全，正常内容
- 0.3-0.7: 存在一定风险，如包含广告、诱导性内容
- 0.7-1.0: 高风险，如钓鱼、欺诈、恶意软件

请只返回 JSON，不要返回其他内容。`

	result, err := c.callAPI(ctx, systemPrompt, content)
	if err != nil {
		return nil, fmt.Errorf("LLM API call failed: %w", err)
	}

	// 清洗 LLM 返回内容：去除 markdown 代码块标记等
	cleanJSON := extractJSON(result)

	var analysis PageAnalysis
	if err := json.Unmarshal([]byte(cleanJSON), &analysis); err != nil {
		logx.Errorw("failed to parse LLM response", logx.LogField{Key: "response", Value: result}, logx.LogField{Key: "err", Value: err.Error()})
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return &analysis, nil
}

// GenerateReadableSlug 基于页面内容生成语义化可读短链
func (c *Client) GenerateReadableSlug(ctx context.Context, content string) (string, error) {
	systemPrompt := `你是一个URL短链生成助手。根据用户提供的网页内容，生成一个语义化的英文短链slug。
规则：
1. 只使用小写字母、数字和连字符
2. 长度不超过20个字符
3. 要能反映页面核心内容
4. 只返回slug本身，不要返回其他内容

示例：
- Go语言教程 -> go-tutorial
- Python数据分析 -> python-data
- 2024年度报告 -> annual-report-2024`

	return c.callAPI(ctx, systemPrompt, content)
}

// callAPI 调用 OpenAI 兼容的 Chat Completions API
func (c *Client) callAPI(ctx context.Context, systemPrompt, userContent string) (string, error) {
	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userContent},
		},
		Temperature: 0.3,
		MaxTokens:   1024,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse API response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("LLM API returned no choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}
