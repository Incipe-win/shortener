package urltool

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

func GetBasePath(targetUrl string) (string, error) {
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		return "", err
	}
	if len(myUrl.Host) == 0 {
		return "", errors.New("no host in targetUrl")
	}
	return path.Base(myUrl.Path), nil
}

// IsCircularURL 检查输入 URL 是否指向本服务短域名，从域名层面直接拦截循环嵌套
// shortDomain 为配置中的短域名，如 "shortener.com"
func IsCircularURL(inputURL, shortDomain string) bool {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsedURL.Hostname())
	shortDomain = strings.ToLower(shortDomain)

	// 精确匹配或子域名匹配
	return host == shortDomain || strings.HasSuffix(host, "."+shortDomain)
}

