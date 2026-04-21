package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	maxContentLength = 2000 // 控制传给 LLM 的文本长度
	requestTimeout   = 10 * time.Second
	userAgent        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

// httpClient 全局 HTTP 客户端
var httpClient = &http.Client{
	Timeout: requestTimeout,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("too many redirects")
		}
		return nil
	},
}

// FetchPageContent 抓取指定 URL 的页面内容，提取标题、描述和正文
// 返回拼接后的文本（截取前 maxContentLength 字符），用于传给 LLM
func FetchPageContent(ctx context.Context, targetURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("URL returned status %d", resp.StatusCode)
	}

	// 限制读取 body 大小（1MB），避免下载过大文件
	limitReader := io.LimitReader(resp.Body, 1<<20)
	doc, err := html.Parse(limitReader)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	title, description, bodyText := extractContent(doc)

	var sb strings.Builder
	if title != "" {
		sb.WriteString("标题: " + title + "\n")
	}
	if description != "" {
		sb.WriteString("描述: " + description + "\n")
	}
	if bodyText != "" {
		sb.WriteString("正文: " + bodyText)
	}

	content := sb.String()
	if len(content) > maxContentLength {
		content = content[:maxContentLength]
	}

	return content, nil
}

// extractContent 从 HTML 文档中提取 title、meta description 和 body 文本
func extractContent(n *html.Node) (title, description, bodyText string) {
	var titleBuilder, bodyBuilder strings.Builder
	var inTitle, inBody, inScript, inStyle bool

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "title":
				inTitle = true
			case "body":
				inBody = true
			case "script", "noscript":
				inScript = true
			case "style":
				inStyle = true
			case "meta":
				var name, content string
				for _, attr := range node.Attr {
					if attr.Key == "name" {
						name = strings.ToLower(attr.Val)
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if name == "description" && content != "" {
					description = content
				}
			}
		}

		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if text != "" {
				if inTitle {
					titleBuilder.WriteString(text)
				}
				if inBody && !inScript && !inStyle {
					bodyBuilder.WriteString(text + " ")
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}

		// 退出当前节点时重置标记
		if node.Type == html.ElementNode {
			switch node.Data {
			case "title":
				inTitle = false
			case "script", "noscript":
				inScript = false
			case "style":
				inStyle = false
			}
		}
	}

	walk(n)

	title = strings.TrimSpace(titleBuilder.String())
	bodyText = strings.TrimSpace(bodyBuilder.String())
	return
}
