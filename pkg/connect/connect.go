package connect

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// client 全局的 HTTP 客户端
var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	Timeout: 2 * time.Second,
}

// Get 判断 url 是否可以请求通
func Get(url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorw("connect client.Get failed", logx.LogField{Key: "error", Value: err.Error()})
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
